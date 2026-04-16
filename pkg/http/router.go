package http

import (
	"context"
	"net/http"
	"source-score/pkg/domain/claim"
	"source-score/pkg/domain/proof"
	"source-score/pkg/domain/source"
	"source-score/pkg/handlers"

	"github.com/gin-gonic/gin"
)

type router struct {
	pingHandler  *handlers.PingHandler
	srcHandler   *handlers.SourceHandler
	claimHandler *handlers.ClaimHandler
	proofHandler *handlers.ProofHandler
}

func NewRouter(
	ctx context.Context,
	sourceSvc source.SourceService,
	claimSvc claim.ClaimService,
	proofSvc proof.ProofService,
) *router {
	return &router{
		pingHandler:  handlers.NewPingHandler(),
		srcHandler:   handlers.NewSourceHandler(ctx, sourceSvc),
		claimHandler: handlers.NewClaimHandler(ctx, claimSvc),
		proofHandler: handlers.NewProofHandler(ctx, proofSvc),
	}
}

func (r *router) PostSource(ctx *gin.Context) {
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

func (r *router) GetClaims(ctx *gin.Context) {
	r.claimHandler.GetClaims(ctx)
}

func (r *router) PostClaim(ctx *gin.Context) {
	r.claimHandler.PostClaim(ctx)
}

func (r *router) GetClaim(ctx *gin.Context, uriDigest string) {
	r.claimHandler.GetClaimByUriDigest(ctx, uriDigest)
}

func (r *router) GetPing(ctx *gin.Context) {
	message := r.pingHandler.GetPing(ctx)

	ctx.JSON(http.StatusOK, gin.H{"data": message})
}

func (r *router) DeleteClaim(ctx *gin.Context, uriDigest string) {
	r.claimHandler.DeleteClaimByUriDigest(ctx, uriDigest)
}

func (r *router) PatchClaim(ctx *gin.Context, claimDigest string) {
	r.claimHandler.PatchClaimByUriDigest(ctx, claimDigest)
}

func (r *router) VerifyAllClaims(ctx *gin.Context) {
	r.claimHandler.VerifyAllClaims(ctx)
}

func (r *router) VerifyClaim(ctx *gin.Context, claimDigest string) {
	// TODO: remove if individual claim verification is not required
	// r.claimHandler.ValidateClaimByUriDigest(ctx, claimDigest)
}

func (r *router) PostProof(ctx *gin.Context) {
	r.proofHandler.PostProof(ctx)
}

func (r *router) DeleteProof(ctx *gin.Context, uriDigest string) {
	r.proofHandler.DeleteProofByUriDigest(ctx, uriDigest)
}

func (r *router) GetProof(ctx *gin.Context, uriDigest string) {
	r.proofHandler.GetProofByUriDigest(ctx, uriDigest)
}

func (r *router) GetProofs(ctx *gin.Context) {
	r.proofHandler.GetProofs(ctx)
}

func (r *router) PatchProof(ctx *gin.Context, uriDigest string) {
	r.proofHandler.PatchProofByUriDigest(ctx, uriDigest)
}

func (r *router) UpdateAllScores(ctx *gin.Context) {
}
