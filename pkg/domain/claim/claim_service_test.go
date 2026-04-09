package claim_test

import (
	"context"
	"errors"
	"strings"

	"source-score/pkg/api"
	"source-score/pkg/apperrors"
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
					Summary:         sampleClaim1.Summary,
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
					Summary:         sampleClaim2.Summary,
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

		When("Retrieving a single claim by uri digest", func() {
			It("Should return the matching claim from repository", func() {
				fakeClaimRepo.GetClaimByUriDigestReturnsOnCall(0, &sampleClaim1, nil)

				claim, err := claimSvc.GetClaimByUriDigest(context.TODO(), claim1Digest)
				Expect(err).ToNot(HaveOccurred())
				Expect(claim).ToNot(BeNil())
				Expect(*claim).To(Equal(sampleClaim1))
				Expect(claim.UriDigest).To(Equal(claim1Digest))
				Expect(fakeClaimRepo.GetClaimByUriDigestCallCount()).To(Equal(1))
				_, arg := fakeClaimRepo.GetClaimByUriDigestArgsForCall(0)
				Expect(arg).To(Equal(claim1Digest))
			})
		})

		When("Deleting a claim by its uri digest", func() {
			It("Should delete the correct claim record via repository", func() {
				fakeClaimRepo.GetClaimByUriDigestReturnsOnCall(1, &sampleClaim1, nil)
				fakeClaimRepo.DeleteClaimByUriDigestReturns(nil)

				err := claimSvc.DeleteClaimByUriDigest(context.TODO(), claim1Digest)
				Expect(err).ToNot(HaveOccurred())
				_, digest := fakeClaimRepo.GetClaimByUriDigestArgsForCall(1)
				Expect(digest).To(Equal(claim1Digest))
				_, c := fakeClaimRepo.DeleteClaimByUriDigestArgsForCall(0)
				Expect(*c).To(Equal(sampleClaim1))
			})
		})

		When("Patching a claim by its uri digest", func() {
			It("Should update the claim via repository and return updated record", func() {
				newTitle := "Updated Claim Title"
				newSummary := "Updated claim summary"
				patchInput := api.ClaimPatchInput{
					Title:   &newTitle,
					Summary: &newSummary,
				}

				fakeClaimRepo.PatchClaimByUriDigestReturns(nil)

				err := claimSvc.PatchClaimByUriDigest(context.TODO(), &patchInput, claim1Digest)
				Expect(err).ToNot(HaveOccurred())
				Expect(fakeClaimRepo.PatchClaimByUriDigestCallCount()).To(Equal(1))
				_, argInput, argDigest := fakeClaimRepo.PatchClaimByUriDigestArgsForCall(0)
				Expect(argDigest).To(Equal(claim1Digest))
				Expect(*argInput.Title).To(Equal(newTitle))
				Expect(*argInput.Summary).To(Equal(newSummary))
			})
		})
	})

	Context("Validation tests", func() {
		When("Posting a claim with empty source uri digest", func() {
			It("Should return ErrInvalidClaim", func() {
				input := api.ClaimInput{
					SourceUriDigest: "",
					Summary:         "non-empty",
					Title:           "title",
					Uri:             "https://ok",
				}

				postClaimCalls := fakeClaimRepo.PostClaimCallCount()
				_, err := claimSvc.PostClaim(context.TODO(), &input)
				Expect(err).To(HaveOccurred())
				Expect(errors.Is(err, apperrors.ErrInvalidClaim)).To(BeTrue())
				Expect(strings.Contains(strings.ToLower(err.Error()), "sourceuridigest")).To(BeTrue())
				Expect(strings.Contains(strings.ToLower(err.Error()), "nonempty")).To(BeTrue())
				Expect(fakeClaimRepo.PostClaimCallCount()).To(Equal(postClaimCalls))
			})
		})

		When("Posting a claim with empty title and summary", func() {
			It("Should return ErrInvalidClaim and mention Title and Summary", func() {
				input := api.ClaimInput{
					SourceUriDigest: "srcdigest",
					Summary:         "",
					Title:           "",
					Uri:             "https://ok",
				}

				postClaimCalls := fakeClaimRepo.PostClaimCallCount()
				_, err := claimSvc.PostClaim(context.TODO(), &input)
				Expect(err).To(HaveOccurred())
				Expect(errors.Is(err, apperrors.ErrInvalidClaim)).To(BeTrue())
				Expect(strings.Contains(strings.ToLower(err.Error()), "title")).To(BeTrue())
				Expect(strings.Contains(strings.ToLower(err.Error()), "summary")).To(BeTrue())
				Expect(strings.Contains(strings.ToLower(err.Error()), "nonempty")).To(BeTrue())
				Expect(fakeClaimRepo.PostClaimCallCount()).To(Equal(postClaimCalls))
			})
		})

		When("Posting a claim with a non-https Uri", func() {
			It("Should return ErrInvalidClaim and mention Uri httpsurl", func() {
				input := api.ClaimInput{
					SourceUriDigest: "srcdigest",
					Summary:         "summary",
					Title:           "title",
					Uri:             "http://not-https",
				}

				postClaimCalls := fakeClaimRepo.PostClaimCallCount()
				_, err := claimSvc.PostClaim(context.TODO(), &input)
				Expect(err).To(HaveOccurred())
				Expect(errors.Is(err, apperrors.ErrInvalidClaim)).To(BeTrue())
				Expect(strings.Contains(strings.ToLower(err.Error()), "uri")).To(BeTrue())
				Expect(strings.Contains(strings.ToLower(err.Error()), "httpsurl")).To(BeTrue())
				Expect(fakeClaimRepo.PostClaimCallCount()).To(Equal(postClaimCalls))
			})
		})
	})
})
