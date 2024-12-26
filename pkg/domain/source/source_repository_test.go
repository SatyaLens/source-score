package source_test

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// const (
// 	uriDigest1 = "8649a4126fb4fc9a750f432b729c8477398cf28ca241403b2cd36a6dc841f441"
// )

var _ = Describe("Source model unit tests", func() {
	Context("Creating sources", func() {
		When("Adding a new source to the DB with valid input", func() {
			It("Should create the source record in the DB", func() {
				err := sourceRepo.PutSource(context.TODO(), &sampleSourceInput1)
				Expect(err).ToNot(HaveOccurred())
			})
		})
	})

	// Context("Updating sources", func() {

	// })

	// Context("Getting sources", func() {

	// })

	// Context("Deleting sources", func() {

	// })
})
