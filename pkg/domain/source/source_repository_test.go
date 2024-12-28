package source_test

import (
	"context"
	"source-score/pkg/api"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

const (
	uriDigest1 = "8649a4126fb4fc9a750f432b729c8477398cf28ca241403b2cd36a6dc841f441"
)

var _ = Describe("Source model unit tests", func() {
	Context("Happy path", Ordered, func() {
		When("Adding a new source to the DB with valid input", func() {
			It("Should create the source record in the DB", func() {
				err := sourceRepo.PutSource(context.TODO(), &sampleSourceInput1)
				Expect(err).ToNot(HaveOccurred())
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

				err := sourceRepo.UpdateSourceByUriDigest(context.TODO(), sourceInput, uriDigest1)
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
