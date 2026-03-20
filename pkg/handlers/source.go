package handlers

import (
	"context"
	"net/http"
	"source-score/pkg/api"
	"source-score/pkg/domain/source"
	"source-score/pkg/helpers"

	"github.com/gin-gonic/gin"
)

type SourceHandler struct {
	sourceSvc source.SourceService
}

func NewSourceHandler(ctx context.Context, sourceSvc source.SourceService) *SourceHandler {
	return &SourceHandler{
		sourceSvc: sourceSvc,
	}
}

func (sh *SourceHandler) DeleteSourceByUriDigest(ctx *gin.Context) {
	uriDigest := ctx.Param("uri_digest")

	// TODO: add basic validation for uriDigest
	if uriDigest == "" {
		ctx.JSON(
			http.StatusBadRequest,
			gin.H{"error": "invalid uri digest"},
		)
		return
	}

	err := sh.sourceSvc.DeleteSourceByUriDigest(ctx, uriDigest)
	// TODO: add proper error wrapping logic
	if err != nil {
		ctx.JSON(
			http.StatusNotFound,
			gin.H{"error": err.Error()},
		)
		return
	}

	ctx.JSON(
		http.StatusNoContent,
		gin.H{},
	)
}

func (sh *SourceHandler) GetSourceByUriDigest(ctx *gin.Context) {
	uriDigest := ctx.Param("uri_digest")

	// TODO: add basic validation for uriDigest
	if err := helpers.ValidateUriDigest(uriDigest); err != nil {
		ctx.JSON(
			http.StatusBadRequest,
			gin.H{"error": err.Error()},
		)
		return
	}

	source, err := sh.sourceSvc.GetSourceByUriDigest(ctx, uriDigest)
	// TODO: add proper error wrapping logic
	if err != nil {
		ctx.JSON(
			http.StatusNotFound,
			gin.H{"error": err.Error()},
		)
		return
	}

	ctx.JSON(
		http.StatusOK,
		source,
	)
}

func (sh *SourceHandler) PostSource(ctx *gin.Context) {
	// var sourceInput *api.SourceInput
	sourceInput := &api.SourceInput{}

	// TODO: add basic validation for source input
	err := ctx.ShouldBindJSON(sourceInput)
	if err != nil {
		ctx.JSON(
			http.StatusBadRequest,
			gin.H{"error": err.Error()},
		)
		return
	}

	err = sh.sourceSvc.PostSource(ctx, sourceInput)
	if err != nil {
		ctx.JSON(
			http.StatusInternalServerError,
			gin.H{"error": err.Error()},
		)
		return
	}

	ctx.JSON(
		http.StatusCreated,
		gin.H{},
	)
}

func (sh *SourceHandler) PutSourceByUriDigest(ctx *gin.Context) {
	uriDigest := ctx.Param("uri_digest")

	// TODO: add basic validation for uriDigest
	if err := helpers.ValidateUriDigest(uriDigest); err != nil {
		ctx.JSON(
			http.StatusBadRequest,
			gin.H{"error": err.Error()},
		)
		return
	}

	var sourceInput *api.SourceInput

	err := ctx.ShouldBindJSON(sourceInput)
	if err != nil {
		ctx.JSON(
			http.StatusBadRequest,
			gin.H{"error": err.Error()},
		)
		return
	}

	err = sh.sourceSvc.PutSourceByUriDigest(ctx, sourceInput, uriDigest)

	if err = sh.sourceSvc.PutSourceByUriDigest(ctx, sourceInput, uriDigest); err != nil {
		ctx.JSON(
			http.StatusNotFound,
			gin.H{"error": err.Error()},
		)
		return
	}

	ctx.JSON(
		http.StatusOK,
		gin.H{},
	)
}
