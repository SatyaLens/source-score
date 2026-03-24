package http

import (
	"context"
	"net/http"
	"source-score/pkg/domain/source"
	"source-score/pkg/handlers"

	"github.com/gin-gonic/gin"
)

type router struct {
	pingHandler *handlers.PingHandler
	srcHandler  *handlers.SourceHandler
}

func NewRouter(ctx context.Context, sourceSvc source.SourceService) *router {
	return &router{
		pingHandler: handlers.NewPingHandler(),
		srcHandler:  handlers.NewSourceHandler(ctx, sourceSvc),
	}
}

func (r *router) CreateSource(ctx *gin.Context) {
	r.srcHandler.PostSource(ctx)
}

func (r *router) DeleteSource(ctx *gin.Context, uriDigest string) {
	r.srcHandler.DeleteSourceByUriDigest(ctx, uriDigest)
}

func (r *router) GetSource(ctx *gin.Context, uriDigest string) {
	r.srcHandler.GetSourceByUriDigest(ctx, uriDigest)
}

func (r *router) GetSources(ctx *gin.Context) {
	r.srcHandler.GetSources(ctx)
}

func (r *router) PatchSource(ctx *gin.Context, uriDigest string) {
	r.srcHandler.PatchSourceByUriDigest(ctx, uriDigest)
}

func (r *router) GetPing(ctx *gin.Context) {
	message := r.pingHandler.GetPing(ctx)

	ctx.JSON(http.StatusOK, gin.H{"data": message})
}
