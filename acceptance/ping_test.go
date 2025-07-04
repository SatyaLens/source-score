package acceptance_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"source-score/pkg/handlers"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"google.golang.org/protobuf/proto"
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

		When("POST request is sent to /ping", func() {
			It("should sent message included in the reponse", func() {
				msg := handlers.IncomingMessage{
					Message: "sample incoming message",
				}
				data, err := proto.Marshal(&msg)
				Expect(err).To(BeNil())

				resp, err := http.Post(
					endpoint,
					"application/x-protobuf",
					bytes.NewReader(data),
				)
				Expect(err).To(BeNil())
				Expect(resp.StatusCode).To(BeEquivalentTo(http.StatusOK))

				defer resp.Body.Close()
				body, err := io.ReadAll(resp.Body)
				Expect(err).To(BeNil())

				var respBody responseBody
				err = json.Unmarshal(body, &respBody)
				Expect(err).To(BeNil())
				Expect(respBody.Data).To(ContainSubstring("sample incoming message"))
			})
		})
	})
})
