package handlers

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"source-score/pkg/api"
	"source-score/pkg/apperrors"
	"source-score/pkg/domain/source"
	"sync/atomic"

	"github.com/gin-gonic/gin"
)

type SourceHandler struct {
	sourceSvc source.SourceService
}

var (
	scoreJobRunning atomic.Bool
)

func NewSourceHandler(ctx context.Context, sourceSvc source.SourceService) *SourceHandler {
	return &SourceHandler{
		sourceSvc: sourceSvc,
	}
}

func (sh *SourceHandler) DeleteSourceByUriDigest(ctx *gin.Context, uriDigest string) {
	err := sh.sourceSvc.DeleteSourceByUriDigest(ctx, uriDigest)
	if err != nil {
		switch {
		case errors.Is(err, apperrors.ErrSourceNotFound):
			ctx.JSON(
				http.StatusNotFound,
				gin.H{"error": err.Error()},
			)
		default:
			slog.Error("failed to delete source", "error", err, "uriDigest", uriDigest)
			ctx.JSON(
				http.StatusInternalServerError,
				gin.H{"error": "internal server error"},
			)
		}
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (sh *SourceHandler) GetSources(ctx *gin.Context) {
	sources, err := sh.sourceSvc.GetSources(ctx)
	if err != nil {
		slog.Error("failed to get sources", "error", err)
		ctx.JSON(
			http.StatusInternalServerError,
			gin.H{"error": "internal server error"},
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
	if err != nil {
		switch {
		case errors.Is(err, apperrors.ErrSourceNotFound):
			ctx.JSON(
				http.StatusNotFound,
				gin.H{"error": err.Error()},
			)
		default:
			slog.Error("failed to get source", "error", err, "uriDigest", uriDigest)
			ctx.JSON(
				http.StatusInternalServerError,
				gin.H{"error": "internal server error"},
			)
		}
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
		slog.Error("failed to bind source input", "error", err)
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
			slog.Error("failed to create source", "error", err)
			ctx.JSON(
				http.StatusInternalServerError,
				gin.H{"error": "internal server error"},
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
		slog.Error("failed to bind source patch input", "error", err, "uriDigest", uriDigest)
		ctx.JSON(
			http.StatusBadRequest,
			gin.H{"error": err.Error()},
		)
		return
	}

	if err = sh.sourceSvc.PatchSourceByUriDigest(ctx, sourceInput, uriDigest); err != nil {
		switch {
		case errors.Is(err, apperrors.ErrInvalidSource):
			ctx.JSON(
				http.StatusBadRequest,
				gin.H{"error": err.Error()},
			)
		case errors.Is(err, apperrors.ErrSourceNotFound):
			ctx.JSON(
				http.StatusNotFound,
				gin.H{"error": err.Error()},
			)
		default:
			slog.Error("failed to patch source", "error", err, "uriDigest", uriDigest)
			ctx.JSON(
				http.StatusInternalServerError,
				gin.H{"error": "internal server error"},
			)
		}
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (sh *SourceHandler) UpdateAllScores(ctx *gin.Context) {
	if scoreJobRunning.CompareAndSwap(false, true) {
		go func(c *gin.Context) {
			defer scoreJobRunning.Store(false)
			if err := sh.sourceSvc.UpdateAllScores(c); err != nil {
				slog.Error(fmt.Sprintf("source score update job failed with error: %v", err))
			}
		}(ctx.Copy())

		ctx.Status(http.StatusAccepted)
		return
	}

	ctx.Status(http.StatusConflict)
}
