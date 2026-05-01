package handlers

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"source-score/pkg/api"
	"source-score/pkg/apperrors"
	"source-score/pkg/domain/proof"

	"github.com/gin-gonic/gin"
)

type ProofHandler struct {
	proofSvc proof.ProofService
}

func NewProofHandler(ctx context.Context, proofSvc proof.ProofService) *ProofHandler {
	return &ProofHandler{proofSvc: proofSvc}
}

func (ph *ProofHandler) GetProofs(ctx *gin.Context) {
	proofs, err := ph.proofSvc.GetProofs(ctx)
	if err != nil {
		slog.Error("failed to get proofs", "error", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	ctx.JSON(http.StatusOK, proofs)
}

func (ph *ProofHandler) PostProof(ctx *gin.Context) {
	proofInput := &api.ProofInput{}
	if err := ctx.ShouldBindJSON(proofInput); err != nil {
		slog.Error("failed to bind proof input", "error", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	digest, err := ph.proofSvc.PostProof(ctx, proofInput)
	if err != nil {
		switch {
		case errors.Is(err, apperrors.ErrInvalidProof):
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		default:
			slog.Error("failed to create proof", "error", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"uriDigest": digest})
}

func (ph *ProofHandler) GetProofByUriDigest(ctx *gin.Context, uriDigest string) {
	p, err := ph.proofSvc.GetProofByUriDigest(ctx, uriDigest)
	if err != nil {
		switch {
		case errors.Is(err, apperrors.ErrProofNotFound):
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		default:
			slog.Error("failed to get proof", "error", err, "uriDigest", uriDigest)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	ctx.JSON(http.StatusOK, p)
}

func (ph *ProofHandler) DeleteProofByUriDigest(ctx *gin.Context, uriDigest string) {
	err := ph.proofSvc.DeleteProofByUriDigest(ctx, uriDigest)
	if err != nil {
		switch {
		case errors.Is(err, apperrors.ErrProofNotFound):
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		default:
			slog.Error("failed to delete proof", "error", err, "uriDigest", uriDigest)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (ph *ProofHandler) PatchProofByUriDigest(ctx *gin.Context, uriDigest string) {
	proofInput := &api.ProofPatchInput{}
	if err := ctx.ShouldBindJSON(proofInput); err != nil {
		slog.Error("failed to bind proof patch input", "error", err, "uriDigest", uriDigest)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := ph.proofSvc.PatchProofByUriDigest(ctx, proofInput, uriDigest); err != nil {
		switch {
		case errors.Is(err, apperrors.ErrInvalidProof):
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		case errors.Is(err, apperrors.ErrProofNotFound):
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		default:
			slog.Error("failed to patch proof", "error", err, "uriDigest", uriDigest)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (ph *ProofHandler) GetProofsByClaimDigest(ctx *gin.Context, claimDigest string) {
	proofs, err := ph.proofSvc.GetProofsByClaimDigest(ctx, claimDigest)
	if err != nil {
		slog.Error("failed to get proofs by claim digest", "error", err, "claimDigest", claimDigest)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	ctx.JSON(http.StatusOK, proofs)
}
