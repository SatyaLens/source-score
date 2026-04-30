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
				srcBody, err := json.Marshal(sourceInput4)
				Expect(err).To(BeNil())

				resp, err := doRequest(http.MethodPost, srcEndpoint, srcBody)
				Expect(err).To(BeNil())
				Expect(resp.StatusCode).To(BeEquivalentTo(http.StatusCreated))

				var srcResp api.CreateSourceResponse
				err = json.NewDecoder(resp.Body).Decode(&srcResp)
				Expect(err).To(BeNil())
				defer resp.Body.Close()
				Expect(srcResp.UriDigest).To(Equal(uriDigest4))

				// create claim attached to that source
				claimInput := api.ClaimInput{
					SourceUriDigest: sampleClaim3.SourceUriDigest,
					Summary:         sampleClaim3.Summary,
					Title:           sampleClaim3.Title,
					Uri:             sampleClaim3.Uri,
				}
				claimBody, err := json.Marshal(claimInput)
				Expect(err).To(BeNil())

				resp, err = doRequest(http.MethodPost, claimEndpoint, claimBody)
				Expect(err).To(BeNil())
				Expect(resp.StatusCode).To(BeEquivalentTo(http.StatusCreated))

				var claimResp map[string]string
				err = json.NewDecoder(resp.Body).Decode(&claimResp)
				Expect(err).To(BeNil())
				resp.Body.Close()
				claimDigest := claimResp["uriDigest"]
				Expect(claimDigest).To(Equal(claim3Digest))

				// create proof 1
				proof1 := api.ProofInput{
					ClaimUriDigest: sampleProof1.ClaimUriDigest,
					ReviewedBy:     sampleProof1.ReviewedBy,
					SupportsClaim:  &sampleProof1.SupportsClaim,
					Uri:            sampleProof1.Uri,
				}
				body1, err := json.Marshal(proof1)
				Expect(err).To(BeNil())

				resp, err = doRequest(http.MethodPost, endpoint, body1)
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
					ClaimUriDigest: sampleProof2.ClaimUriDigest,
					ReviewedBy:     sampleProof2.ReviewedBy,
					SupportsClaim:  &sampleProof2.SupportsClaim,
					Uri:            sampleProof2.Uri,
				}
				body2, err := json.Marshal(proof2)
				Expect(err).To(BeNil())

				resp, err = doRequest(http.MethodPost, endpoint, body2)
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
				resp, err := doRequest(http.MethodGet, proofsEndpoint, nil)
				Expect(err).To(BeNil())
				defer resp.Body.Close()
				Expect(resp.StatusCode).To(Equal(http.StatusOK))

				var proofs []api.Proof
				err = json.NewDecoder(resp.Body).Decode(&proofs)
				Expect(err).To(BeNil())
				Expect(len(proofs)).To(BeNumerically(">=", 2))
				Expect(proofs).To(ContainElements(sampleProof1, sampleProof2))
			})

			It("should return the correct proof when requested by uri digest", func() {
				proofUrl, err := url.JoinPath(endpoint, proof1Digest)
				Expect(err).To(BeNil())

				resp, err := doRequest(http.MethodGet, proofUrl, nil)
				Expect(err).To(BeNil())
				defer resp.Body.Close()
				Expect(resp.StatusCode).To(Equal(http.StatusOK))

				var p api.Proof
				err = json.NewDecoder(resp.Body).Decode(&p)
				Expect(err).To(BeNil())
				defer resp.Body.Close()
				Expect(p.UriDigest).To(Equal(proof1Digest))
				Expect(p).To(Equal(sampleProof1))
			})

			It("should return all proofs for the given claim digest", func() {
				claimProofsUrl, err := url.JoinPath(baseUrl, "/api/v1/claim", claim3Digest, "proofs")
				Expect(err).To(BeNil())

				resp, err := doRequest(http.MethodGet, claimProofsUrl, nil)
				Expect(err).To(BeNil())
				defer resp.Body.Close()
				Expect(resp.StatusCode).To(Equal(http.StatusOK))

				var proofs []api.Proof
				err = json.NewDecoder(resp.Body).Decode(&proofs)
				Expect(err).To(BeNil())
				Expect(len(proofs)).To(BeNumerically(">=", 2))
				Expect(proofs).To(ContainElements(sampleProof1, sampleProof2))
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
				addCommonHeaders(req)
				req.Header.Set("Content-Type", "application/json")

				resp, err := client.Do(req)
				Expect(err).To(BeNil())
				defer resp.Body.Close()
				Expect(resp.StatusCode).To(Equal(http.StatusNoContent))

				// verify update
				resp, err = doRequest(http.MethodGet, proofUrl, nil)
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
				addCommonHeaders(req)

				resp, err := client.Do(req)
				Expect(err).To(BeNil())
				defer resp.Body.Close()
				Expect(resp.StatusCode).To(Equal(http.StatusNoContent))

				// verify deletion
				resp, err = doRequest(http.MethodGet, proofUrl, nil)
				Expect(err).To(BeNil())
				defer resp.Body.Close()
				Expect(resp.StatusCode).To(Equal(http.StatusNotFound))
			})
		})
	})

	Context("Validation tests", func() {
		When("POST request with space in ClaimUriDigest is sent", func() {
			It("should return 400 with validation error mentioning nospace", func() {
				supports := true
				proof := api.ProofInput{
					ClaimUriDigest: "claim digest",
					ReviewedBy:     "ValidReviewer",
					SupportsClaim:  &supports,
					Uri:            "https://example.com",
				}
				body, err := json.Marshal(proof)
				Expect(err).To(BeNil())

				resp, err := doRequest(http.MethodPost, endpoint, body)
				Expect(err).To(BeNil())
				defer resp.Body.Close()
				Expect(resp.StatusCode).To(Equal(http.StatusBadRequest))

				var errResp map[string]string
				err = json.NewDecoder(resp.Body).Decode(&errResp)
				Expect(err).To(BeNil())
				Expect(strings.ToLower(errResp["error"])).To(ContainSubstring("claimuridigest"))
				Expect(strings.ToLower(errResp["error"])).To(ContainSubstring("nospace"))
			})
		})

		When("POST request with space in ReviewedBy is sent", func() {
			It("should return 400 with validation error mentioning nospace", func() {
				supports := true
				proof := api.ProofInput{
					ClaimUriDigest: "validclaimdigest",
					ReviewedBy:     "Reviewer Name",
					SupportsClaim:  &supports,
					Uri:            "https://example.com",
				}
				body, err := json.Marshal(proof)
				Expect(err).To(BeNil())

				resp, err := doRequest(http.MethodPost, endpoint, body)
				Expect(err).To(BeNil())
				defer resp.Body.Close()
				Expect(resp.StatusCode).To(Equal(http.StatusBadRequest))

				var errResp map[string]string
				err = json.NewDecoder(resp.Body).Decode(&errResp)
				Expect(err).To(BeNil())
				Expect(strings.ToLower(errResp["error"])).To(ContainSubstring("reviewedby"))
				Expect(strings.ToLower(errResp["error"])).To(ContainSubstring("nospace"))
			})
		})

		When("POST request with empty ReviewedBy and ClaimUriDigest is sent", func() {
			It("should return 400 with validation error mentioning nonempty", func() {
				supports := true
				proof := api.ProofInput{
					ClaimUriDigest: "",
					ReviewedBy:     "",
					SupportsClaim:  &supports,
					Uri:            "https://example.com",
				}
				body, err := json.Marshal(proof)
				Expect(err).To(BeNil())

				resp, err := doRequest(http.MethodPost, endpoint, body)
				Expect(err).To(BeNil())
				defer resp.Body.Close()
				Expect(resp.StatusCode).To(Equal(http.StatusBadRequest))

				var errResp map[string]string
				err = json.NewDecoder(resp.Body).Decode(&errResp)
				Expect(err).To(BeNil())
				Expect(strings.ToLower(errResp["error"])).To(ContainSubstring("claimuridigest"))
				Expect(strings.ToLower(errResp["error"])).To(ContainSubstring("reviewedby"))
				Expect(strings.ToLower(errResp["error"])).To(ContainSubstring("required"))
			})
		})

		When("POST request without SupportsClaim is sent", func() {
			It("should return 400 with validation error mentioning supportsClaim", func() {
				proof := api.ProofInput{
					ClaimUriDigest: "validclaimdigest",
					ReviewedBy:     "ValidReviewer",
					Uri:            "https://example.com",
				}
				body, err := json.Marshal(proof)
				Expect(err).To(BeNil())

				resp, err := doRequest(http.MethodPost, endpoint, body)
				Expect(err).To(BeNil())
				defer resp.Body.Close()
				Expect(resp.StatusCode).To(Equal(http.StatusBadRequest))

				var errResp map[string]string
				err = json.NewDecoder(resp.Body).Decode(&errResp)
				Expect(err).To(BeNil())
				Expect(strings.ToLower(errResp["error"])).To(ContainSubstring("supportsclaim"))
				Expect(strings.ToLower(errResp["error"])).To(ContainSubstring("required"))
			})
		})

		When("POST request with invalid url in Uri is sent", func() {
			It("should return 400 with validation error mentioning httpsurl", func() {
				supports := true
				proof := api.ProofInput{
					ClaimUriDigest: "validclaimdigest",
					ReviewedBy:     "ValidReviewer",
					SupportsClaim:  &supports,
					Uri:            "http://not-https.com",
				}
				body, err := json.Marshal(proof)
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

		When("DELETE request is sent for a non-existent proof", func() {
			It("should return 404 with proof not found error", func() {
				proofUrl, err := url.JoinPath(endpoint, "doesnotexist")
				Expect(err).To(BeNil())

				req, err := http.NewRequest(http.MethodDelete, proofUrl, nil)
				Expect(err).To(BeNil())
				addCommonHeaders(req)

				resp, err := client.Do(req)
				Expect(err).To(BeNil())
				defer resp.Body.Close()
				Expect(resp.StatusCode).To(Equal(http.StatusNotFound))

				var errResp map[string]string
				err = json.NewDecoder(resp.Body).Decode(&errResp)
				Expect(err).To(BeNil())
				Expect(strings.ToLower(errResp["error"])).To(ContainSubstring("proof not found"))
			})
		})

		When("PATCH request with empty ReviewedBy is sent", func() {
			It("should return 400 with validation error mentioning nonempty", func() {
				proofUrl, err := url.JoinPath(endpoint, proof1Digest)
				Expect(err).To(BeNil())

				patchBody := api.ProofPatchInput{ReviewedBy: ""}
				body, err := json.Marshal(patchBody)
				Expect(err).To(BeNil())

				req, err := http.NewRequest(http.MethodPatch, proofUrl, bytes.NewBuffer(body))
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
				Expect(strings.ToLower(errResp["error"])).To(ContainSubstring("reviewedby"))
				Expect(strings.ToLower(errResp["error"])).To(ContainSubstring("required"))
			})
		})

		When("PATCH request with space in ReviewedBy value is sent", func() {
			It("should return 400 with validation error mentioning nospace", func() {
				proofUrl, err := url.JoinPath(endpoint, proof1Digest)
				Expect(err).To(BeNil())

				patchBody := api.ProofPatchInput{ReviewedBy: "Reviewer Name"}
				body, err := json.Marshal(patchBody)
				Expect(err).To(BeNil())

				req, err := http.NewRequest(http.MethodPatch, proofUrl, bytes.NewBuffer(body))
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
				Expect(strings.ToLower(errResp["error"])).To(ContainSubstring("reviewedby"))
				Expect(strings.ToLower(errResp["error"])).To(ContainSubstring("nospace"))
			})
		})
	})
})
