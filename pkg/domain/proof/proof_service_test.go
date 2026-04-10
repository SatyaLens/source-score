package proof_test

import (
	"context"

	"source-score/pkg/api"
	"source-score/pkg/domain/proof"
	"source-score/pkg/domain/proof/prooffakes"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var (
	fakeProofRepo = prooffakes.FakeProofRepository{}
	proofSvc      = proof.NewProofService(context.TODO(), &fakeProofRepo)
)

var _ = Describe("Proof model service layer unit tests", Ordered, func() {
	Context("Happy path", func() {
		When("Posting new proofs", func() {
			It("Should pass data to the repository and return digest", func() {
				fakeProofRepo.PostProofReturnsOnCall(0, proof1Digest, nil)
				fakeProofRepo.PostProofReturnsOnCall(1, proof2Digest, nil)

				input := api.ProofInput{
					ClaimUriDigest: sampleProof1.ClaimUriDigest,
					ReviewedBy:     sampleProof1.ReviewedBy,
					SupportsClaim:  sampleProof1.SupportsClaim,
					Uri:            sampleProof1.Uri,
				}

				digest, err := proofSvc.PostProof(context.TODO(), &input)
				Expect(err).ToNot(HaveOccurred())
				Expect(digest).To(Equal(proof1Digest))
				Expect(fakeProofRepo.PostProofCallCount()).To(Equal(1))
				_, arg := fakeProofRepo.PostProofArgsForCall(0)
				Expect(arg).To(Equal(&input))

				input = api.ProofInput{
					ClaimUriDigest: sampleProof2.ClaimUriDigest,
					ReviewedBy:     sampleProof2.ReviewedBy,
					SupportsClaim:  sampleProof2.SupportsClaim,
					Uri:            sampleProof2.Uri,
				}

				digest, err = proofSvc.PostProof(context.TODO(), &input)
				Expect(err).ToNot(HaveOccurred())
				Expect(digest).To(Equal(proof2Digest))
				Expect(fakeProofRepo.PostProofCallCount()).To(Equal(2))
				_, arg = fakeProofRepo.PostProofArgsForCall(1)
				Expect(arg).To(Equal(&input))
			})
		})

		When("Retrieving all proofs", func() {
			It("Should return proofs from repository", func() {
				expected := []api.Proof{sampleProof1, sampleProof2}
				fakeProofRepo.GetProofsReturnsOnCall(0, expected, nil)

				proofs, err := proofSvc.GetProofs(context.TODO())
				Expect(err).ToNot(HaveOccurred())
				Expect(len(proofs)).To(Equal(2))
				Expect(proofs).To(ContainElements(expected))
				Expect(fakeProofRepo.GetProofsCallCount()).To(Equal(1))
			})
		})

		When("Retrieving a single proof by uri digest", func() {
			It("Should return the matching proof from repository", func() {
				fakeProofRepo.GetProofByUriDigestReturnsOnCall(0, &sampleProof1, nil)

				p, err := proofSvc.GetProofByUriDigest(context.TODO(), proof1Digest)
				Expect(err).ToNot(HaveOccurred())
				Expect(p).ToNot(BeNil())
				Expect(*p).To(Equal(sampleProof1))
				Expect(p.UriDigest).To(Equal(proof1Digest))
				Expect(fakeProofRepo.GetProofByUriDigestCallCount()).To(Equal(1))
				_, arg := fakeProofRepo.GetProofByUriDigestArgsForCall(0)
				Expect(arg).To(Equal(proof1Digest))
			})
		})		

		When("Patching a proof by its uri digest", func() {
			It("Should update the proof via repository", func() {
				patchInput := api.ProofPatchInput{ReviewedBy: "UpdatedReviewer"}
				fakeProofRepo.PatchProofByUriDigestReturns(nil)

				err := proofSvc.PatchProofByUriDigest(context.TODO(), &patchInput, proof1Digest)
				Expect(err).ToNot(HaveOccurred())
				Expect(fakeProofRepo.PatchProofByUriDigestCallCount()).To(Equal(1))
				_, argInput, argDigest := fakeProofRepo.PatchProofByUriDigestArgsForCall(0)
				Expect(argDigest).To(Equal(proof1Digest))
				Expect(argInput.ReviewedBy).To(Equal(patchInput.ReviewedBy))
			})
		})

		When("Deleting a proof by its uri digest", func() {
			It("Should delete the correct proof record via repository", func() {
				fakeProofRepo.GetProofByUriDigestReturnsOnCall(1, &sampleProof1, nil)
				fakeProofRepo.DeleteProofByUriDigestReturns(nil)

				err := proofSvc.DeleteProofByUriDigest(context.TODO(), proof1Digest)
				Expect(err).ToNot(HaveOccurred())
				_, digest := fakeProofRepo.GetProofByUriDigestArgsForCall(1)
				Expect(digest).To(Equal(proof1Digest))
				_, c := fakeProofRepo.DeleteProofByUriDigestArgsForCall(0)
				Expect(*c).To(Equal(sampleProof1))
			})
		})
	})
})
