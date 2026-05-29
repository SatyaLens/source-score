package http

import (
	"context"
	"source-score/pkg/api"
	"source-score/pkg/conf"
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
	authHandler  *handlers.AuthHandler
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
		authHandler:  handlers.NewAuthHandler(conf.Cfg.JwtSecret),
	}
}

func (r *router) PostSource(ctx *gin.Context, params api.PostSourceParams) {
	r.srcHandler.PostSource(ctx)
}

func (r *router) DeleteSource(ctx *gin.Context, uriDigest string, params api.DeleteSourceParams) {
	r.srcHandler.DeleteSourceByUriDigest(ctx, uriDigest)
}

func (r *router) GetSource(ctx *gin.Context, uriDigest string, params api.GetSourceParams) {
	r.srcHandler.GetSourceByUriDigest(ctx, uriDigest)
}

func (r *router) GetSources(ctx *gin.Context, params api.GetSourcesParams) {
	r.srcHandler.GetSources(ctx)
}

func (r *router) PatchSource(ctx *gin.Context, uriDigest string, params api.PatchSourceParams) {
	r.srcHandler.PatchSourceByUriDigest(ctx, uriDigest)
}

func (r *router) GetClaims(ctx *gin.Context, params api.GetClaimsParams) {
	r.claimHandler.GetClaims(ctx, params)
}

func (r *router) PostClaim(ctx *gin.Context, params api.PostClaimParams) {
	r.claimHandler.PostClaim(ctx)
}

func (r *router) GetClaim(ctx *gin.Context, uriDigest string, params api.GetClaimParams) {
	r.claimHandler.GetClaimByUriDigest(ctx, uriDigest)
}

func (r *router) DeleteClaim(ctx *gin.Context, uriDigest string, params api.DeleteClaimParams) {
	r.claimHandler.DeleteClaimByUriDigest(ctx, uriDigest)
}

func (r *router) PatchClaim(ctx *gin.Context, claimDigest string, params api.PatchClaimParams) {
	r.claimHandler.PatchClaimByUriDigest(ctx, claimDigest)
}

func (r *router) VerifyAllClaims(ctx *gin.Context, params api.VerifyAllClaimsParams) {
	r.claimHandler.VerifyAllClaims(ctx)
}

func (r *router) VerifyClaim(ctx *gin.Context, claimDigest string, params api.VerifyClaimParams) {
	// TODO: remove if individual claim verification is not required
	// r.claimHandler.ValidateClaimByUriDigest(ctx, claimDigest)
}

func (r *router) PostProof(ctx *gin.Context, params api.PostProofParams) {
	r.proofHandler.PostProof(ctx)
}

func (r *router) DeleteProof(ctx *gin.Context, uriDigest string, params api.DeleteProofParams) {
	r.proofHandler.DeleteProofByUriDigest(ctx, uriDigest)
}

func (r *router) GetProof(ctx *gin.Context, uriDigest string, params api.GetProofParams) {
	r.proofHandler.GetProofByUriDigest(ctx, uriDigest)
}

func (r *router) GetProofs(ctx *gin.Context, params api.GetProofsParams) {
	r.proofHandler.GetProofs(ctx)
}

func (r *router) PatchProof(ctx *gin.Context, uriDigest string, params api.PatchProofParams) {
	r.proofHandler.PatchProofByUriDigest(ctx, uriDigest)
}

func (r *router) UpdateAllScores(ctx *gin.Context, params api.UpdateAllScoresParams) {
	r.srcHandler.UpdateAllScores(ctx)
}

func (r *router) GetClaimsBySourceDigest(ctx *gin.Context, sourceDigest string, params api.GetClaimsBySourceDigestParams) {
	r.claimHandler.GetClaimsBySourceDigest(ctx, sourceDigest)
}

func (r *router) GetProofsByClaimDigest(ctx *gin.Context, claimDigest string, params api.GetProofsByClaimDigestParams) {
	r.proofHandler.GetProofsByClaimDigest(ctx, claimDigest)
}

func (r *router) GetAuthToken(c *gin.Context, params api.GetAuthTokenParams) {
	r.authHandler.GetAuthToken(c)
}
