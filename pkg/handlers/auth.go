package handlers

import (
	"fmt"
	"log/slog"
	"net/http"
	"source-score/pkg/middleware"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type AuthHandler struct {
	jwtSecret string
}

func NewAuthHandler(secret string) *AuthHandler {
	return &AuthHandler{jwtSecret: secret}
}

func (ah *AuthHandler) GetAuthToken(c *gin.Context) {
	clientID := c.GetHeader(middleware.ClientIDHeader)
	if clientID == "" {
		c.JSON(
			http.StatusBadRequest,
			gin.H{"error": fmt.Sprintf("%s header missing", middleware.ClientIDHeader)},
		)
		return
	}

	token, err := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(4 * time.Hour)),
			Audience:  []string{clientID},
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    middleware.TokenIssuer,
		},
	).SignedString(ah.jwtSecret)
	if err != nil {
		slog.Error("failed to generate auth token", "error", err)
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": "failed to generate auth token"},
		)
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}
