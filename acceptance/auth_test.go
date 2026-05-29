package acceptance_test

import (
	"net/http"
	"net/url"
	"source-score/pkg/middleware"
	"time"

	"github.com/golang-jwt/jwt/v5"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Auth middleware tests", func() {
	protectedEndpoint, err := url.JoinPath(baseUrl, "/api/v1/claims")
	Expect(err).To(BeNil())

	XContext("Validation tests", func() {
		It("should reject requests with an invalid jwt token", func() {
			resp, err := doRequestWithAuthToken(http.MethodGet, protectedEndpoint, "invalid.jwt.token")
			Expect(err).To(BeNil())
			defer resp.Body.Close()

			Expect(resp.StatusCode).To(Equal(http.StatusUnauthorized))
		})

		It("should reject requests with an expired jwt token", func() {
			token := signAcceptanceToken(jwt.RegisteredClaims{
				Audience:  []string{commonHeaders[middleware.ClientIDHeader]},
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(-time.Hour)),
				IssuedAt:  jwt.NewNumericDate(time.Now().Add(-2 * time.Hour)),
				Issuer:    middleware.TokenIssuer,
			})

			resp, err := doRequestWithAuthToken(http.MethodPost, protectedEndpoint, token)
			Expect(err).To(BeNil())
			defer resp.Body.Close()

			Expect(resp.StatusCode).To(Equal(http.StatusUnauthorized))
		})

		It("should reject requests when jwt audience does not match the client id header", func() {
			token := signAcceptanceToken(jwt.RegisteredClaims{
				Audience:  []string{"another-client"},
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
				IssuedAt:  jwt.NewNumericDate(time.Now()),
				Issuer:    middleware.TokenIssuer,
			})

			resp, err := doRequestWithAuthToken(http.MethodPost, protectedEndpoint, token)
			Expect(err).To(BeNil())
			defer resp.Body.Close()

			Expect(resp.StatusCode).To(Equal(http.StatusUnauthorized))
		})

		It("should reject requests when jwt issuer does not match the configured token issuer", func() {
			token := signAcceptanceToken(jwt.RegisteredClaims{
				Audience:  []string{commonHeaders[middleware.ClientIDHeader]},
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
				IssuedAt:  jwt.NewNumericDate(time.Now()),
				Issuer:    "another-issuer",
			})

			resp, err := doRequestWithAuthToken(http.MethodPost, protectedEndpoint, token)
			Expect(err).To(BeNil())
			defer resp.Body.Close()

			Expect(resp.StatusCode).To(Equal(http.StatusUnauthorized))
		})
	})
})

func doRequestWithAuthToken(method, endpoint, authToken string) (*http.Response, error) {
	req, err := http.NewRequest(method, endpoint, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set(middleware.ClientIDHeader, commonHeaders[middleware.ClientIDHeader])
	req.Header.Set("Authorization", "Bearer "+authToken)

	return client.Do(req)
}

func signAcceptanceToken(claims jwt.RegisteredClaims) string {
	// signing secret is the same as that hardcoded for app container in docker compose
	signedToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte("default-secret-string"))
	Expect(err).To(BeNil())

	return signedToken
}
