package claim_test

import (
	"context"
	"source-score/pkg/api"
	"source-score/pkg/domain/claim"
	"source-score/pkg/domain/claim/claimfakes"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var (
	fakeClaimRepo = claimfakes.FakeClaimRepository{}
	claimSvc      = claim.NewClaimService(context.TODO(), &fakeClaimRepo)
)

var _ = Describe("Claim model service layer unit tests", Ordered, func() {
	Context("Happy path", func() {
		When("Posting a new claim with valid input", func() {
			It("Should pass data to the repository and return digest", func() {
				fakeClaimRepo.PostClaimReturnsOnCall(0, claim1Digest, nil)
				fakeClaimRepo.PostClaimReturnsOnCall(1, claim2Digest, nil)

				input := api.ClaimInput{
					SourceUriDigest: sampleClaim1.SourceUriDigest,
					Summary:         &sampleClaim1.Summary,
					Title:           sampleClaim1.Title,
					Uri:             sampleClaim1.Uri,
				}
				digest, err := claimSvc.PostClaim(context.TODO(), &input)
				Expect(err).ToNot(HaveOccurred())
				Expect(digest).To(Equal(claim1Digest))
				Expect(fakeClaimRepo.PostClaimCallCount()).To(Equal(1))
				_, arg := fakeClaimRepo.PostClaimArgsForCall(0)
				Expect(arg).To(Equal(&input))

				input = api.ClaimInput{
					SourceUriDigest: sampleClaim2.SourceUriDigest,
					Summary:         &sampleClaim2.Summary,
					Title:           sampleClaim2.Title,
					Uri:             sampleClaim2.Uri,
				}
				digest, err = claimSvc.PostClaim(context.TODO(), &input)
				Expect(err).ToNot(HaveOccurred())
				Expect(digest).To(Equal(claim2Digest))
				Expect(fakeClaimRepo.PostClaimCallCount()).To(Equal(2))
				_, arg = fakeClaimRepo.PostClaimArgsForCall(1)
				Expect(arg).To(Equal(&input))
			})
		})

		When("Retrieving all claims", func() {
			It("Should return claims from repository", func() {
				expected := []api.Claim{sampleClaim1, sampleClaim2}
				fakeClaimRepo.GetClaimsReturnsOnCall(0, expected, nil)

				claims, err := claimSvc.GetClaims(context.TODO())
				Expect(err).ToNot(HaveOccurred())
				Expect(len(claims)).To(Equal(2))
				Expect(claims).To(ContainElements(expected))
				Expect(fakeClaimRepo.GetClaimsCallCount()).To(Equal(1))
			})
		})
	})
})
