package api

import (
	"log"
	"log/slog"
	"net/http"
	"source-score/pkg/handlers"

	"github.com/gin-gonic/gin"
)

type router struct {
	pingHandler *handlers.PingHandler
}

func NewRouter() *router {
	return &router{
		pingHandler: handlers.NewPingHandler(),
	}
}

func (r *router) CreateSource(ctx *gin.Context) {
	body := SourceInput{}
	// using BindJson method to serialize body with struct
	if err := ctx.BindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		_ = ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
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
	// requestId := ctx.GetHeader(helpers.RequestIdHeader)
	// if requestId == "" {
	// 	requestId = uuid.NewString()
	// }
	// TODO :: log using a context that has request header and verify request ID is printed
	slog.InfoContext(ctx, "printing using slog")
	message := r.pingHandler.GetPing()

	ctx.JSON(http.StatusOK, gin.H{"data": message})
}
