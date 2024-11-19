package api

import (
	"net/http"
	"source-score/pkg/handlers"

	"github.com/gin-gonic/gin"
)

type router struct {
	pingHandler   *handlers.PingHandler
}

func NewRouter() *router {
	return &router{
		pingHandler:   handlers.NewPingHandler(),
	}
}

func (r *router) GetPing(ctx *gin.Context) {
	message := r.pingHandler.GetPing()

	ctx.JSON(http.StatusOK, gin.H{"data": message})
}
