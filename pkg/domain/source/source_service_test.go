package source_test

import (
	"context"
	"source-score/pkg/api"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Source model service layer unit test", func() {
	Context("Happy path", Ordered, func() {
		When("Adding a new source with valid input", func() {
			It("Should pass the data to the repository layer", func() {
				err := sourceSvc.PutSource(context.TODO(), &sampleSourceInput1)
				Expect(err).ToNot(HaveOccurred())

				Expect(fakeSourceRepo.PutSourceCallCount()).To(Equal(1))
				_, srcInput := fakeSourceRepo.PutSourceArgsForCall(0)
				Expect(srcInput).To(Equal(&sampleSourceInput1))
			})
		})

		When("Retrieving a source by its uri digest", func() {
			fakeSourceRepo.GetSourceByUriDigestReturnsOnCall(0, &sampleSource1, nil)
			It("Should pass the digest to the repo layer", func() {
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
				sourceInput := &api.SourceInput{
					Name:    "Updated Sample Source 1",
					Summary: "Updated Sample summary",
					Tags:    "updated-tag1",
				}
				updatedSource := sampleSource1
				updatedSource.Name = "Updated Sample Source 1"
				updatedSource.Summary = "Updated Sample summary"
				updatedSource.Tags = "updated-tag1"
				fakeSourceRepo.GetSourceByUriDigestReturnsOnCall(1, &sampleSource1, nil)
				fakeSourceRepo.GetSourceByUriDigestReturnsOnCall(2, &updatedSource, nil)

				err := sourceSvc.UpdateSourceByUriDigest(context.TODO(), sourceInput, uriDigest1)
				Expect(err).ToNot(HaveOccurred())

				source, err := sourceSvc.GetSourceByUriDigest(context.TODO(), uriDigest1)
				Expect(err).ToNot(HaveOccurred())
				Expect(source.Name).To(BeEquivalentTo(sourceInput.Name))
				Expect(source.Summary).To(BeEquivalentTo(sourceInput.Summary))
				Expect(source.Tags).To(BeEquivalentTo(sourceInput.Tags))
				Expect(source.Uri).To(BeEquivalentTo(sampleSourceInput1.Uri))
				Expect(source.UriDigest).To(BeEquivalentTo(uriDigest1))			})
		})

		// When("Deleting a source by its uri digest", func() {
		// 	It("Should delete the correct source record from the DB", func() {
		// 		source := &api.Source{
		// 			UriDigest: uriDigest1,
		// 		}

		// 		err := sourceRepo.DeleteSourceByUriDigest(context.TODO(), source)
		// 		Expect(err).ToNot(HaveOccurred())

		// 		_, err = sourceRepo.GetSourceByUriDigest(context.TODO(), uriDigest1)
		// 		Expect(err).To(HaveOccurred())
		// 		Expect(err.Error()).To(ContainSubstring("record not found"))
		// 	})
		// })
	})
})
