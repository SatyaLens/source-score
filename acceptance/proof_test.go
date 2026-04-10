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

const (
	proof1Digest = "6f2479b5249b1c27c4935da5594bc72bb0b9e59e704aea9af50780bc6178c357"
	proof2Digest = "8df5229f310ae8322062834f3ba45a38ecef8ded549665d1170e15c8249b7cd0"
)

var _ = Describe("Proof model tests", func() {
	endpoint, err := url.JoinPath(baseUrl, "/api/v1/proof")
	Expect(err).To(BeNil())

	proofsEndpoint, err := url.JoinPath(baseUrl, "/api/v1/proofs")
	Expect(err).To(BeNil())

	// we need a source and a claim to attach proofs to
	srcEndpoint, err := url.JoinPath(baseUrl, "/api/v1/source")
	Expect(err).To(BeNil())

	claimEndpoint, err := url.JoinPath(baseUrl, "/api/v1/claim")
	Expect(err).To(BeNil())

	Context("Happy path tests", Ordered, func() {
		When("valid POST requests are sent to create proofs", func() {
			It("should create the proofs and return their uri digests", func() {
				// create source
				srcBody, err := json.Marshal(sourceInput3)
				Expect(err).To(BeNil())

				resp, err := http.Post(srcEndpoint, "application/json", bytes.NewBuffer(srcBody))
				Expect(err).To(BeNil())
				Expect(resp.StatusCode).To(BeEquivalentTo(http.StatusCreated))

				var srcResp api.CreateSourceResponse
				err = json.NewDecoder(resp.Body).Decode(&srcResp)
				Expect(err).To(BeNil())
				resp.Body.Close()

				// create claim attached to that source
				claimInput := api.ClaimInput{
					SourceUriDigest: srcResp.UriDigest,
					Summary:         sampleClaim1.Summary,
					Title:           sampleClaim1.Title,
					Uri:             sampleClaim1.Uri,
				}
				claimBody, err := json.Marshal(claimInput)
				Expect(err).To(BeNil())

				resp, err = http.Post(claimEndpoint, "application/json", bytes.NewBuffer(claimBody))
				Expect(err).To(BeNil())
				Expect(resp.StatusCode).To(BeEquivalentTo(http.StatusCreated))

				var claimResp map[string]string
				err = json.NewDecoder(resp.Body).Decode(&claimResp)
				Expect(err).To(BeNil())
				resp.Body.Close()
				claimDigest := claimResp["uriDigest"]

				// create proof 1
				proof1 := api.ProofInput{
					ClaimUriDigest: claimDigest,
					ReviewedBy:     "ReviewerA",
					SupportsClaim:  true,
					Uri:            "https://sample-proof-1",
				}
				body1, err := json.Marshal(proof1)
				Expect(err).To(BeNil())

				resp, err = http.Post(endpoint, "application/json", bytes.NewBuffer(body1))
				Expect(err).To(BeNil())
				Expect(resp.StatusCode).To(BeEquivalentTo(http.StatusCreated))

				var respBody map[string]string
				err = json.NewDecoder(resp.Body).Decode(&respBody)
				Expect(err).To(BeNil())
				resp.Body.Close()
				d1 := respBody["uriDigest"]
				Expect(d1).To(Equal(proof1Digest))

				// create proof 2
				proof2 := api.ProofInput{
					ClaimUriDigest: claimDigest,
					ReviewedBy:     "ReviewerB",
					SupportsClaim:  false,
					Uri:            "https://sample-proof-2",
				}
				body2, err := json.Marshal(proof2)
				Expect(err).To(BeNil())

				resp, err = http.Post(endpoint, "application/json", bytes.NewBuffer(body2))
				Expect(err).To(BeNil())
				Expect(resp.StatusCode).To(BeEquivalentTo(http.StatusCreated))

				err = json.NewDecoder(resp.Body).Decode(&respBody)
				Expect(err).To(BeNil())
				resp.Body.Close()
				d2 := respBody["uriDigest"]
				Expect(d2).To(Equal(proof2Digest))
			})
		})

		When("GET requests are sent to fetch proofs", func() {
			It("should return all the created proofs and individual proof", func() {
				resp, err := http.Get(proofsEndpoint)
				Expect(err).To(BeNil())
				defer resp.Body.Close()
				Expect(resp.StatusCode).To(Equal(http.StatusOK))

				var proofs []api.Proof
				err = json.NewDecoder(resp.Body).Decode(&proofs)
				Expect(err).To(BeNil())
				Expect(len(proofs)).To(BeNumerically(">=", 2))

				// GET single proof
				proofUrl, err := url.JoinPath(endpoint, proof1Digest)
				Expect(err).To(BeNil())
				resp, err = http.Get(proofUrl)
				Expect(err).To(BeNil())
				defer resp.Body.Close()
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
			})

			It("should return the correct proof when requested by uri digest", func() {
				proofUrl, err := url.JoinPath(endpoint, proof1Digest)
				Expect(err).To(BeNil())

				resp, err := http.Get(proofUrl)
				Expect(err).To(BeNil())
				defer resp.Body.Close()
				Expect(resp.StatusCode).To(Equal(http.StatusOK))

				var p api.Proof
				err = json.NewDecoder(resp.Body).Decode(&p)
				Expect(err).To(BeNil())
				Expect(p.UriDigest).To(Equal(proof1Digest))
			})
		})

		When("PATCH request is sent to update a proof", func() {
			It("should update the proof and subsequent GET returns updated record", func() {
				proofUrl, err := url.JoinPath(endpoint, proof1Digest)
				Expect(err).To(BeNil())

				patchBody := api.ProofPatchInput{ReviewedBy: "UpdatedReviewer"}
				body, err := json.Marshal(patchBody)
				Expect(err).To(BeNil())

				req, err := http.NewRequest(http.MethodPatch, proofUrl, bytes.NewBuffer(body))
				Expect(err).To(BeNil())
				req.Header.Set("Content-Type", "application/json")

				resp, err := client.Do(req)
				Expect(err).To(BeNil())
				defer resp.Body.Close()
				Expect(resp.StatusCode).To(Equal(http.StatusNoContent))

				// verify update
				resp, err = http.Get(proofUrl)
				Expect(err).To(BeNil())
				defer resp.Body.Close()
				Expect(resp.StatusCode).To(Equal(http.StatusOK))

				var p api.Proof
				err = json.NewDecoder(resp.Body).Decode(&p)
				Expect(err).To(BeNil())
				Expect(p.ReviewedBy).To(Equal("UpdatedReviewer"))
			})
		})

		When("DELETE request is sent to delete a proof", func() {
			It("should delete the created proof and subsequent GET returns 404", func() {
				proofUrl, err := url.JoinPath(endpoint, proof1Digest)
				Expect(err).To(BeNil())

				req, err := http.NewRequest(http.MethodDelete, proofUrl, nil)
				Expect(err).To(BeNil())

				resp, err := client.Do(req)
				Expect(err).To(BeNil())
				defer resp.Body.Close()
				Expect(resp.StatusCode).To(Equal(http.StatusNoContent))

				// verify deletion
				resp, err = http.Get(proofUrl)
				Expect(err).To(BeNil())
				defer resp.Body.Close()
				Expect(resp.StatusCode).To(Equal(http.StatusNotFound))
			})
		})
	})
})
