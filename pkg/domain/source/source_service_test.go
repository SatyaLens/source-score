package source_test

import (
	"context"
	"errors"
	"source-score/pkg/api"
	"source-score/pkg/apperrors"
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"gorm.io/gorm"
)

var _ = Describe("Source model service layer unit test", Ordered, func() {
	Context("Happy path", Ordered, func() {
		When("Adding a new source with valid input", func() {
			It("Should pass the data to the repository layer", func() {
				fakeSourceRepo.PostSourceReturnsOnCall(0, uriDigest1, nil)
				digest, err := sourceSvc.PostSource(context.TODO(), &sampleSourceInput1)
				Expect(err).ToNot(HaveOccurred())
				Expect(digest).To(Equal(uriDigest1))

				Expect(fakeSourceRepo.PostSourceCallCount()).To(Equal(1))
				_, srcInput := fakeSourceRepo.PostSourceArgsForCall(0)
				Expect(srcInput).To(Equal(&sampleSourceInput1))
			})
		})

		When("Retrieving all sources", func() {
			It("Should pass the request to the repository layer", func() {
				expectedSources := []api.Source{sampleSource1, sampleSource2}
				fakeSourceRepo.GetSourcesReturnsOnCall(0, expectedSources, nil)
				sources, err := sourceSvc.GetSources(context.TODO())
				Expect(err).ToNot(HaveOccurred())
				Expect(len(sources)).To(Equal(2))
				Expect(sources).To(Equal(expectedSources))
				Expect(fakeSourceRepo.GetSourcesCallCount()).To(Equal(1))
			})
		})

		When("Retrieving a source by its uri digest", func() {
			It("Should pass the digest to the repo layer", func() {
				fakeSourceRepo.GetSourceByUriDigestReturnsOnCall(0, &sampleSource1, nil)
				src, err := sourceSvc.GetSourceByUriDigest(context.TODO(), uriDigest1)
				Expect(err).ToNot(HaveOccurred())
				Expect(src).To(Equal(&sampleSource1))
				Expect(fakeSourceRepo.GetSourceByUriDigestCallCount()).To(Equal(1))
				_, digest := fakeSourceRepo.GetSourceByUriDigestArgsForCall(0)
				Expect(digest).To(Equal(uriDigest1))
			})
		})

		When("Updating a source by its uri digest", func() {
			It("Should update the correct source record in the DB", func() {
				name := "Updated Sample Source 1"
				summary := "Updated Sample summary"
				tags := "updated-tag1"
				sourceInput := &api.SourcePatchInput{
					Name:    &name,
					Summary: &summary,
					Tags:    &tags,
				}
				updatedSource = sampleSource1
				updatedSource.Name = "Updated Sample Source 1"
				updatedSource.Summary = "Updated Sample summary"
				updatedSource.Tags = "updated-tag1"
				fakeSourceRepo.GetSourceByUriDigestReturnsOnCall(1, &updatedSource, nil)

				err := sourceSvc.PatchSourceByUriDigest(context.TODO(), sourceInput, uriDigest1)
				Expect(err).ToNot(HaveOccurred())
				_, srcInput := fakeSourceRepo.PostSourceArgsForCall(0)
				Expect(srcInput.Name).To(Equal(sampleSource1.Name))
				Expect(srcInput.Summary).To(Equal(sampleSource1.Summary))
				Expect(srcInput.Tags).To(Equal(sampleSource1.Tags))

				source, err := sourceSvc.GetSourceByUriDigest(context.TODO(), uriDigest1)
				Expect(err).ToNot(HaveOccurred())
				_, digest := fakeSourceRepo.GetSourceByUriDigestArgsForCall(1)
				Expect(digest).To(Equal(uriDigest1))
				Expect(source.Name).To(BeEquivalentTo(*sourceInput.Name))
				Expect(source.Summary).To(BeEquivalentTo(*sourceInput.Summary))
				Expect(source.Tags).To(BeEquivalentTo(*sourceInput.Tags))
				Expect(source.Uri).To(BeEquivalentTo(sampleSourceInput1.Uri))
				Expect(source.UriDigest).To(BeEquivalentTo(uriDigest1))
			})
		})

		When("Deleting a source by its uri digest", func() {
			It("Should delete the correct source record from the DB", func() {
				fakeSourceRepo.GetSourceByUriDigestReturnsOnCall(2, &updatedSource, nil)

				err := sourceSvc.DeleteSourceByUriDigest(context.TODO(), uriDigest1)
				Expect(err).ToNot(HaveOccurred())
				_, digest := fakeSourceRepo.GetSourceByUriDigestArgsForCall(2)
				Expect(digest).To(Equal(uriDigest1))
				_, src := fakeSourceRepo.DeleteSourceByUriDigestArgsForCall(0)
				Expect(*src).To(Equal(updatedSource))
			})
		})
	})

	Context("Source POST validation tests", func() {
		When("Creating a source with tags containing spaces", func() {
			It("Should return invalid source error with nospace validation message", func() {
				invalidInput := &api.SourceInput{
					Name:    "Valid Name",
					Summary: "Valid Summary",
					Tags:    "tag1, tag2",
					Uri:     "https://example.com",
				}

				_, err := sourceSvc.PostSource(context.TODO(), invalidInput)

				Expect(err).ToNot(BeNil())
				Expect(errors.Is(err, apperrors.ErrInvalidSource)).To(BeTrue())
				Expect(strings.ToLower(err.Error())).To(ContainSubstring("tags validation failed"))
				Expect(strings.ToLower(err.Error())).To(ContainSubstring("nospace"))
			})
		})

		When("Creating a source with non-https uri", func() {
			It("Should return invalid source error with httpsurl validation message", func() {
				invalidInput := &api.SourceInput{
					Name:    "Valid Name",
					Summary: "Valid Summary",
					Tags:    "tag1,tag2",
					Uri:     "http://example.com",
				}

				_, err := sourceSvc.PostSource(context.TODO(), invalidInput)

				Expect(err).ToNot(BeNil())
				Expect(errors.Is(err, apperrors.ErrInvalidSource)).To(BeTrue())
				Expect(strings.ToLower(err.Error())).To(ContainSubstring("uri validation failed"))
				Expect(strings.ToLower(err.Error())).To(ContainSubstring("httpsurl"))
			})
		})

		When("Creating a source with empty summary and name", func() {
			It("Should return invalid source error with httpsurl validation message", func() {
				invalidInput := &api.SourceInput{
					Name:    "",
					Summary: "",
					Tags:    "tag1,tag2",
					Uri:     "https://example.com",
				}

				_, err := sourceSvc.PostSource(context.TODO(), invalidInput)

				Expect(err).ToNot(BeNil())
				Expect(errors.Is(err, apperrors.ErrInvalidSource)).To(BeTrue())
				Expect(strings.ToLower(err.Error())).To(ContainSubstring("name validation failed"))
				Expect(strings.ToLower(err.Error())).To(ContainSubstring("summary validation failed"))
				Expect(strings.ToLower(err.Error())).To(ContainSubstring("nonempty"))
				Expect(strings.ToLower(err.Error())).ToNot(ContainSubstring("tags"))
				Expect(strings.ToLower(err.Error())).ToNot(ContainSubstring("uri"))
			})
		})
	})

	Context("Source PACTH validation tests", func() {
		When("Patching a source that does not exist", func() {
			It("Should return source not found error", func() {
				fakeSourceRepo.PatchSourceByUriDigestReturns(gorm.ErrRecordNotFound)

				name := "Updated Sample Source 1"
				summary := "Updated Sample summary"
				tags := "updated-tag1"
				sourceInput := &api.SourcePatchInput{
					Name:    &name,
					Summary: &summary,
					Tags:    &tags,
				}
				err := sourceSvc.PatchSourceByUriDigest(context.TODO(), sourceInput, "invalid-digest")

				Expect(err).ToNot(BeNil())
				Expect(errors.Is(err, apperrors.ErrSourceNotFound)).To(BeTrue())
			})
		})

		When("Patching a source with empty name", func() {
			It("Should return invalid source error with nonempty validation message", func() {
				emptyName := ""
				invalidInput := &api.SourcePatchInput{
					Name: &emptyName,
				}

				err := sourceSvc.PatchSourceByUriDigest(context.TODO(), invalidInput, uriDigest1)

				Expect(err).ToNot(BeNil())
				Expect(errors.Is(err, apperrors.ErrInvalidSource)).To(BeTrue())
				Expect(strings.ToLower(err.Error())).To(ContainSubstring("name validation failed"))
				Expect(strings.ToLower(err.Error())).To(ContainSubstring("nonempty"))
			})
		})

		When("Patching a source with tags containing spaces", func() {
			It("Should return invalid source error with nospace validation message", func() {
				invalidTags := "tag1, tag2"
				invalidInput := &api.SourcePatchInput{
					Tags: &invalidTags,
				}

				err := sourceSvc.PatchSourceByUriDigest(context.TODO(), invalidInput, uriDigest1)

				Expect(err).ToNot(BeNil())
				Expect(errors.Is(err, apperrors.ErrInvalidSource)).To(BeTrue())
				Expect(err.Error()).To(ContainSubstring("Tags"))
				Expect(err.Error()).To(ContainSubstring("nospace"))
			})
		})

		When("Patching a source with empty tag and summary", func() {
			It("Should return invalid source error with nospace validation message", func() {
				validName := "valid-name"
				invalidSummary := ""
				invalidTags := ""
				invalidInput := &api.SourcePatchInput{
					Name:    &validName,
					Summary: &invalidSummary,
					Tags:    &invalidTags,
				}

				err := sourceSvc.PatchSourceByUriDigest(context.TODO(), invalidInput, uriDigest1)
				Expect(err).ToNot(BeNil())
				Expect(errors.Is(err, apperrors.ErrInvalidSource)).To(BeTrue())
				Expect(strings.ToLower(err.Error())).To(ContainSubstring("summary validation failed"))
				Expect(strings.ToLower(err.Error())).To(ContainSubstring("tags validation failed"))
				Expect(strings.ToLower(err.Error())).To(ContainSubstring("nonempty"))
			})
		})
	})
})
