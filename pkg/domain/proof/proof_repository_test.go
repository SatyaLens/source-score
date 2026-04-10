package proof_test

import (
	"context"
	"source-score/pkg/api"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Proof repository layer unit tests", Ordered, func() {
	Context("Happy path", Ordered, func() {
		When("Posting new proofs", func() {
			It("Should create the proofs and return their uri digests", func() {
				input := api.ProofInput{
					ClaimUriDigest: sampleProof1.ClaimUriDigest,
					ReviewedBy:     sampleProof1.ReviewedBy,
					SupportsClaim:  sampleProof1.SupportsClaim,
					Uri:            sampleProof1.Uri,
				}

				digest, err := proofRepo.PostProof(context.TODO(), &input)
				Expect(err).ToNot(HaveOccurred())
				Expect(digest).ToNot(BeEmpty())
				Expect(digest).To(Equal(proof1Digest))

				input = api.ProofInput{
					ClaimUriDigest: sampleProof2.ClaimUriDigest,
					ReviewedBy:     sampleProof2.ReviewedBy,
					SupportsClaim:  sampleProof2.SupportsClaim,
					Uri:            sampleProof2.Uri,
				}

				digest, err = proofRepo.PostProof(context.TODO(), &input)
				Expect(err).ToNot(HaveOccurred())
				Expect(digest).ToNot(BeEmpty())
				Expect(digest).To(Equal(proof2Digest))
			})
		})

		When("Retrieving all proofs from the DB", func() {
			It("Should return all proof records from the DB", func() {
				proofs, err := proofRepo.GetProofs(context.TODO())
				Expect(err).ToNot(HaveOccurred())
				Expect(len(proofs)).To(Equal(2))

				Expect(proofs).To(ContainElements(
					sampleProof1,
					sampleProof2,
				))
			})
		})

		When("Retrieving a single proof by uri digest", func() {
			It("Should return the matching proof record", func() {
				p, err := proofRepo.GetProofByUriDigest(context.TODO(), proof1Digest)
				Expect(err).ToNot(HaveOccurred())
				Expect(p).ToNot(BeNil())
				Expect(*p).To(Equal(sampleProof1))
			})
		})

		When("Patching a proof by its uri digest", func() {
			It("Should update the correct proof record in the DB", func() {
				patchInput := &api.ProofPatchInput{
					ReviewedBy: "Updated Reviewer",
				}

				err := proofRepo.PatchProofByUriDigest(context.TODO(), patchInput, proof1Digest)
				Expect(err).ToNot(HaveOccurred())

				updated, err := proofRepo.GetProofByUriDigest(context.TODO(), proof1Digest)
				Expect(err).ToNot(HaveOccurred())
				Expect(updated).ToNot(BeNil())
				Expect(updated.ReviewedBy).To(Equal(patchInput.ReviewedBy))
			})
		})

		When("Deleting a proof by its uri digest", func() {
			It("Should delete the correct proof record from the DB", func() {
				p := &api.Proof{
					UriDigest: proof1Digest,
				}

				err := proofRepo.DeleteProofByUriDigest(context.TODO(), p)
				Expect(err).ToNot(HaveOccurred())

				_, err = proofRepo.GetProofByUriDigest(context.TODO(), proof1Digest)
				Expect(err).To(HaveOccurred())
			})
		})
	})
})
