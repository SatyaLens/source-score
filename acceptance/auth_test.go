package acceptance_test

import (
	"encoding/json"
	"net/http"
	"net/url"
	"source-score/pkg/middleware"
	"time"

	"github.com/golang-jwt/jwt/v5"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var (
	protectedEndpoint string
	err               error
)

var _ = Describe("Auth middleware tests", Ordered, func() {
	protectedEndpoint, err = url.JoinPath(baseUrl, "/api/v1/claims")
	Expect(err).To(BeNil())

	Context("Validation tests", func() {
		It("should reject requests without the client id header", func() {
			resp, err := doRequestWithHeaders(http.MethodGet, protectedEndpoint, map[string]string{
				"Authorization": commonHeaders["Authorization"],
			})
			Expect(err).To(BeNil())
			Expect(resp.StatusCode).To(Equal(http.StatusBadRequest))

			defer resp.Body.Close()
			var respBody map[string]string
			err = json.NewDecoder(resp.Body).Decode(&respBody)
			Expect(err).To(BeNil())
			Expect(respBody["msg"]).To(ContainSubstring("Client-ID is required"))
		})

		It("should reject requests without the authorization header", func() {
			resp, err := doRequestWithHeaders(http.MethodGet, protectedEndpoint, map[string]string{
				middleware.ClientIDHeader: commonHeaders[middleware.ClientIDHeader],
			})
			Expect(err).To(BeNil())
			Expect(resp.StatusCode).To(Equal(http.StatusUnauthorized))

			defer resp.Body.Close()
			var respBody map[string]string
			err = json.NewDecoder(resp.Body).Decode(&respBody)
			Expect(err).To(BeNil())
			Expect(respBody["error"]).To(ContainSubstring("missing token"))
		})

		It("should reject requests with an invalid jwt token", func() {
			resp, err := doRequestWithAuthToken(http.MethodGet, protectedEndpoint, "invalid.jwt.token")
			Expect(err).To(BeNil())
			Expect(resp.StatusCode).To(Equal(http.StatusUnauthorized))

			defer resp.Body.Close()
			var respBody map[string]string
			err = json.NewDecoder(resp.Body).Decode(&respBody)
			Expect(err).To(BeNil())
			Expect(respBody["error"]).To(ContainSubstring("invalid or expired token"))
		})

		It("should reject requests with an expired jwt token", func() {
			token := signAcceptanceToken(jwt.RegisteredClaims{
				Audience:  []string{commonHeaders[middleware.ClientIDHeader]},
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(-time.Hour)),
				IssuedAt:  jwt.NewNumericDate(time.Now().Add(-2 * time.Hour)),
				Issuer:    middleware.TokenIssuer,
			})

			resp, err := doRequestWithAuthToken(http.MethodGet, protectedEndpoint, token)
			Expect(err).To(BeNil())
			Expect(resp.StatusCode).To(Equal(http.StatusUnauthorized))

			defer resp.Body.Close()
			var respBody map[string]string
			err = json.NewDecoder(resp.Body).Decode(&respBody)
			Expect(err).To(BeNil())
			Expect(respBody["error"]).To(ContainSubstring("invalid or expired token"))
		})

		It("should reject requests when jwt audience does not match the client id header", func() {
			token := signAcceptanceToken(jwt.RegisteredClaims{
				Audience:  []string{"another-client"},
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
				IssuedAt:  jwt.NewNumericDate(time.Now()),
				Issuer:    middleware.TokenIssuer,
			})

			resp, err := doRequestWithAuthToken(http.MethodGet, protectedEndpoint, token)
			Expect(err).To(BeNil())
			Expect(resp.StatusCode).To(Equal(http.StatusUnauthorized))

			defer resp.Body.Close()
			var respBody map[string]string
			err = json.NewDecoder(resp.Body).Decode(&respBody)
			Expect(err).To(BeNil())
			Expect(respBody["error"]).To(ContainSubstring("invalid or expired token"))
		})

		It("should reject requests when jwt issuer does not match the configured token issuer", func() {
			token := signAcceptanceToken(jwt.RegisteredClaims{
				Audience:  []string{commonHeaders[middleware.ClientIDHeader]},
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
				IssuedAt:  jwt.NewNumericDate(time.Now()),
				Issuer:    "another-issuer",
			})

			resp, err := doRequestWithAuthToken(http.MethodGet, protectedEndpoint, token)
			Expect(err).To(BeNil())
			Expect(resp.StatusCode).To(Equal(http.StatusUnauthorized))

			defer resp.Body.Close()
			var respBody map[string]string
			err = json.NewDecoder(resp.Body).Decode(&respBody)
			Expect(err).To(BeNil())
			Expect(respBody["error"]).To(ContainSubstring("invalid or expired token"))
		})
	})
})

func doRequestWithAuthToken(method, endpoint, authToken string) (*http.Response, error) {
	return doRequestWithHeaders(method, endpoint, map[string]string{
		middleware.ClientIDHeader: commonHeaders[middleware.ClientIDHeader],
		"Authorization":           "Bearer " + authToken,
	})
}

func doRequestWithHeaders(method, endpoint string, headers map[string]string) (*http.Response, error) {
	req, err := http.NewRequest(method, endpoint, nil)
	if err != nil {
		return nil, err
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	return client.Do(req)
}

func signAcceptanceToken(claims jwt.RegisteredClaims) string {
	// signing secret is the same as that hardcoded for app container in docker compose
	signedToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte("default-secret-string"))
	Expect(err).To(BeNil())

	return signedToken
}
