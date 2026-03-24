package acceptance_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"
	"source-score/pkg/api"
	"strings"

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
	body1, err := json.Marshal(sourceInput1)
	Expect(err).To(BeNil())
	body2, err := json.Marshal(sourceInput2)
	Expect(err).To(BeNil())

	Context("Happy path tests", Ordered, func() {
		When("valid POST requests are sent with source model inputs", func() {
			It("should return successful response", func() {
				var respBody api.CreateSourceResponse

				resp, err := http.Post(
					endpoint,
					"application/json",
					bytes.NewBuffer(body1),
				)
				Expect(err).To(BeNil())
				Expect(resp.StatusCode).To(BeEquivalentTo(http.StatusCreated))
				err = json.NewDecoder(resp.Body).Decode(&respBody)
				Expect(err).To(BeNil())
				Expect(respBody.UriDigest).To(Equal(uriDigest1))
				resp.Body.Close()

				resp, err = http.Post(
					endpoint,
					"application/json",
					bytes.NewBuffer(body2),
				)
				Expect(err).To(BeNil())
				Expect(resp.StatusCode).To(BeEquivalentTo(http.StatusCreated))
				err = json.NewDecoder(resp.Body).Decode(&respBody)
				Expect(err).To(BeNil())
				defer resp.Body.Close()
				Expect(respBody.UriDigest).To(Equal(uriDigest2))
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

		When("GET request is sent to retrieve all sources", func() {
			It("should return all sources", func() {
				sourcesUrl, err := url.JoinPath(baseUrl, "/api/v1/sources")
				Expect(err).To(BeNil())
				resp, err := http.Get(sourcesUrl)
				Expect(err).To(BeNil())
				defer resp.Body.Close()
				Expect(resp.StatusCode).To(Equal(http.StatusOK))

				var sources []api.Source
				err = json.NewDecoder(resp.Body).Decode(&sources)
				Expect(err).To(BeNil())
				Expect(len(sources)).To(Equal(2))
				Expect(sources).To(ContainElements(
					sampleSource1,
					sampleSource2,
				))
			})
		})

		When("PATCH request is sent to update all the fields of the created source", func() {
			It("should update the source record", func() {
				updatedSrcInput := api.SourcePatchInput{
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

		When("PATCH request is sent to update some fields of the created source", func() {
			It("should update the source record", func() {
				updatedSrcInput := api.SourcePatchInput{
					Name: "twice updated name",
					Tags: "twice-updated-tag",
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
				Expect(src.Name).To(Equal("twice updated name"))
				Expect(src.Summary).To(Equal(updatedSummary))
				Expect(src.Tags).To(Equal("twice-updated-tag"))
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

	Context("Validation tests", func() {
		When("POST request with missing required fields is sent", func() {
			It("should return 400 Bad Request with error message", func() {
				invalidBody := []byte(`{}`)
				resp, err := http.Post(endpoint, "application/json", bytes.NewBuffer(invalidBody))
				Expect(err).To(BeNil())
				defer resp.Body.Close()
				Expect(resp.StatusCode).To(Equal(http.StatusBadRequest))

				var errResp map[string]interface{}
				err = json.NewDecoder(resp.Body).Decode(&errResp)
				Expect(err).To(BeNil())
				Expect(errResp["error"]).ToNot(BeNil())
			})
		})

		When("POST request with missing tags field is sent", func() {
			It("should return 400 Bad Request with error message", func() {
				invalidBody := []byte(`{"name":"valid name","summary":"valid summary","uri":"https://example.com"}`)
				resp, err := http.Post(endpoint, "application/json", bytes.NewBuffer(invalidBody))
				Expect(err).To(BeNil())
				defer resp.Body.Close()
				Expect(resp.StatusCode).To(Equal(http.StatusBadRequest))

				var errResp map[string]interface{}
				err = json.NewDecoder(resp.Body).Decode(&errResp)
				Expect(err).To(BeNil())
				Expect(errResp["error"]).ToNot(BeNil())
			})
		})

		When("POST request with empty name is sent", func() {
			It("should return 400 Bad Request with error message", func() {
				invalidInput := api.SourceInput{
					Name:    "",
					Summary: "valid summary",
					Tags:    "tag1,tag2",
					Uri:     "https://example.com",
				}
				body, _ := json.Marshal(invalidInput)
				resp, err := http.Post(endpoint, "application/json", bytes.NewBuffer(body))
				Expect(err).To(BeNil())
				defer resp.Body.Close()
				Expect(resp.StatusCode).To(Equal(http.StatusBadRequest))

				var errResp map[string]any
				err = json.NewDecoder(resp.Body).Decode(&errResp)
				Expect(err).To(BeNil())
				Expect(strings.ToLower(errResp["error"].(string))).To(ContainSubstring("name validation failed"))
			})
		})

		When("POST request with empty summary is sent", func() {
			It("should return 400 Bad Request with error message", func() {
				invalidInput := api.SourceInput{
					Name:    "valid name",
					Summary: "",
					Tags:    "tag1,tag2",
					Uri:     "https://example.com",
				}
				body, _ := json.Marshal(invalidInput)
				resp, err := http.Post(endpoint, "application/json", bytes.NewBuffer(body))
				Expect(err).To(BeNil())
				defer resp.Body.Close()
				Expect(resp.StatusCode).To(Equal(http.StatusBadRequest))

				var errResp map[string]any
				err = json.NewDecoder(resp.Body).Decode(&errResp)
				Expect(err).To(BeNil())
				Expect(strings.ToLower(errResp["error"].(string))).To(ContainSubstring("summary validation failed"))
			})
		})

		When("POST request with empty tags is sent", func() {
			It("should return 400 Bad Request with error message", func() {
				invalidInput := api.SourceInput{
					Name:    "valid name",
					Summary: "valid summary",
					Tags:    "",
					Uri:     "https://example.com",
				}
				body, _ := json.Marshal(invalidInput)
				resp, err := http.Post(endpoint, "application/json", bytes.NewBuffer(body))
				Expect(err).To(BeNil())
				defer resp.Body.Close()
				Expect(resp.StatusCode).To(Equal(http.StatusBadRequest))

				var errResp map[string]any
				err = json.NewDecoder(resp.Body).Decode(&errResp)
				Expect(err).To(BeNil())
				Expect(strings.ToLower(errResp["error"].(string))).To(ContainSubstring("tags validation failed"))
			})
		})

		When("POST request with tags containing spaces is sent", func() {
			It("should return 400 Bad Request with error message", func() {
				invalidInput := api.SourceInput{
					Name:    "valid name",
					Summary: "valid summary",
					Tags:    "tag1, tag2",
					Uri:     "https://example.com",
				}
				body, _ := json.Marshal(invalidInput)
				resp, err := http.Post(endpoint, "application/json", bytes.NewBuffer(body))
				Expect(err).To(BeNil())
				defer resp.Body.Close()
				Expect(resp.StatusCode).To(Equal(http.StatusBadRequest))

				var errResp map[string]any
				err = json.NewDecoder(resp.Body).Decode(&errResp)
				Expect(err).To(BeNil())
				Expect(strings.ToLower(errResp["error"].(string))).To(ContainSubstring("tags validation failed"))
			})
		})

		When("POST request with empty uri is sent", func() {
			It("should return 400 Bad Request with error message", func() {
				invalidInput := api.SourceInput{
					Name:    "valid name",
					Summary: "valid summary",
					Tags:    "tag1,tag2",
					Uri:     "",
				}
				body, _ := json.Marshal(invalidInput)
				resp, err := http.Post(endpoint, "application/json", bytes.NewBuffer(body))
				Expect(err).To(BeNil())
				defer resp.Body.Close()
				Expect(resp.StatusCode).To(Equal(http.StatusBadRequest))

				var errResp map[string]any
				err = json.NewDecoder(resp.Body).Decode(&errResp)
				Expect(err).To(BeNil())
				Expect(strings.ToLower(errResp["error"].(string))).To(ContainSubstring("uri validation failed"))
			})
		})

		When("POST request with non-https uri is sent", func() {
			It("should return 400 Bad Request with error message", func() {
				invalidInput := api.SourceInput{
					Name:    "valid name",
					Summary: "valid summary",
					Tags:    "tag1,tag2",
					Uri:     "http://example.com",
				}
				body, _ := json.Marshal(invalidInput)
				resp, err := http.Post(endpoint, "application/json", bytes.NewBuffer(body))
				Expect(err).To(BeNil())
				defer resp.Body.Close()
				Expect(resp.StatusCode).To(Equal(http.StatusBadRequest))

				var errResp map[string]any
				err = json.NewDecoder(resp.Body).Decode(&errResp)
				Expect(err).To(BeNil())
				Expect(strings.ToLower(errResp["error"].(string))).To(ContainSubstring("uri validation failed"))
			})
		})
	})
})
