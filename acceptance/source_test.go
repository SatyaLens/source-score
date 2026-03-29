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

var (
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
					Name:    &updatedName,
					Summary: &updatedSummary,
					Tags:    &updatedTags,
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
				name := "twice updated name"
				tags := "twice-updated-tag"
				updatedSrcInput := api.SourcePatchInput{
					Name: &name,
					Tags: &tags,
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
		When("GET request is sent for an invalid source", func() {
			It("should return 404 error", func() {
				srcUrl, err := url.JoinPath(endpoint, "invalid-digest")
				Expect(err).To(BeNil())
				resp, err := http.Get(srcUrl)
				Expect(err).To(BeNil())
				defer resp.Body.Close()
				Expect(resp.StatusCode).To(Equal(http.StatusNotFound))
			})
		})

		When("POST request with missing required fields is sent", func() {
			It("should return 400 Bad Request with error message", func() {
				invalidBody := []byte(`{}`)
				resp, err := http.Post(endpoint, "application/json", bytes.NewBuffer(invalidBody))
				Expect(err).To(BeNil())
				defer resp.Body.Close()
				Expect(resp.StatusCode).To(Equal(http.StatusBadRequest))

				var errResp map[string]any
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

				var errResp map[string]any
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
				Expect(strings.ToLower(errResp["error"].(string))).
					To(ContainSubstring("field validation for 'name' failed on the 'required' tag"))
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
				Expect(strings.ToLower(errResp["error"].(string))).
					To(ContainSubstring("field validation for 'summary' failed on the 'required' tag"))
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
				Expect(strings.ToLower(errResp["error"].(string))).
					To(ContainSubstring("field validation for 'tags' failed on the 'required' tag"))
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
				Expect(strings.ToLower(errResp["error"].(string))).
					To(ContainSubstring("field validation for 'uri' failed on the 'required' tag"))
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

		When("PATCH request with empty name is sent", func() {
			It("should return 400 Bad Request with error message", func() {
				name := ""
				invalidInput := api.SourcePatchInput{
					Name: &name,
				}
				body, _ := json.Marshal(invalidInput)
				srcUrl, err := url.JoinPath(endpoint, uriDigest1)
				Expect(err).To(BeNil())
				req, err := http.NewRequest(http.MethodPatch, srcUrl, bytes.NewBuffer(body))
				Expect(err).To(BeNil())
				req.Header.Set("Content-Type", "application/json")
				resp, err := client.Do(req)
				Expect(err).To(BeNil())
				defer resp.Body.Close()
				Expect(resp.StatusCode).To(Equal(http.StatusBadRequest))

				var errResp map[string]any
				err = json.NewDecoder(resp.Body).Decode(&errResp)
				Expect(err).To(BeNil())
				Expect(strings.ToLower(errResp["error"].(string))).To(ContainSubstring("name validation failed"))

				By("verifying validation is not running for nil fields")
				Expect(strings.ToLower(errResp["error"].(string))).ToNot(ContainSubstring("tags validation failed"))
				Expect(strings.ToLower(errResp["error"].(string))).ToNot(ContainSubstring("summary validation failed"))
			})
		})

		When("PATCH request with empty summary is sent", func() {
			It("should return 400 Bad Request with error message", func() {
				summary := ""
				invalidInput := api.SourcePatchInput{
					Summary: &summary,
				}
				body, _ := json.Marshal(invalidInput)
				srcUrl, err := url.JoinPath(endpoint, uriDigest1)
				Expect(err).To(BeNil())
				req, err := http.NewRequest(http.MethodPatch, srcUrl, bytes.NewBuffer(body))
				Expect(err).To(BeNil())
				req.Header.Set("Content-Type", "application/json")
				resp, err := client.Do(req)
				Expect(err).To(BeNil())
				defer resp.Body.Close()
				Expect(resp.StatusCode).To(Equal(http.StatusBadRequest))

				var errResp map[string]any
				err = json.NewDecoder(resp.Body).Decode(&errResp)
				Expect(err).To(BeNil())
				Expect(strings.ToLower(errResp["error"].(string))).To(ContainSubstring("summary validation failed"))

				By("verifying validation is not running for nil fields")
				Expect(strings.ToLower(errResp["error"].(string))).ToNot(ContainSubstring("tags validation failed"))
				Expect(strings.ToLower(errResp["error"].(string))).ToNot(ContainSubstring("name validation failed"))
			})
		})

		When("PATCH request with empty tags is sent", func() {
			It("should return 400 Bad Request with error message", func() {
				tags := ""
				invalidInput := api.SourcePatchInput{
					Tags: &tags,
				}
				body, _ := json.Marshal(invalidInput)
				srcUrl, err := url.JoinPath(endpoint, uriDigest1)
				Expect(err).To(BeNil())
				req, err := http.NewRequest(http.MethodPatch, srcUrl, bytes.NewBuffer(body))
				Expect(err).To(BeNil())
				req.Header.Set("Content-Type", "application/json")
				resp, err := client.Do(req)
				Expect(err).To(BeNil())
				defer resp.Body.Close()
				Expect(resp.StatusCode).To(Equal(http.StatusBadRequest))

				var errResp map[string]any
				err = json.NewDecoder(resp.Body).Decode(&errResp)
				Expect(err).To(BeNil())
				Expect(strings.ToLower(errResp["error"].(string))).To(ContainSubstring("tags validation failed"))

				By("verifying validation is not running for nil fields")
				Expect(strings.ToLower(errResp["error"].(string))).ToNot(ContainSubstring("summary validation failed"))
				Expect(strings.ToLower(errResp["error"].(string))).ToNot(ContainSubstring("name validation failed"))
			})
		})

		When("PATCH request with tags containing spaces is sent", func() {
			It("should return 400 Bad Request with error message", func() {
				tags := "tag1, tag2"
				invalidInput := api.SourcePatchInput{
					Tags: &tags,
				}
				body, _ := json.Marshal(invalidInput)
				srcUrl, err := url.JoinPath(endpoint, uriDigest1)
				Expect(err).To(BeNil())
				req, err := http.NewRequest(http.MethodPatch, srcUrl, bytes.NewBuffer(body))
				Expect(err).To(BeNil())
				req.Header.Set("Content-Type", "application/json")
				resp, err := client.Do(req)
				Expect(err).To(BeNil())
				defer resp.Body.Close()
				Expect(resp.StatusCode).To(Equal(http.StatusBadRequest))

				var errResp map[string]any
				err = json.NewDecoder(resp.Body).Decode(&errResp)
				Expect(err).To(BeNil())
				Expect(strings.ToLower(errResp["error"].(string))).To(ContainSubstring("tags validation failed"))

				By("verifying validation is not running for nil fields")
				Expect(strings.ToLower(errResp["error"].(string))).ToNot(ContainSubstring("summary validation failed"))
				Expect(strings.ToLower(errResp["error"].(string))).ToNot(ContainSubstring("name validation failed"))
			})
		})

		When("PATCH request for a source that doesn't exist is sent", func() {
			It("should return 404 Not Found with error message", func() {
				srcName := "new name"
				validInput := api.SourcePatchInput{
					Name: &srcName,
				}
				body, _ := json.Marshal(validInput)
				srcUrl, err := url.JoinPath(endpoint, "invalid-uri-digest")
				Expect(err).To(BeNil())
				req, err := http.NewRequest(http.MethodPatch, srcUrl, bytes.NewBuffer(body))
				Expect(err).To(BeNil())
				req.Header.Set("Content-Type", "application/json")
				resp, err := client.Do(req)
				Expect(err).To(BeNil())
				defer resp.Body.Close()
				Expect(resp.StatusCode).To(Equal(http.StatusNotFound))

				var errResp map[string]any
				err = json.NewDecoder(resp.Body).Decode(&errResp)
				Expect(err).To(BeNil())
				Expect(errResp["error"]).ToNot(BeNil())
				Expect(strings.ToLower(errResp["error"].(string))).To(ContainSubstring("source not found"))
			})
		})

	})
})
