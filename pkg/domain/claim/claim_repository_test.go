package claim_test

import (
	"context"
	"errors"
	"time"

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
					Summary:         sampleClaim1.Summary,
					Title:           sampleClaim1.Title,
					Uri:             sampleClaim1.Uri,
				}

				digest, err := claimRepo.PostClaim(context.TODO(), &input)
				Expect(err).ToNot(HaveOccurred())
				Expect(digest).ToNot(BeEmpty())
				Expect(digest).To(Equal(claim1Digest))

				input = api.ClaimInput{
					SourceUriDigest: sampleClaim2.SourceUriDigest,
					Summary:         sampleClaim2.Summary,
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

		When("Validating a claim by its uri digest", func() {
			It("Should update the claim's checked and validity fields", func() {
				validity := true
				verification := &api.ClaimVerification{
					Validity: &validity,
				}

				err := claimRepo.VerifyClaimByUriDigest(context.TODO(), verification, claim2Digest)
				Expect(err).ToNot(HaveOccurred())

				validated, err := claimRepo.GetClaimByUriDigest(context.TODO(), claim2Digest)
				Expect(err).ToNot(HaveOccurred())
				Expect(validated).ToNot(BeNil())
				Expect(validated.Checked).To(BeTrue())
				Expect(validated.Validity).To(BeTrue())
			})
		})

		When("Verifying multiple claims at once", func() {
			It("Should update checked and validity fields for all provided claims", func() {
				// Create a new source for this test
				newSource := api.Source{
					Name:      "Test Source for VerifyClaims",
					Score:     0,
					Summary:   "Test source summary",
					Tags:      "test-tag",
					Uri:       "https://test-source-verify",
					UriDigest: "test-source-digest-123",
				}
				result := testDB.Create(&newSource)
				Expect(result.Error).ToNot(HaveOccurred())

				// Create 3 claims associated with the new source
				claim3Input := api.ClaimInput{
					SourceUriDigest: newSource.UriDigest,
					Summary:         "Test claim 3 summary",
					Title:           "Test Claim 3",
					Uri:             "https://test-claim-3",
				}
				digest3, err := claimRepo.PostClaim(context.TODO(), &claim3Input)
				Expect(err).ToNot(HaveOccurred())

				claim4Input := api.ClaimInput{
					SourceUriDigest: newSource.UriDigest,
					Summary:         "Test claim 4 summary",
					Title:           "Test Claim 4",
					Uri:             "https://test-claim-4",
				}
				digest4, err := claimRepo.PostClaim(context.TODO(), &claim4Input)
				Expect(err).ToNot(HaveOccurred())

				claim5Input := api.ClaimInput{
					SourceUriDigest: newSource.UriDigest,
					Summary:         "Test claim 5 summary",
					Title:           "Test Claim 5",
					Uri:             "https://test-claim-5",
				}
				digest5, err := claimRepo.PostClaim(context.TODO(), &claim5Input)
				Expect(err).ToNot(HaveOccurred())

				// Prepare updated claims - 2 claims with different validity values
				updatedClaims := []api.Claim{
					{UriDigest: digest3, Validity: true},
					{UriDigest: digest4, Validity: false},
				}

				// Call VerifyClaims
				err = claimRepo.VerifyAllClaims(context.TODO(), updatedClaims)
				Expect(err).ToNot(HaveOccurred())

				// Verify claim 3 - should be checked=true, validity=true
				claim3, err := claimRepo.GetClaimByUriDigest(context.TODO(), digest3)
				Expect(err).ToNot(HaveOccurred())
				Expect(claim3.Checked).To(BeTrue())
				Expect(claim3.Validity).To(BeTrue())
				Expect(claim3.Title).To(Equal(claim3Input.Title))
				Expect(claim3.Summary).To(Equal(claim3Input.Summary))

				// Verify claim 4 - should be checked=true, validity=false
				claim4, err := claimRepo.GetClaimByUriDigest(context.TODO(), digest4)
				Expect(err).ToNot(HaveOccurred())
				Expect(claim4.Checked).To(BeTrue())
				Expect(claim4.Validity).To(BeFalse())
				Expect(claim4.Title).To(Equal(claim4Input.Title))
				Expect(claim4.Summary).To(Equal(claim4Input.Summary))

				// Verify claim 5 - should remain unchanged (checked=false, validity=false)
				claim5, err := claimRepo.GetClaimByUriDigest(context.TODO(), digest5)
				Expect(err).ToNot(HaveOccurred())
				Expect(claim5.Checked).To(BeFalse())
				Expect(claim5.Validity).To(BeFalse())
				Expect(claim5.Title).To(Equal(claim5Input.Title))
				Expect(claim5.Summary).To(Equal(claim5Input.Summary))
			})
		})

		When("Getting checked claims grouped by sources", func() {
			It("Should return only checked claims grouped by source uri digest", func() {
				// Create a new source
				newSource := api.Source{
					Name:      "Test Source for GetCheckedClaims",
					Score:     0,
					Summary:   "Test source summary",
					Tags:      "test-tag",
					Uri:       "https://test-source-checked",
					UriDigest: "test-source-digest-456",
				}
				result := testDB.Create(&newSource)
				Expect(result.Error).ToNot(HaveOccurred())

				// Create 2 claims
				claim6Input := api.ClaimInput{
					SourceUriDigest: newSource.UriDigest,
					Summary:         "Test claim 6 summary",
					Title:           "Test Claim 6",
					Uri:             "https://test-claim-6",
				}
				digest6, err := claimRepo.PostClaim(context.TODO(), &claim6Input)
				Expect(err).ToNot(HaveOccurred())

				claim7Input := api.ClaimInput{
					SourceUriDigest: newSource.UriDigest,
					Summary:         "Test claim 7 summary",
					Title:           "Test Claim 7",
					Uri:             "https://test-claim-7",
				}
				digest7, err := claimRepo.PostClaim(context.TODO(), &claim7Input)
				Expect(err).ToNot(HaveOccurred())

				// Create 2 proofs for claim 6 (2 supporting, 0 refuting)
				proof1 := api.Proof{
					ClaimUriDigest: digest6,
					ReviewedBy:     "ReviewerX",
					SupportsClaim:  true,
					Uri:            "https://proof-1-for-claim6",
					UriDigest:      "proof1digest",
				}
				result = testDB.Create(&proof1)
				Expect(result.Error).ToNot(HaveOccurred())

				proof2 := api.Proof{
					ClaimUriDigest: digest6,
					ReviewedBy:     "ReviewerY",
					SupportsClaim:  true,
					Uri:            "https://proof-2-for-claim6",
					UriDigest:      "proof2digest",
				}
				result = testDB.Create(&proof2)
				Expect(result.Error).ToNot(HaveOccurred())

				// Call VerifyClaims to mark claim 6 as checked
				updatedClaims := []api.Claim{
					{UriDigest: digest6, Validity: true},
				}
				err = claimRepo.VerifyAllClaims(context.TODO(), updatedClaims)
				Expect(err).ToNot(HaveOccurred())

				// Eventually call GetCheckedClaimsBySources and verify claim 6 is present
				Eventually(func() bool {
					srcsClaims, err := claimRepo.GetCheckedClaimsBySources(context.TODO())
					if err != nil {
						return false
					}

					// Check if the source digest exists in the map
					claims, exists := srcsClaims[newSource.UriDigest]
					if !exists {
						return false
					}

					// Check if claim 6 is in the list
					for _, claim := range claims {
						if claim.UriDigest == digest6 && claim.Checked == true {
							return true
						}
					}
					return false
				}, 10*time.Second, 1*time.Second).Should(BeTrue())

				// Verify claim 7 is NOT in the checked claims (it was not verified)
				srcsClaims, err := claimRepo.GetCheckedClaimsBySources(context.TODO())
				Expect(err).ToNot(HaveOccurred())

				claims := srcsClaims[newSource.UriDigest]
				provingClaim, err := claimRepo.GetClaimByUriDigest(context.TODO(), digest6)
				Expect(err).ToNot(HaveOccurred())
				uncheckedClaim, err := claimRepo.GetClaimByUriDigest(context.TODO(), digest7)
				Expect(err).ToNot(HaveOccurred())
				Expect(claims).To(ContainElement(*provingClaim))
				Expect(claims).ToNot(ContainElement(*uncheckedClaim))
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

	Context("Validation tests", func() {
		When("Retrieving a non-existent claim by uri digest", func() {
			It("Should return gorm.ErrRecordNotFound", func() {
				_, err := claimRepo.GetClaimByUriDigest(context.TODO(), "doesnotexist")
				Expect(err).To(HaveOccurred())
				Expect(errors.Is(err, gorm.ErrRecordNotFound)).To(BeTrue())
			})
		})

		When("Patching a non-existent claim by uri digest", func() {
			It("Should return gorm.ErrRecordNotFound", func() {
				newTitle := "New Title"
				patchInput := &api.ClaimPatchInput{Title: &newTitle}

				err := claimRepo.PatchClaimByUriDigest(context.TODO(), patchInput, "doesnotexist")
				Expect(err).To(HaveOccurred())
				Expect(errors.Is(err, gorm.ErrRecordNotFound)).To(BeTrue())
			})
		})

		When("Validating a non-existent claim by uri digest", func() {
			It("Should return gorm.ErrRecordNotFound", func() {
				validity := true
				verification := &api.ClaimVerification{Validity: &validity}

				err := claimRepo.VerifyClaimByUriDigest(context.TODO(), verification, "doesnotexist")
				Expect(err).To(HaveOccurred())
				Expect(errors.Is(err, gorm.ErrRecordNotFound)).To(BeTrue())
			})
		})
	})
})
