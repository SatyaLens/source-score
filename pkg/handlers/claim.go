package handlers

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"source-score/pkg/api"
	"source-score/pkg/apperrors"
	"source-score/pkg/domain/claim"
	"sync/atomic"

	"github.com/gin-gonic/gin"
)

type ClaimHandler struct {
	claimSvc claim.ClaimService
}

var (
	verificationJobRunning atomic.Bool
)

func NewClaimHandler(ctx context.Context, claimSvc claim.ClaimService) *ClaimHandler {
	return &ClaimHandler{claimSvc: claimSvc}
}

func (ch *ClaimHandler) GetClaims(ctx *gin.Context) {
	claims, err := ch.claimSvc.GetClaims(ctx)
	if err != nil {
		slog.Error("failed to get claims", "error", err)
		ctx.JSON(
			http.StatusInternalServerError,
			gin.H{"error": "internal server error"},
		)
		return
	}

	ctx.JSON(
		http.StatusOK,
		claims,
	)
}

func (ch *ClaimHandler) PostClaim(ctx *gin.Context) {
	claimInput := &api.ClaimInput{}
	if err := ctx.ShouldBindJSON(claimInput); err != nil {
		slog.Error("failed to bind claim input", "error", err)
		ctx.JSON(
			http.StatusBadRequest,
			gin.H{"error": err.Error()},
		)
		return
	}

	digest, err := ch.claimSvc.PostClaim(ctx, claimInput)
	if err != nil {
		switch {
		case errors.Is(err, apperrors.ErrInvalidClaim):
			ctx.JSON(
				http.StatusBadRequest,
				gin.H{"error": err.Error()},
			)
		default:
			slog.Error("failed to create claim", "error", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"uriDigest": digest})
}

func (ch *ClaimHandler) GetClaimByUriDigest(ctx *gin.Context, uriDigest string) {
	claim, err := ch.claimSvc.GetClaimByUriDigest(ctx, uriDigest)
	if err != nil {
		switch {
		case errors.Is(err, apperrors.ErrClaimNotFound):
			ctx.JSON(
				http.StatusNotFound,
				gin.H{"error": err.Error()},
			)
		default:
			slog.Error("failed to get claim", "error", err, "uriDigest", uriDigest)
			ctx.JSON(
				http.StatusInternalServerError,
				gin.H{"error": "internal server error"},
			)
		}
		return
	}

	ctx.JSON(
		http.StatusOK,
		claim,
	)
}

func (ch *ClaimHandler) DeleteClaimByUriDigest(ctx *gin.Context, uriDigest string) {
	err := ch.claimSvc.DeleteClaimByUriDigest(ctx, uriDigest)
	if err != nil {
		switch {
		case errors.Is(err, apperrors.ErrClaimNotFound):
			ctx.JSON(
				http.StatusNotFound,
				gin.H{"error": err.Error()},
			)
		default:
			slog.Error("failed to delete claim", "error", err, "uriDigest", uriDigest)
			ctx.JSON(
				http.StatusInternalServerError,
				gin.H{"error": "internal server error"},
			)
		}
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (ch *ClaimHandler) PatchClaimByUriDigest(ctx *gin.Context, uriDigest string) {
	claimInput := &api.ClaimPatchInput{}
	if err := ctx.ShouldBindJSON(claimInput); err != nil {
		slog.Error("failed to bind claim patch input", "error", err, "uriDigest", uriDigest)
		ctx.JSON(
			http.StatusBadRequest,
			gin.H{"error": err.Error()},
		)
		return
	}

	if err := ch.claimSvc.PatchClaimByUriDigest(ctx, claimInput, uriDigest); err != nil {
		switch {
		case errors.Is(err, apperrors.ErrInvalidClaim):
			ctx.JSON(
				http.StatusBadRequest,
				gin.H{"error": err.Error()},
			)
		case errors.Is(err, apperrors.ErrClaimNotFound):
			ctx.JSON(
				http.StatusNotFound,
				gin.H{"error": err.Error()},
			)
		default:
			slog.Error("failed to patch claim", "error", err, "uriDigest", uriDigest)
			ctx.JSON(
				http.StatusInternalServerError,
				gin.H{"error": "internal server error"},
			)
		}
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (ch *ClaimHandler) ValidateClaimByUriDigest(ctx *gin.Context, uriDigest string) {
	claimVerification := &api.ClaimVerification{}
	if err := ctx.ShouldBindJSON(claimVerification); err != nil {
		slog.Error("failed to bind claim verification input", "error", err, "uriDigest", uriDigest)
		ctx.JSON(
			http.StatusBadRequest,
			gin.H{"error": err.Error()},
		)
		return
	}

	if err := ch.claimSvc.VerifyClaimByUriDigest(ctx, claimVerification, uriDigest); err != nil {
		switch {
		case errors.Is(err, apperrors.ErrInvalidClaimVerification):
			ctx.JSON(
				http.StatusBadRequest,
				gin.H{"error": err.Error()},
			)
		case errors.Is(err, apperrors.ErrClaimNotFound):
			ctx.JSON(
				http.StatusNotFound,
				gin.H{"error": err.Error()},
			)
		default:
			slog.Error("failed to verify claim", "error", err, "uriDigest", uriDigest)
			ctx.JSON(
				http.StatusInternalServerError,
				gin.H{"error": "internal server error"},
			)
		}
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (ch *ClaimHandler) VerifyAllClaims(ctx *gin.Context) {
	if verificationJobRunning.CompareAndSwap(false, true) {
		go func(c *gin.Context) {
			defer verificationJobRunning.Store(false)
			if err := ch.claimSvc.VerifyAllClaims(c); err != nil {
				slog.Error(fmt.Sprintf("claim verification job failed with error: %v", err))
			}
		}(ctx.Copy())

		ctx.Status(http.StatusAccepted)
		return
	}

	ctx.Status(http.StatusConflict)
}
