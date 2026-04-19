package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type PingHandler struct {
	message string
}

func NewPingHandler() *PingHandler {
	return &PingHandler{
		message: "Pong",
	}
}

func (ph PingHandler) GetPing(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"data": ph.message})
}
