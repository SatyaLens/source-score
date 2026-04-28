package acceptance_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"source-score/pkg/api"
	"strings"
	"time"

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

				resp, err := doRequest(http.MethodPost, srcEndpoint, srcBody)
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

				resp, err = doRequest(http.MethodPost, endpoint, body1)
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

				resp, err = doRequest(http.MethodPost, endpoint, body2)
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
				resp, err := doRequest(http.MethodGet, claimsEndpoint, nil)
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

				resp, err := doRequest(http.MethodGet, claimUrl, nil)
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
				addCommonHeaders(req)
				req.Header.Set("Content-Type", "application/json")

				resp, err := client.Do(req)
				Expect(err).To(BeNil())
				defer resp.Body.Close()
				Expect(resp.StatusCode).To(Equal(http.StatusNoContent))

				// verify update
				resp, err = doRequest(http.MethodGet, claimUrl, nil)
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
				addCommonHeaders(req)
				req.Header.Set("Content-Type", "application/json")

				resp, err := client.Do(req)
				Expect(err).To(BeNil())
				defer resp.Body.Close()
				Expect(resp.StatusCode).To(Equal(http.StatusNoContent))

				// verify update
				resp, err = doRequest(http.MethodGet, claimUrl, nil)
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
				addCommonHeaders(req)
				req.Header.Set("Content-Type", "application/json")

				resp, err := client.Do(req)
				Expect(err).To(BeNil())
				defer resp.Body.Close()
				Expect(resp.StatusCode).To(Equal(http.StatusNoContent))

				// verify update
				resp, err = doRequest(http.MethodGet, claimUrl, nil)
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

		// TODO: uncomment when single claim verification is enabled
		// When("POST request is sent to verify a claim", func() {
		// 	It("should verify the claim and subsequent GET returns updated checked and validity fields", func() {
		// 		claimUrl, err := url.JoinPath(endpoint, claim2Digest)
		// 		Expect(err).To(BeNil())

		// 		// verify claim with validity true
		// 		validity := true
		// 		verifyBody := api.ClaimVerification{
		// 			Validity: &validity,
		// 		}
		// 		body, err := json.Marshal(verifyBody)
		// 		Expect(err).To(BeNil())

		// 		resp, err := doRequest(http.MethodPost, claimUrl, body)
		// 		Expect(err).To(BeNil())
		// 		defer resp.Body.Close()
		// 		Expect(resp.StatusCode).To(Equal(http.StatusNoContent))

		// 		// verify the claim was updated
		// 		resp, err = doRequest(http.MethodGet, claimUrl, nil)
		// 		Expect(err).To(BeNil())
		// 		defer resp.Body.Close()
		// 		Expect(resp.StatusCode).To(Equal(http.StatusOK))

		// 		var c api.Claim
		// 		err = json.NewDecoder(resp.Body).Decode(&c)
		// 		Expect(err).To(BeNil())
		// 		Expect(c.Checked).To(BeTrue())
		// 		Expect(c.Validity).To(BeTrue())
		// 	})
		// })

		When("DELETE request is sent to delete the created claim", func() {
			It("should delete the created claim and subsequent GET returns 404", func() {
				claimUrl, err := url.JoinPath(endpoint, claim1Digest)
				Expect(err).To(BeNil())

				req, err := http.NewRequest(http.MethodDelete, claimUrl, nil)
				Expect(err).To(BeNil())
				addCommonHeaders(req)

				resp, err := client.Do(req)
				Expect(err).To(BeNil())
				defer resp.Body.Close()
				Expect(resp.StatusCode).To(Equal(http.StatusNoContent))

				// verify deletion
				resp, err = doRequest(http.MethodGet, claimUrl, nil)
				Expect(err).To(BeNil())
				defer resp.Body.Close()
				Expect(resp.StatusCode).To(Equal(http.StatusNotFound))
			})
		})

		When("Verifying all claims based on their proofs", func() {
			It("Should verify claims with validity based on proof counts", func() {
				// Create a new source
				srcInput := api.SourceInput{
					Name:    "Test Source for Verify All",
					Summary: "sample summary",
					Tags:    "tag45",
					Uri:     "https://test-source-verify-all",
				}
				srcBody, err := json.Marshal(srcInput)
				Expect(err).To(BeNil())

				resp, err := doRequest(http.MethodPost, srcEndpoint, srcBody)
				Expect(err).To(BeNil())
				Expect(resp.StatusCode).To(BeEquivalentTo(http.StatusCreated))

				var srcResp api.CreateSourceResponse
				err = json.NewDecoder(resp.Body).Decode(&srcResp)
				Expect(err).To(BeNil())
				resp.Body.Close()

				// Create claim 1
				claim1Input := api.ClaimInput{
					SourceUriDigest: srcResp.UriDigest,
					Summary:         "Test claim 1 for verify all",
					Title:           "Test Claim 1",
					Uri:             "https://test-claim-verify-1",
				}
				body1, err := json.Marshal(claim1Input)
				Expect(err).To(BeNil())

				resp, err = doRequest(http.MethodPost, endpoint, body1)
				Expect(err).To(BeNil())
				Expect(resp.StatusCode).To(BeEquivalentTo(http.StatusCreated))

				var respBody map[string]string
				err = json.NewDecoder(resp.Body).Decode(&respBody)
				Expect(err).To(BeNil())
				resp.Body.Close()
				testClaim1Digest := respBody["uriDigest"]

				// Create claim 2
				claim2Input := api.ClaimInput{
					SourceUriDigest: srcResp.UriDigest,
					Summary:         "Test claim 2 for verify all",
					Title:           "Test Claim 2",
					Uri:             "https://test-claim-verify-2",
				}
				body2, err := json.Marshal(claim2Input)
				Expect(err).To(BeNil())

				resp, err = doRequest(http.MethodPost, endpoint, body2)
				Expect(err).To(BeNil())
				Expect(resp.StatusCode).To(BeEquivalentTo(http.StatusCreated))

				err = json.NewDecoder(resp.Body).Decode(&respBody)
				Expect(err).To(BeNil())
				resp.Body.Close()
				testClaim2Digest := respBody["uriDigest"]

				// Create proofs for claim 1 (3 supporting, 1 refuting) -> validity = true
				proofEndpoint, err := url.JoinPath(baseUrl, "/api/v1/proof")
				Expect(err).To(BeNil())

				for i := range 3 {
					supports := true
					proofInput := api.ProofInput{
						ClaimUriDigest: testClaim1Digest,
						ReviewedBy:     "ReviewerA",
						SupportsClaim:  supports,
						Uri:            fmt.Sprintf("https://proof-claim1-support-%d", i),
					}
					proofBody, _ := json.Marshal(proofInput)
					resp, err = doRequest(http.MethodPost, proofEndpoint, proofBody)
					Expect(err).To(BeNil())
					Expect(resp.StatusCode).To(Equal(http.StatusCreated))
					resp.Body.Close()
				}

				supports := false
				proofInput := api.ProofInput{
					ClaimUriDigest: testClaim1Digest,
					ReviewedBy:     "ReviewerB",
					SupportsClaim:  supports,
					Uri:            "https://proof-claim1-refute",
				}
				proofBody, _ := json.Marshal(proofInput)
				resp, err = doRequest(http.MethodPost, proofEndpoint, proofBody)
				Expect(err).To(BeNil())
				Expect(resp.StatusCode).To(Equal(http.StatusCreated))
				resp.Body.Close()

				// Create proofs for claim 2 (1 supporting, 2 refuting) -> validity = false
				supports = true
				proofInput = api.ProofInput{
					ClaimUriDigest: testClaim2Digest,
					ReviewedBy:     "ReviewerC",
					SupportsClaim:  supports,
					Uri:            "https://proof-claim2-support",
				}
				proofBody, _ = json.Marshal(proofInput)
				resp, err = doRequest(http.MethodPost, proofEndpoint, proofBody)
				Expect(err).To(BeNil())
				Expect(resp.StatusCode).To(Equal(http.StatusCreated))
				resp.Body.Close()

				for i := range 2 {
					supports := false
					proofInput := api.ProofInput{
						ClaimUriDigest: testClaim2Digest,
						ReviewedBy:     "ReviewerD",
						SupportsClaim:  supports,
						Uri:            fmt.Sprintf("https://proof-claim2-refute-%d", i),
					}
					proofBody, _ := json.Marshal(proofInput)
					resp, err = doRequest(http.MethodPost, proofEndpoint, proofBody)
					Expect(err).To(BeNil())
					Expect(resp.StatusCode).To(Equal(http.StatusCreated))
					resp.Body.Close()
				}

				// Hit the verify all claims endpoint
				verifyAllUrl, err := url.JoinPath(baseUrl, "/api/v1/claims/verify")
				Expect(err).To(BeNil())

				resp, err = doRequest(http.MethodPost, verifyAllUrl, nil)
				Expect(err).To(BeNil())
				defer resp.Body.Close()
				Expect(resp.StatusCode).To(Equal(http.StatusAccepted))

				Eventually(func(g Gomega) {
					// Assert claim 1 has validity=true (3 supporting > 1 refuting)
					claim1Url, err := url.JoinPath(endpoint, testClaim1Digest)
					g.Expect(err).To(BeNil())

					resp, err = doRequest(http.MethodGet, claim1Url, nil)
					g.Expect(err).To(BeNil())
					g.Expect(resp.StatusCode).To(Equal(http.StatusOK))

					var c1 api.Claim
					err = json.NewDecoder(resp.Body).Decode(&c1)
					g.Expect(err).To(BeNil())
					g.Expect(c1.Checked).To(BeTrue())
					g.Expect(c1.Validity).To(BeTrue())
					resp.Body.Close()

					// Assert claim 2 has validity=false (1 supporting < 2 refuting)
					claim2Url, err := url.JoinPath(endpoint, testClaim2Digest)
					g.Expect(err).To(BeNil())

					resp, err = doRequest(http.MethodGet, claim2Url, nil)
					g.Expect(err).To(BeNil())
					g.Expect(resp.StatusCode).To(Equal(http.StatusOK))

					var c2 api.Claim
					err = json.NewDecoder(resp.Body).Decode(&c2)
					g.Expect(err).To(BeNil())
					g.Expect(c2.Checked).To(BeTrue())
					g.Expect(c2.Validity).To(BeFalse())
					resp.Body.Close()
				}, 10*time.Second, 1*time.Second).Should(Succeed())
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

				resp, err := doRequest(http.MethodPost, endpoint, body)
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

				resp, err := doRequest(http.MethodPost, endpoint, body)
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

				resp, err := doRequest(http.MethodPost, endpoint, body)
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

				resp, err := doRequest(http.MethodGet, claimUrl, nil)
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
				addCommonHeaders(req)

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
				addCommonHeaders(req)
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
				addCommonHeaders(req)
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

		// TODO: uncomment when single claim verification is enabled
		// When("POST request to verify claim without Validity field is sent", func() {
		// 	It("should return 400 with validation error mentioning validity", func() {
		// 		claimUrl, err := url.JoinPath(endpoint, claim2Digest)
		// 		Expect(err).To(BeNil())

		// 		verifyBody := api.ClaimVerification{}
		// 		body, err := json.Marshal(verifyBody)
		// 		Expect(err).To(BeNil())

		// 		resp, err := doRequest(http.MethodPost, claimUrl, body)
		// 		Expect(err).To(BeNil())
		// 		defer resp.Body.Close()
		// 		Expect(resp.StatusCode).To(Equal(http.StatusBadRequest))

		// 		var errResp map[string]string
		// 		err = json.NewDecoder(resp.Body).Decode(&errResp)
		// 		Expect(err).To(BeNil())
		// 		Expect(strings.ToLower(errResp["error"])).To(ContainSubstring("validity"))
		// 		Expect(strings.ToLower(errResp["error"])).To(ContainSubstring("required"))
		// 	})
		// })

		// TODO: uncomment when single claim verification is enabled
		// When("POST request to verify a non-existent claim is sent", func() {
		// 	It("should return 404 with claim not found error", func() {
		// 		claimUrl, err := url.JoinPath(endpoint, "doesnotexist")
		// 		Expect(err).To(BeNil())

		// 		validity := true
		// 		verifyBody := api.ClaimVerification{
		// 			Validity: &validity,
		// 		}
		// 		body, err := json.Marshal(verifyBody)
		// 		Expect(err).To(BeNil())

		// 		resp, err := doRequest(http.MethodPost, claimUrl, body)
		// 		Expect(err).To(BeNil())
		// 		defer resp.Body.Close()
		// 		Expect(resp.StatusCode).To(Equal(http.StatusNotFound))

		// 		var errResp map[string]string
		// 		err = json.NewDecoder(resp.Body).Decode(&errResp)
		// 		Expect(err).To(BeNil())
		// 		Expect(strings.ToLower(errResp["error"])).To(ContainSubstring("claim not found"))
		// 	})
		// })
	})
})
