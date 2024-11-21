package acceptance

import (
	"net/http"
	"path"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Acceptance Tests", func() {
	BeforeAll(func() {
		if port == "" {
			port = "8080"
		}

		baseUrl = "localhost:" + port
	})
	Context("Testing /ping endpoint", func() {
		endpoint := path.Join(baseUrl, "ping")
		resp, err := http.Get(endpoint)

		Expect(err).To(BeNil())
		Expect(resp.StatusCode).To(BeEquivalentTo(http.StatusOK))

		defer resp.Body.Close()
	})
})
