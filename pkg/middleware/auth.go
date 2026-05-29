package middleware

import (
	"fmt"
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

		clientID := c.GetHeader(ClientIDHeader)
		if clientID == "" {
			c.JSON(
				http.StatusBadRequest,
				gin.H{"error": fmt.Sprintf("%s header missing", ClientIDHeader)},
			)
			return
		}

		token, err := jwt.ParseWithClaims(
			strings.TrimPrefix(authHeader, "Bearer "),
			&jwt.RegisteredClaims{},
			func(*jwt.Token) (any, error) { return []byte(jwtSecret), nil },
			jwt.WithAudience(clientID),
			jwt.WithIssuer(TokenIssuer),
		)
		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		c.Next()
	}
}
