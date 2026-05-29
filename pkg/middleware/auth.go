package middleware

import (
	"fmt"
	"log/slog"
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
		if !strings.HasPrefix(strings.ToLower(authHeader), "bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing token"})
			return
		}

		clientID := c.GetHeader(ClientIDHeader)
		if clientID == "" {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("%s header missing", ClientIDHeader)})
			return
		}

		headerValArr := strings.Fields(strings.TrimSpace(authHeader))
		if len(headerValArr) != 2 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization header format"})
			return
		}

		token, err := jwt.ParseWithClaims(
			headerValArr[1],
			&jwt.RegisteredClaims{},
			func(*jwt.Token) (any, error) { return []byte(jwtSecret), nil },
			jwt.WithAudience(clientID),
			jwt.WithIssuer(TokenIssuer),
			jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Name}),
		)
		if err != nil || !token.Valid {
			if err != nil {
				slog.Error("failed to parse jwt token", "error", err)
			}
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
			return
		}

		c.Next()
	}
}
