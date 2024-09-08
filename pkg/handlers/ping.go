package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type pong struct {
	message string
}

func NewPong() *pong {
	return &pong{
		message: "pong",
	}
}

func (pg pong) GetPing(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"data": pg.message})
}