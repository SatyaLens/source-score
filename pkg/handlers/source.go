package handlers

import (
	"context"
	"errors"
	"net/http"
	"source-score/pkg/api"
	"source-score/pkg/apperrors"
	"source-score/pkg/domain/source"

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

func (sh *SourceHandler) DeleteSourceByUriDigest(ctx *gin.Context, uriDigest string) {
	err := sh.sourceSvc.DeleteSourceByUriDigest(ctx, uriDigest)
	// TODO: add proper error wrapping logic
	if err != nil {
		ctx.JSON(
			http.StatusNotFound,
			gin.H{"error": err.Error()},
		)
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (sh *SourceHandler) GetSources(ctx *gin.Context) {
	sources, err := sh.sourceSvc.GetSources(ctx)
	if err != nil {
		ctx.JSON(
			http.StatusInternalServerError,
			gin.H{"error": err.Error()},
		)
		return
	}

	ctx.JSON(
		http.StatusOK,
		sources,
	)
}

func (sh *SourceHandler) GetSourceByUriDigest(ctx *gin.Context, uriDigest string) {
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
	err := ctx.ShouldBindJSON(sourceInput)
	if err != nil {
		ctx.JSON(
			http.StatusBadRequest,
			gin.H{"error": err.Error()},
		)
		return
	}

	digest, err := sh.sourceSvc.PostSource(ctx, sourceInput)
	if err != nil {
		switch {
		case errors.Is(err, apperrors.ErrInvalidSource):
			ctx.JSON(
				http.StatusBadRequest,
				gin.H{"error": err.Error()},
			)
		default:
			ctx.JSON(
				http.StatusInternalServerError,
				gin.H{"error": err.Error()},
			)
		}
		return
	}

	ctx.JSON(
		http.StatusCreated,
		api.CreateSourceResponse{UriDigest: digest},
	)
}

func (sh *SourceHandler) PatchSourceByUriDigest(ctx *gin.Context, uriDigest string) {
	sourceInput := &api.SourcePatchInput{}
	err := ctx.ShouldBindJSON(sourceInput)
	if err != nil {
		ctx.JSON(
			http.StatusBadRequest,
			gin.H{"error": err.Error()},
		)
		return
	}

	if err = sh.sourceSvc.PatchSourceByUriDigest(ctx, sourceInput, uriDigest); err != nil {
		switch {
		case errors.Is(err, apperrors.ErrInvalidSource):
			ctx.Status(http.StatusBadRequest)
		default:
			ctx.Status(http.StatusInternalServerError)
		}
		return
	}

	ctx.Status(http.StatusNoContent)
}
