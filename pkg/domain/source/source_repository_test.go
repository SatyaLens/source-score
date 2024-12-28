package source_test

import (
	"context"

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

		// When("Updating a source by its uri digest", func() {
		// 	It("Should update the source record in the DB", func() {
				
		// 	})
		// })
	})

	// Context("Updating sources", func() {

	// })

	// Context("Getting sources", func() {

	// })

	// Context("Deleting sources", func() {

	// })
})
