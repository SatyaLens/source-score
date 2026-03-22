package acceptance_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"
	"source-score/pkg/api"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

const (
	updatedName    = "updated source"
	updatedSummary = "updated summary"
	updatedTags    = "tag1,tag2"
)

var _ = Describe("Source model tests", func() {
	endpoint, err := url.JoinPath(baseUrl, "/api/v1/source")
	Expect(err).To(BeNil())
	body, err := json.Marshal(sourceInput1)
	Expect(err).To(BeNil())

	Context("Happy path tests", Ordered, func() {
		When("valid POST request is sent with source model input", func() {
			It("should return successful response", func() {
				resp, err := http.Post(
					endpoint,
					"application/json",
					bytes.NewBuffer(body),
				)
				Expect(err).To(BeNil())
				defer resp.Body.Close()
				Expect(resp.StatusCode).To(BeEquivalentTo(http.StatusCreated))
			})
		})

		When("GET request is sent to query the created source", func() {
			It("should return the created source", func() {
				srcUrl, err := url.JoinPath(endpoint, uriDigest1)
				Expect(err).To(BeNil())
				resp, err := http.Get(srcUrl)
				Expect(err).To(BeNil())
				defer resp.Body.Close()
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
			})
		})

		When("PATCH request is sent to update the created source", func() {
			It("should update the source record", func() {
				updatedSrcInput := api.SourceInput{
					Name:    updatedName,
					Summary: updatedSummary,
					Tags:    updatedTags,
				}
				reqBody, err := json.Marshal(updatedSrcInput)
				Expect(err).To(BeNil())

				srcUrl, err := url.JoinPath(endpoint, uriDigest1)
				Expect(err).To(BeNil())
				req, err := http.NewRequest(
					http.MethodPatch,
					srcUrl,
					bytes.NewBuffer(reqBody))
				Expect(err).To(BeNil())
				req.Header.Set("Content-Type", "application/json")
				resp, err := client.Do(req)
				Expect(err).To(BeNil())
				Expect(resp.StatusCode).To(Equal(http.StatusNoContent))
				resp.Body.Close()

				By("verifying source got updated")
				var src api.Source
				resp, err = http.Get(srcUrl)
				Expect(err).To(BeNil())
				defer resp.Body.Close()
				Expect(resp.StatusCode).To(Equal(http.StatusOK))

				err = json.NewDecoder(resp.Body).Decode(&src)
				Expect(err).To(BeNil())
				Expect(src.Name).To(Equal(updatedName))
				Expect(src.Summary).To(Equal(updatedSummary))
				Expect(src.Tags).To(Equal(updatedTags))
			})
		})

		When("DELETE request is sent to delete the created source", func() {
			It("should delete the created source", func() {
				srcUrl, err := url.JoinPath(endpoint, uriDigest1)
				Expect(err).To(BeNil())
				req, err := http.NewRequest(
					http.MethodDelete,
					srcUrl,
					nil,
				)
				Expect(err).To(BeNil())

				resp, err := client.Do(req)
				Expect(err).To(BeNil())
				defer resp.Body.Close()
				Expect(resp.StatusCode).To(Equal(http.StatusNoContent))

				By("verifying source got deleted")
				resp, err = http.Get(srcUrl)
				Expect(err).To(BeNil())
				Expect(resp.StatusCode).To(Equal(http.StatusNotFound))
			})
		})
	})
})
