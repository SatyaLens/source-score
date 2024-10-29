package api

import (
	"net/http"
	"source-score/pkg/handlers"

	"github.com/gin-gonic/gin"
)

type router struct {
	pingHandler   *handlers.PingHandler
	sourceHandler *handlers.SourceHandler
}

func NewRouter() *router {
	return &router{
		pingHandler:   handlers.NewPingHandler(),
		sourceHandler: handlers.NewSourceHandler(),
	}
}

func (r *router) GetPing(ctx *gin.Context) {
	message := r.pingHandler.GetPing()

	ctx.JSON(http.StatusOK, gin.H{"data": message})
}

func (r *router) CreateSource(ctx *gin.Context) {
	message := r.sourceHandler.CreateSource()

	ctx.JSON(http.StatusOK, gin.H{"data": message})
}

func (r *router) DeleteSource(ctx *gin.Context, uriDigest string) {
	message := r.sourceHandler.DeleteSource()

	ctx.JSON(http.StatusOK, gin.H{"data": message})
}

func (r *router) GetSource(ctx *gin.Context, uriDigest string) {
	message := r.sourceHandler.GetSource()

	ctx.JSON(http.StatusOK, gin.H{"data": message})
}

func (r *router) UpdateSource(ctx *gin.Context, uriDigest string) {
	message := r.sourceHandler.UpdateSource()

	ctx.JSON(http.StatusOK, gin.H{"data": message})
}