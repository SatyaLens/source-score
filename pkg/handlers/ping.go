package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type PongHandler struct {
	message string
}

func NewPongHandler() *PongHandler {
	return &PongHandler{
		message: "Pong",
	}
}

func (ph PongHandler) GetPing(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"data": ph.message})
}