package source_test

import (
	"context"
	"errors"
	"source-score/pkg/api"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"gorm.io/gorm"
)

var _ = Describe("Source model repository layer unit tests", Ordered, func() {
	Context("Happy path", Ordered, func() {
		When("Adding new sources to the DB with valid input", func() {
			It("Should create the source record in the DB", func() {
				digest, err := sourceRepo.PostSource(context.TODO(), &sampleSourceInput1)
				Expect(err).ToNot(HaveOccurred())
				Expect(digest).To(Equal(uriDigest1))

				digest, err = sourceRepo.PostSource(context.TODO(), &sampleSourceInput2)
				Expect(err).ToNot(HaveOccurred())
				Expect(digest).To(Equal(uriDigest2))
			})
		})

		When("Retrieving all sources from the DB", func() {
			It("Should return all source records from the DB", func() {
				sources, err := sourceRepo.GetSources(context.TODO())
				Expect(err).ToNot(HaveOccurred())
				Expect(len(sources)).To(Equal(2))

				Expect(sources).To(ContainElements(
					sampleSource1,
					sampleSource2,
				))
			})
		})

		When("Retrieving a source by its uri digest", func() {
			It("Should return the correct source record from the DB", func() {
				source, err := sourceRepo.GetSourceByUriDigest(context.TODO(), uriDigest1)
				Expect(err).ToNot(HaveOccurred())
				Expect(source.Name).To(BeEquivalentTo(sampleSourceInput1.Name))
				Expect(source.Summary).To(BeEquivalentTo(sampleSourceInput1.Summary))
				Expect(source.Tags).To(BeEquivalentTo(sampleSourceInput1.Tags))
				Expect(source.Uri).To(BeEquivalentTo(sampleSourceInput1.Uri))
				Expect(source.UriDigest).To(BeEquivalentTo(uriDigest1))
			})
		})

		When("Updating all the fields of source by its uri digest", func() {
			It("Should update the correct source record in the DB", func() {
				name := "Updated Sample Source 1"
				summary := "Updated Sample summary"
				tags := "updated-tag1"
				sourceInput := &api.SourcePatchInput{
					Name:    &name,
					Summary: &summary,
					Tags:    &tags,
				}

				err := sourceRepo.PatchSourceByUriDigest(context.TODO(), sourceInput, uriDigest1)
				Expect(err).ToNot(HaveOccurred())

				source, err := sourceRepo.GetSourceByUriDigest(context.TODO(), uriDigest1)
				Expect(err).ToNot(HaveOccurred())
				Expect(source.Name).To(BeEquivalentTo(*sourceInput.Name))
				Expect(source.Summary).To(BeEquivalentTo(*sourceInput.Summary))
				Expect(source.Tags).To(BeEquivalentTo(*sourceInput.Tags))
				Expect(source.Uri).To(BeEquivalentTo(sampleSourceInput1.Uri))
				Expect(source.UriDigest).To(BeEquivalentTo(uriDigest1))
			})
		})

		When("Updating some fields of source by its uri digest", func() {
			It("Should update the correct source record in the DB", func() {
				name := "Twice Updated Sample Source 1"
				tags := "twice-updated-tag1"
				sourceInput := &api.SourcePatchInput{
					Name: &name,
					Tags: &tags,
				}

				err := sourceRepo.PatchSourceByUriDigest(context.TODO(), sourceInput, uriDigest1)
				Expect(err).ToNot(HaveOccurred())

				source, err := sourceRepo.GetSourceByUriDigest(context.TODO(), uriDigest1)
				Expect(err).ToNot(HaveOccurred())
				Expect(source.Name).To(BeEquivalentTo(*sourceInput.Name))
				Expect(source.Summary).To(BeEquivalentTo("Updated Sample summary"))
				Expect(source.Tags).To(BeEquivalentTo(*sourceInput.Tags))
				Expect(source.Uri).To(BeEquivalentTo(sampleSourceInput1.Uri))
				Expect(source.UriDigest).To(BeEquivalentTo(uriDigest1))
			})
		})

		When("Deleting a source by its uri digest", func() {
			It("Should delete the correct source record from the DB", func() {
				source := &api.Source{
					UriDigest: uriDigest1,
				}

				err := sourceRepo.DeleteSourceByUriDigest(context.TODO(), source)
				Expect(err).ToNot(HaveOccurred())

				_, err = sourceRepo.GetSourceByUriDigest(context.TODO(), uriDigest1)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("record not found"))
			})
		})

		When("Updating scores for multiple sources", func() {
			It("Should update only the score field for all provided sources", func() {
				// Create 3 new sources
				source5Input := api.SourceInput{
					Name:    "Sample Source 5",
					Summary: "Sample summary 5",
					Tags:    "tag3",
					Uri:     "https://sample-uri-5",
				}
				digest5, err := sourceRepo.PostSource(context.TODO(), &source5Input)
				Expect(err).ToNot(HaveOccurred())

				source6Input := api.SourceInput{
					Name:    "Sample Source 6",
					Summary: "Sample summary 6",
					Tags:    "tag4",
					Uri:     "https://sample-uri-6",
				}
				digest6, err := sourceRepo.PostSource(context.TODO(), &source6Input)
				Expect(err).ToNot(HaveOccurred())

				source7Input := api.SourceInput{
					Name:    "Sample Source 7",
					Summary: "Sample summary 7",
					Tags:    "tag5",
					Uri:     "https://sample-uri-7",
				}
				digest7, err := sourceRepo.PostSource(context.TODO(), &source7Input)
				Expect(err).ToNot(HaveOccurred())

				// Prepare updated sources with new scores
				updatedSources := []api.Source{
					{UriDigest: digest5, Score: 0.4},
					{UriDigest: digest6, Score: 1},
					{UriDigest: digest7, Score: 0},
				}

				// Update scores
				err = sourceRepo.UpdateSourceScores(context.TODO(), &updatedSources)
				Expect(err).ToNot(HaveOccurred())

				// Verify scores were updated
				src5, err := sourceRepo.GetSourceByUriDigest(context.TODO(), digest5)
				Expect(err).ToNot(HaveOccurred())
				Expect(src5.Score).To(Equal(updatedSources[0].Score))
				Expect(src5.Name).To(Equal(source5Input.Name))
				Expect(src5.Summary).To(Equal(source5Input.Summary))

				src6, err := sourceRepo.GetSourceByUriDigest(context.TODO(), digest6)
				Expect(err).ToNot(HaveOccurred())
				Expect(src6.Score).To(Equal(updatedSources[1].Score))
				Expect(src6.Name).To(Equal(source6Input.Name))
				Expect(src6.Summary).To(Equal(source6Input.Summary))

				src7, err := sourceRepo.GetSourceByUriDigest(context.TODO(), digest7)
				Expect(err).ToNot(HaveOccurred())
				Expect(src7.Score).To(Equal(updatedSources[2].Score))
				Expect(src7.Name).To(Equal(source7Input.Name))
				Expect(src7.Summary).To(Equal(source7Input.Summary))
			})
		})

	})

	Context("Validation tests", func() {
		When("Patching a source that does not exist", func() {
			It("Should return record not found error", func() {
				name := "Twice Updated Sample Source 1"
				tags := "twice-updated-tag1"
				sourceInput := &api.SourcePatchInput{
					Name: &name,
					Tags: &tags,
				}

				err := sourceRepo.PatchSourceByUriDigest(context.TODO(), sourceInput, "invalid-digest")
				Expect(err).To(HaveOccurred())
				Expect(errors.Is(err, gorm.ErrRecordNotFound)).To(BeTrue())
			})
		})

		When("Getting a source that does not exist", func() {
			It("Should return record not found error", func() {
				_, err := sourceRepo.GetSourceByUriDigest(context.TODO(), "invalid-digest")
				Expect(err).To(HaveOccurred())
				Expect(errors.Is(err, gorm.ErrRecordNotFound)).To(BeTrue())
			})
		})
	})
})
