package claim_test

import (
	"context"
	"errors"

	"source-score/pkg/api"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"gorm.io/gorm"
)

var _ = Describe("Claim repository layer unit tests", func() {
	Context("Happy path", Ordered, func() {
		When("Posting new claims", func() {
			It("Should create the claims and return their uri digests", func() {
				input := api.ClaimInput{
					SourceUriDigest: sampleClaim1.SourceUriDigest,
					Summary:         &sampleClaim1.Summary,
					Title:           sampleClaim1.Title,
					Uri:             sampleClaim1.Uri,
				}

				digest, err := claimRepo.PostClaim(context.TODO(), &input)
				Expect(err).ToNot(HaveOccurred())
				Expect(digest).ToNot(BeEmpty())
				Expect(digest).To(Equal(claim1Digest))

				input = api.ClaimInput{
					SourceUriDigest: sampleClaim2.SourceUriDigest,
					Summary:         &sampleClaim2.Summary,
					Title:           sampleClaim2.Title,
					Uri:             sampleClaim2.Uri,
				}

				digest, err = claimRepo.PostClaim(context.TODO(), &input)
				Expect(err).ToNot(HaveOccurred())
				Expect(digest).ToNot(BeEmpty())
				Expect(digest).To(Equal(claim2Digest))
			})
		})

		When("Retrieving all claims from the DB", func() {
			It("Should return all claim records from the DB", func() {
				claims, err := claimRepo.GetClaims(context.TODO())
				Expect(err).ToNot(HaveOccurred())
				Expect(len(claims)).To(Equal(2))

				Expect(claims).To(ContainElements(
					sampleClaim1,
					sampleClaim2,
				))
			})
		})

		When("Retrieving a single claim by uri digest", func() {
			It("Should return the matching claim record", func() {
				claim, err := claimRepo.GetClaimByUriDigest(context.TODO(), claim1Digest)
				Expect(err).ToNot(HaveOccurred())
				Expect(claim).ToNot(BeNil())
				Expect(*claim).To(Equal(sampleClaim1))
			})
		})

		When("Patching a claim by its uri digest", func() {
			It("Should update the correct claim record in the DB", func() {
				newSummary := "Updated claim summary"
				newTitle := "Updated Claim Title"
				patchInput := &api.ClaimPatchInput{
					Summary: &newSummary,
					Title:   &newTitle,
				}

				err := claimRepo.PatchClaimByUriDigest(context.TODO(), patchInput, claim1Digest)
				Expect(err).ToNot(HaveOccurred())

				updated, err := claimRepo.GetClaimByUriDigest(context.TODO(), claim1Digest)
				Expect(err).ToNot(HaveOccurred())
				Expect(updated).ToNot(BeNil())
				Expect(updated.Summary).To(Equal(newSummary))
				Expect(updated.Title).To(Equal(newTitle))
			})
		})

		When("Deleting a claim by its uri digest", func() {
			It("Should delete the correct claim record from the DB", func() {
				claim := &api.Claim{
					UriDigest: claim1Digest,
				}

				err := claimRepo.DeleteClaimByUriDigest(context.TODO(), claim)
				Expect(err).ToNot(HaveOccurred())

				_, err = claimRepo.GetClaimByUriDigest(context.TODO(), claim1Digest)
				Expect(err).To(HaveOccurred())
				Expect(errors.Is(err, gorm.ErrRecordNotFound)).To(BeTrue())
			})
		})
	})
})
