package acceptance_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"
	"source-score/pkg/api"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

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
					Summary:         &sampleClaim1.Summary,
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
					Summary:         &sampleClaim2.Summary,
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

		When("GET requests are sent to create claims", func() {
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
})
