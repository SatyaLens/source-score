package api

import (
	"log"
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

func (r *router) CreateSource(ctx *gin.Context) {
	log.Println("unimplemented")
}

func (r *router) DeleteSource(ctx *gin.Context, uriDigest string) {
	log.Println("unimplemented")
}

func (r *router) GetSource(ctx *gin.Context, uriDigest string) {
	log.Println("unimplemented")
}

func (r *router) UpdateSource(ctx *gin.Context, uriDigest string) {
	log.Println("unimplemented")
}

func (r *router) GetPing(ctx *gin.Context) {
	message := r.pingHandler.GetPing()

	ctx.JSON(http.StatusOK, gin.H{"data": message})
}
