package handlers

import (
	"context"
	"net/http"
	"source-score/pkg/api"
	"source-score/pkg/domain/claim"

	"github.com/gin-gonic/gin"
)

type ClaimHandler struct {
	claimSvc claim.ClaimService
}

func NewClaimHandler(ctx context.Context, claimSvc claim.ClaimService) *ClaimHandler {
	return &ClaimHandler{claimSvc: claimSvc}
}

func (ch *ClaimHandler) GetClaims(ctx *gin.Context) {
	claims, err := ch.claimSvc.GetClaims(ctx)
	if err != nil {
		ctx.JSON(
			http.StatusInternalServerError,
			gin.H{"error": err.Error()},
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
		ctx.JSON(
			http.StatusBadRequest,
			gin.H{"error": err.Error()},
		)
		return
	}

	digest, err := ch.claimSvc.PostClaim(ctx, claimInput)
	if err != nil {
		switch {
		default:
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"uriDigest": digest})
}
