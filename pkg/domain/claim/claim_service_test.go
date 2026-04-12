package claim_test

import (
	"context"
	"errors"
	"strings"

	"source-score/pkg/api"
	"source-score/pkg/apperrors"
	"source-score/pkg/domain/claim"
	"source-score/pkg/domain/claim/claimfakes"

	"gorm.io/gorm"

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

			It("Should patch only the summary when title is nil", func() {
				newSummary := "Summary only update"
				patchInput := api.ClaimPatchInput{
					Summary: &newSummary,
					Title:   nil,
				}

				before := fakeClaimRepo.PatchClaimByUriDigestCallCount()
				fakeClaimRepo.PatchClaimByUriDigestReturns(nil)
				err := claimSvc.PatchClaimByUriDigest(context.TODO(), &patchInput, claim1Digest)
				Expect(err).ToNot(HaveOccurred())
				Expect(fakeClaimRepo.PatchClaimByUriDigestCallCount()).To(Equal(before + 1))
				_, argInput, argDigest := fakeClaimRepo.PatchClaimByUriDigestArgsForCall(before)
				Expect(argDigest).To(Equal(claim1Digest))
				Expect(argInput.Title).To(BeNil())
				Expect(*argInput.Summary).To(Equal(newSummary))
			})

			It("Should patch only the title when summary is nil", func() {
				newTitle := "Title only update"
				patchInput := api.ClaimPatchInput{
					Summary: nil,
					Title:   &newTitle,
				}

				before := fakeClaimRepo.PatchClaimByUriDigestCallCount()
				fakeClaimRepo.PatchClaimByUriDigestReturns(nil)
				err := claimSvc.PatchClaimByUriDigest(context.TODO(), &patchInput, claim1Digest)
				Expect(err).ToNot(HaveOccurred())
				Expect(fakeClaimRepo.PatchClaimByUriDigestCallCount()).To(Equal(before + 1))
				_, argInput, argDigest := fakeClaimRepo.PatchClaimByUriDigestArgsForCall(before)
				Expect(argDigest).To(Equal(claim1Digest))
				Expect(*argInput.Title).To(Equal(newTitle))
				Expect(argInput.Summary).To(BeNil())
			})

		})

		When("Validating a claim by its uri digest", func() {
			It("Should validate the claim via repository with validity false", func() {
				validity := false
				verification := &api.ClaimVerification{
					Validity: &validity,
				}

				before := fakeClaimRepo.VerifyClaimByUriDigestCallCount()
				fakeClaimRepo.VerifyClaimByUriDigestReturns(nil)

				err := claimSvc.VerifyClaimByUriDigest(context.TODO(), verification, claim2Digest)
				Expect(err).ToNot(HaveOccurred())
				Expect(fakeClaimRepo.VerifyClaimByUriDigestCallCount()).To(Equal(before + 1))
				_, argVerification, argDigest := fakeClaimRepo.VerifyClaimByUriDigestArgsForCall(before)
				Expect(argDigest).To(Equal(claim2Digest))
				Expect(argVerification.Validity).ToNot(BeNil())
				Expect(*argVerification.Validity).To(BeFalse())
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

		When("Deleting a non-existent claim by uri digest", func() {
			It("Should return ErrClaimNotFound and not call Delete on repo", func() {
				fakeClaimRepo.GetClaimByUriDigestReturns(nil, gorm.ErrRecordNotFound)
				before := fakeClaimRepo.DeleteClaimByUriDigestCallCount()

				err := claimSvc.DeleteClaimByUriDigest(context.TODO(), "doesnotexist")
				Expect(err).To(HaveOccurred())
				Expect(errors.Is(err, apperrors.ErrClaimNotFound)).To(BeTrue())
				Expect(fakeClaimRepo.DeleteClaimByUriDigestCallCount()).To(Equal(before))
			})
		})

		When("Getting a non-existent claim by uri digest", func() {
			It("Should return ErrClaimNotFound and nil claim", func() {
				fakeClaimRepo.GetClaimByUriDigestReturns(nil, gorm.ErrRecordNotFound)

				c, err := claimSvc.GetClaimByUriDigest(context.TODO(), "doesnotexist")
				Expect(err).To(HaveOccurred())
				Expect(c).To(BeNil())
				Expect(errors.Is(err, apperrors.ErrClaimNotFound)).To(BeTrue())
			})
		})

		When("Patching a non-existent claim by uri digest", func() {
			It("Should return ErrClaimNotFound", func() {
				fakeClaimRepo.PatchClaimByUriDigestReturns(gorm.ErrRecordNotFound)

				patchInput := api.ClaimPatchInput{Title: nil, Summary: nil}
				err := claimSvc.PatchClaimByUriDigest(context.TODO(), &patchInput, "doesnotexist")
				Expect(err).To(HaveOccurred())
				Expect(errors.Is(err, apperrors.ErrClaimNotFound)).To(BeTrue())
			})
		})

		When("Patching a claim with empty title and summary", func() {
			It("Should return ErrInvalidClaim and not call repo", func() {
				empty := ""
				patchInput := api.ClaimPatchInput{Title: &empty, Summary: &empty}
				before := fakeClaimRepo.PatchClaimByUriDigestCallCount()

				err := claimSvc.PatchClaimByUriDigest(context.TODO(), &patchInput, claim1Digest)
				Expect(err).To(HaveOccurred())
				Expect(errors.Is(err, apperrors.ErrInvalidClaim)).To(BeTrue())
				Expect(strings.Contains(strings.ToLower((err.Error())), "title")).To(BeTrue())
				Expect(strings.Contains(strings.ToLower(err.Error()), "summary")).To(BeTrue())
				Expect(strings.Contains(strings.ToLower(err.Error()), "nonempty")).To(BeTrue())
				Expect(fakeClaimRepo.PatchClaimByUriDigestCallCount()).To(Equal(before))
			})
		})

		When("Validating a claim without Validity field", func() {
			It("Should return ErrInvalidClaimVerification and not call repo", func() {
				verification := &api.ClaimVerification{
					Validity: nil,
				}
				before := fakeClaimRepo.VerifyClaimByUriDigestCallCount()

				err := claimSvc.VerifyClaimByUriDigest(context.TODO(), verification, claim1Digest)
				Expect(err).To(HaveOccurred())
				Expect(errors.Is(err, apperrors.ErrInvalidClaimVerification)).To(BeTrue())
				Expect(strings.Contains(strings.ToLower(err.Error()), "validity")).To(BeTrue())
				Expect(strings.Contains(strings.ToLower(err.Error()), "nonempty")).To(BeTrue())
				Expect(fakeClaimRepo.VerifyClaimByUriDigestCallCount()).To(Equal(before))
			})
		})

		When("Validating a non-existent claim by uri digest", func() {
			It("Should return ErrClaimNotFound", func() {
				validity := true
				verification := &api.ClaimVerification{Validity: &validity}
				before := fakeClaimRepo.VerifyClaimByUriDigestCallCount()
				fakeClaimRepo.VerifyClaimByUriDigestReturns(gorm.ErrRecordNotFound)

				err := claimSvc.VerifyClaimByUriDigest(context.TODO(), verification, "doesnotexist")
				Expect(err).To(HaveOccurred())
				Expect(errors.Is(err, apperrors.ErrClaimNotFound)).To(BeTrue())
				Expect(fakeClaimRepo.VerifyClaimByUriDigestCallCount()).To(Equal(before + 1))
			})
		})
	})
})
