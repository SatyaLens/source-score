package claim_test

import (
	"context"

	"source-score/pkg/api"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
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
	})
})
