package source_test

import (
	"context"
	"source-score/pkg/api"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Source model repository layer unit tests", func() {
	Context("Happy path", Ordered, func() {
		When("Adding a new source to the DB with valid input", func() {
			It("Should create the source record in the DB", func() {
				digest, err := sourceRepo.PostSource(context.TODO(), &sampleSourceInput1)
				Expect(err).ToNot(HaveOccurred())
				Expect(digest).To(Equal(uriDigest1))
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

		When("Updating a source by its uri digest", func() {
			It("Should update the correct source record in the DB", func() {
				sourceInput := &api.SourceInput{
					Name:    "Updated Sample Source 1",
					Summary: "Updated Sample summary",
					Tags:    "updated-tag1",
				}

				err := sourceRepo.PatchSourceByUriDigest(context.TODO(), sourceInput, uriDigest1)
				Expect(err).ToNot(HaveOccurred())

				source, err := sourceRepo.GetSourceByUriDigest(context.TODO(), uriDigest1)
				Expect(err).ToNot(HaveOccurred())
				Expect(source.Name).To(BeEquivalentTo(sourceInput.Name))
				Expect(source.Summary).To(BeEquivalentTo(sourceInput.Summary))
				Expect(source.Tags).To(BeEquivalentTo(sourceInput.Tags))
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
	})
})
