package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

const (
	ClientIDHeader = "Client-ID"
	TokenIssuer    = "source-score"
)

func AuthTokenMiddleware(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing token"})
			return
		}

		token, err := jwt.ParseWithClaims(
			strings.TrimPrefix(authHeader, "Bearer "),
			jwt.RegisteredClaims{},
			func(*jwt.Token) (any, error) { return []byte(jwtSecret), nil },
		)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}

		c.Next()
	}
}
