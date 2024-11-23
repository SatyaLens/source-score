package acceptance_test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("API Tests", func() {
	Context("Testing /ping endpoint", func() {
		endpoint, err := url.JoinPath(baseUrl, "ping")
		Expect(err).To(BeNil())

		When("GET request is sent to /ping", func() {
			It("should get Pong message in reponse", func() {
				resp, err := http.Get(endpoint)
				Expect(err).To(BeNil())
				Expect(resp.StatusCode).To(BeEquivalentTo(http.StatusOK))

				defer resp.Body.Close()
				body, err := io.ReadAll(resp.Body)
				Expect(err).To(BeNil())

				var respBody responseBody
				err = json.Unmarshal(body, &respBody)
				Expect(err).To(BeNil())	
				Expect(respBody.Data).To(BeEquivalentTo("Pong"))
			})
		})
	})
})
