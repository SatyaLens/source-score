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
	uriDigest := ctx.Param("uriDigest")

	// TODO: add basic validation for uriDigest
	if err := helpers.ValidateUriDigest(uriDigest); err != nil {
		ctx.JSON(
			http.StatusBadRequest,
			gin.H{"error": err.Error()},
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
	uriDigest := ctx.Param("uriDigest")

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

func (sh *SourceHandler) PatchSourceByUriDigest(ctx *gin.Context) {
	uriDigest := ctx.Param("uriDigest")

	// TODO: add basic validation for uriDigest
	if err := helpers.ValidateUriDigest(uriDigest); err != nil {
		ctx.JSON(
			http.StatusBadRequest,
			gin.H{"error": err.Error()},
		)
		return
	}

	sourceInput := &api.SourceInput{}

	err := ctx.ShouldBindJSON(sourceInput)
	if err != nil {
		ctx.JSON(
			http.StatusBadRequest,
			gin.H{"error": err.Error()},
		)
		return
	}

	if err = sh.sourceSvc.PatchSourceByUriDigest(ctx, sourceInput, uriDigest); err != nil {
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
