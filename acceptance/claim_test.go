package acceptance_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"
	"source-score/pkg/api"
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func ptrString(s string) *string { return &s }

var _ = Describe("Claim model tests", func() {
	// claim endpoints
	endpoint, err := url.JoinPath(baseUrl, "/api/v1/claim")
	Expect(err).To(BeNil())

	claimsEndpoint, err := url.JoinPath(baseUrl, "/api/v1/claims")
	Expect(err).To(BeNil())

	// we'll create a source to attach claims to
	srcEndpoint, err := url.JoinPath(baseUrl, "/api/v1/source")
	Expect(err).To(BeNil())

	Context("Happy path tests", Ordered, func() {
		When("valid POST requests are sent to create claims", func() {
			It("should create the claims and return their uri digests", func() {
				srcBody, err := json.Marshal(sourceInput3)
				Expect(err).To(BeNil())

				resp, err := http.Post(srcEndpoint, "application/json", bytes.NewBuffer(srcBody))
				Expect(err).To(BeNil())
				Expect(resp.StatusCode).To(BeEquivalentTo(http.StatusCreated))

				var srcResp api.CreateSourceResponse
				err = json.NewDecoder(resp.Body).Decode(&srcResp)
				Expect(err).To(BeNil())
				resp.Body.Close()

				claim1 := api.ClaimInput{
					SourceUriDigest: srcResp.UriDigest,
					Summary:         sampleClaim1.Summary,
					Title:           sampleClaim1.Title,
					Uri:             sampleClaim1.Uri,
				}
				body1, err := json.Marshal(claim1)
				Expect(err).To(BeNil())

				resp, err = http.Post(endpoint, "application/json", bytes.NewBuffer(body1))
				Expect(err).To(BeNil())
				Expect(resp.StatusCode).To(BeEquivalentTo(http.StatusCreated))

				var respBody map[string]string
				err = json.NewDecoder(resp.Body).Decode(&respBody)
				Expect(err).To(BeNil())
				resp.Body.Close()
				digest1 := respBody["uriDigest"]
				Expect(digest1).To(Equal(claim1Digest))

				claim2 := api.ClaimInput{
					SourceUriDigest: srcResp.UriDigest,
					Summary:         sampleClaim2.Summary,
					Title:           sampleClaim2.Title,
					Uri:             sampleClaim2.Uri,
				}
				body2, err := json.Marshal(claim2)
				Expect(err).To(BeNil())

				resp, err = http.Post(endpoint, "application/json", bytes.NewBuffer(body2))
				Expect(err).To(BeNil())
				Expect(resp.StatusCode).To(BeEquivalentTo(http.StatusCreated))

				err = json.NewDecoder(resp.Body).Decode(&respBody)
				Expect(err).To(BeNil())
				resp.Body.Close()
				digest2 := respBody["uriDigest"]
				Expect(digest2).To(Equal(claim2Digest))

			})
		})

		When("GET requests are sent to fetch claims", func() {
			It("should return all the created claims", func() {
				resp, err := http.Get(claimsEndpoint)
				Expect(err).To(BeNil())
				defer resp.Body.Close()
				Expect(resp.StatusCode).To(Equal(http.StatusOK))

				var claims []api.Claim
				err = json.NewDecoder(resp.Body).Decode(&claims)
				Expect(err).To(BeNil())
				Expect(len(claims)).To(BeNumerically(">=", 2))

				// assert the created claims are present
				Expect(claims).To(ContainElements(
					sampleClaim1,
					sampleClaim2,
				))
			})
		})

		When("GET request is sent to retrieve a single claim by digest", func() {
			It("should return the created claim", func() {
				claimUrl, err := url.JoinPath(endpoint, claim1Digest)
				Expect(err).To(BeNil())

				resp, err := http.Get(claimUrl)
				Expect(err).To(BeNil())
				defer resp.Body.Close()
				Expect(resp.StatusCode).To(Equal(http.StatusOK))

				var c api.Claim
				err = json.NewDecoder(resp.Body).Decode(&c)
				Expect(err).To(BeNil())
				Expect(c).To(Equal(sampleClaim1))
			})
		})

		When("PATCH request is sent to update the created claim", func() {
			It("should update the claim and subsequent GET returns updated record", func() {
				claimUrl, err := url.JoinPath(endpoint, claim1Digest)
				Expect(err).To(BeNil())

				updatedTitle := "Patched Claim Title"
				updatedSummary := "Patched claim summary"
				patchBody := api.ClaimPatchInput{
					Title:   &updatedTitle,
					Summary: &updatedSummary,
				}
				body, err := json.Marshal(patchBody)
				Expect(err).To(BeNil())

				req, err := http.NewRequest(http.MethodPatch, claimUrl, bytes.NewBuffer(body))
				Expect(err).To(BeNil())
				req.Header.Set("Content-Type", "application/json")

				resp, err := client.Do(req)
				Expect(err).To(BeNil())
				defer resp.Body.Close()
				Expect(resp.StatusCode).To(Equal(http.StatusNoContent))

				// verify update
				resp, err = http.Get(claimUrl)
				Expect(err).To(BeNil())
				defer resp.Body.Close()
				Expect(resp.StatusCode).To(Equal(http.StatusOK))

				var c api.Claim
				err = json.NewDecoder(resp.Body).Decode(&c)
				Expect(err).To(BeNil())
				Expect(c.Title).To(Equal(updatedTitle))
				Expect(c.Summary).To(Equal(updatedSummary))
			})

			It("should update only the summary when only summary is provided", func() {
				claimUrl, err := url.JoinPath(endpoint, claim1Digest)
				Expect(err).To(BeNil())

				updatedSummary2 := "Patched Claim Summary Only"
				patchBody := api.ClaimPatchInput{
					Summary: &updatedSummary2,
				}
				body, err := json.Marshal(patchBody)
				Expect(err).To(BeNil())

				req, err := http.NewRequest(http.MethodPatch, claimUrl, bytes.NewBuffer(body))
				Expect(err).To(BeNil())
				req.Header.Set("Content-Type", "application/json")

				resp, err := client.Do(req)
				Expect(err).To(BeNil())
				defer resp.Body.Close()
				Expect(resp.StatusCode).To(Equal(http.StatusNoContent))

				// verify update
				resp, err = http.Get(claimUrl)
				Expect(err).To(BeNil())
				defer resp.Body.Close()
				Expect(resp.StatusCode).To(Equal(http.StatusOK))

				var c api.Claim
				err = json.NewDecoder(resp.Body).Decode(&c)
				Expect(err).To(BeNil())
				// title should remain as previously patched
				Expect(c.Title).To(Equal("Patched Claim Title"))
				Expect(c.Summary).To(Equal(updatedSummary2))
			})

			It("should update only the title when only title is provided", func() {
				claimUrl, err := url.JoinPath(endpoint, claim1Digest)
				Expect(err).To(BeNil())

				updatedTitle2 := "Patched Claim Title Only"
				patchBody := api.ClaimPatchInput{
					Title: &updatedTitle2,
				}
				body, err := json.Marshal(patchBody)
				Expect(err).To(BeNil())

				req, err := http.NewRequest(http.MethodPatch, claimUrl, bytes.NewBuffer(body))
				Expect(err).To(BeNil())
				req.Header.Set("Content-Type", "application/json")

				resp, err := client.Do(req)
				Expect(err).To(BeNil())
				defer resp.Body.Close()
				Expect(resp.StatusCode).To(Equal(http.StatusNoContent))

				// verify update
				resp, err = http.Get(claimUrl)
				Expect(err).To(BeNil())
				defer resp.Body.Close()
				Expect(resp.StatusCode).To(Equal(http.StatusOK))

				var c api.Claim
				err = json.NewDecoder(resp.Body).Decode(&c)
				Expect(err).To(BeNil())
				// summary should remain as previously patched in prior test
				Expect(c.Summary).To(Equal("Patched Claim Summary Only"))
				Expect(c.Title).To(Equal(updatedTitle2))
			})
		})

		When("DELETE request is sent to delete the created claim", func() {
			It("should delete the created claim and subsequent GET returns 404", func() {
				claimUrl, err := url.JoinPath(endpoint, claim1Digest)
				Expect(err).To(BeNil())

				req, err := http.NewRequest(http.MethodDelete, claimUrl, nil)
				Expect(err).To(BeNil())

				resp, err := client.Do(req)
				Expect(err).To(BeNil())
				defer resp.Body.Close()
				Expect(resp.StatusCode).To(Equal(http.StatusNoContent))

				// verify deletion
				resp, err = http.Get(claimUrl)
				Expect(err).To(BeNil())
				defer resp.Body.Close()
				Expect(resp.StatusCode).To(Equal(http.StatusNotFound))
			})
		})
	})

	Context("Validation tests", func() {
		When("POST request with empty sourceUriDigest is sent", func() {
			It("should return 400 with validation error", func() {
				claim := api.ClaimInput{
					SourceUriDigest: "",
					Summary:         "summary",
					Title:           "title",
					Uri:             "https://ok",
				}
				body, err := json.Marshal(claim)
				Expect(err).To(BeNil())

				resp, err := http.Post(endpoint, "application/json", bytes.NewBuffer(body))
				Expect(err).To(BeNil())
				defer resp.Body.Close()
				Expect(resp.StatusCode).To(Equal(http.StatusBadRequest))

				var errResp map[string]string
				err = json.NewDecoder(resp.Body).Decode(&errResp)
				Expect(err).To(BeNil())
				Expect(strings.ToLower(errResp["error"])).To(ContainSubstring("sourceuridigest"))
				Expect(strings.ToLower(errResp["error"])).To(ContainSubstring("required"))
			})
		})

		When("POST request with empty title and summary is sent", func() {
			It("should return 400 with validation error mentioning Title and Summary", func() {
				claim := api.ClaimInput{
					SourceUriDigest: "somedigest",
					Summary:         "",
					Title:           "",
					Uri:             "https://ok",
				}
				body, err := json.Marshal(claim)
				Expect(err).To(BeNil())

				resp, err := http.Post(endpoint, "application/json", bytes.NewBuffer(body))
				Expect(err).To(BeNil())
				defer resp.Body.Close()
				Expect(resp.StatusCode).To(Equal(http.StatusBadRequest))

				var errResp map[string]string
				err = json.NewDecoder(resp.Body).Decode(&errResp)
				Expect(err).To(BeNil())
				Expect(strings.ToLower(errResp["error"])).To(ContainSubstring("title"))
				Expect(strings.ToLower(errResp["error"])).To(ContainSubstring("summary"))
				Expect(strings.ToLower(errResp["error"])).To(ContainSubstring("required"))
			})
		})

		When("POST request with non-https Uri is sent", func() {
			It("should return 400 with validation error mentioning Uri and httpsurl", func() {
				claim := api.ClaimInput{
					SourceUriDigest: "somedigest",
					Summary:         "summary",
					Title:           "title",
					Uri:             "http://not-https",
				}
				body, err := json.Marshal(claim)
				Expect(err).To(BeNil())

				resp, err := http.Post(endpoint, "application/json", bytes.NewBuffer(body))
				Expect(err).To(BeNil())
				defer resp.Body.Close()
				Expect(resp.StatusCode).To(Equal(http.StatusBadRequest))

				var errResp map[string]string
				err = json.NewDecoder(resp.Body).Decode(&errResp)
				Expect(err).To(BeNil())
				Expect(strings.ToLower(errResp["error"])).To(ContainSubstring("uri"))
				Expect(strings.ToLower(errResp["error"])).To(ContainSubstring("httpsurl"))
			})
		})

		When("GET request is sent for a non-existent claim", func() {
			It("should return 404 with claim not found error", func() {
				claimUrl, err := url.JoinPath(endpoint, "doesnotexist")
				Expect(err).To(BeNil())

				resp, err := http.Get(claimUrl)
				Expect(err).To(BeNil())
				defer resp.Body.Close()
				Expect(resp.StatusCode).To(Equal(http.StatusNotFound))

				var errResp map[string]string
				err = json.NewDecoder(resp.Body).Decode(&errResp)
				Expect(err).To(BeNil())
				Expect(strings.ToLower(errResp["error"])).To(ContainSubstring("claim not found"))
			})
		})

		When("DELETE request is sent for a non-existent claim", func() {
			It("should return 404 with claim not found error", func() {
				claimUrl, err := url.JoinPath(endpoint, "doesnotexistdigest")
				Expect(err).To(BeNil())

				req, err := http.NewRequest(http.MethodDelete, claimUrl, nil)
				Expect(err).To(BeNil())

				resp, err := client.Do(req)
				Expect(err).To(BeNil())
				defer resp.Body.Close()
				Expect(resp.StatusCode).To(Equal(http.StatusNotFound))

				var errResp map[string]string
				err = json.NewDecoder(resp.Body).Decode(&errResp)
				Expect(err).To(BeNil())
				Expect(strings.ToLower(errResp["error"])).To(ContainSubstring("claim not found"))
			})
		})

		When("PATCH request is sent for a non-existent claim", func() {
			It("should return 404 with claim not found error", func() {
				claimUrl, err := url.JoinPath(endpoint, "doesnotexist")
				Expect(err).To(BeNil())

				patchBody := api.ClaimPatchInput{
					Title:   ptrString("irrelevant"),
					Summary: ptrString("irrelevant summary"),
				}
				body, err := json.Marshal(patchBody)
				Expect(err).To(BeNil())

				req, err := http.NewRequest(http.MethodPatch, claimUrl, bytes.NewBuffer(body))
				Expect(err).To(BeNil())
				req.Header.Set("Content-Type", "application/json")

				resp, err := client.Do(req)
				Expect(err).To(BeNil())
				defer resp.Body.Close()
				Expect(resp.StatusCode).To(Equal(http.StatusNotFound))

				var errResp map[string]string
				err = json.NewDecoder(resp.Body).Decode(&errResp)
				Expect(err).To(BeNil())
				Expect(strings.ToLower(errResp["error"])).To(ContainSubstring("claim not found"))
			})
		})

		When("PATCH request with empty title and summary is sent", func() {
			It("should return 400 with validation error mentioning Title and Summary", func() {
				claimUrl, err := url.JoinPath(endpoint, claim1Digest)
				Expect(err).To(BeNil())

				patchBody := api.ClaimPatchInput{
					Title:   ptrString(""),
					Summary: ptrString(""),
				}
				body, err := json.Marshal(patchBody)
				Expect(err).To(BeNil())

				req, err := http.NewRequest(http.MethodPatch, claimUrl, bytes.NewBuffer(body))
				Expect(err).To(BeNil())
				req.Header.Set("Content-Type", "application/json")

				resp, err := client.Do(req)
				Expect(err).To(BeNil())
				defer resp.Body.Close()
				Expect(resp.StatusCode).To(Equal(http.StatusBadRequest))

				var errResp map[string]string
				err = json.NewDecoder(resp.Body).Decode(&errResp)
				Expect(err).To(BeNil())
				Expect(strings.ToLower(errResp["error"])).To(ContainSubstring("title"))
				Expect(strings.ToLower(errResp["error"])).To(ContainSubstring("summary"))
				Expect(strings.ToLower(errResp["error"])).To(ContainSubstring("nonempty"))
			})
		})
	})
})
