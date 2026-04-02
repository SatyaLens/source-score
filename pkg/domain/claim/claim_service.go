package claim

import (
	"context"
	"source-score/pkg/api"
)

//go:generate go tool counterfeiter . ClaimService
type ClaimService interface {
	GetClaims(ctx context.Context) ([]api.Claim, error)
	PostClaim(ctx context.Context, claimInput *api.ClaimInput) (string, error)
}

type claimService struct {
	claimRepo ClaimRepository
}

func NewClaimService(ctx context.Context, claimRepo ClaimRepository) ClaimService {
	return &claimService{claimRepo: claimRepo}
}

func (svc *claimService) GetClaims(ctx context.Context) ([]api.Claim, error) {
	return svc.claimRepo.GetClaims(ctx)
}

func (svc *claimService) PostClaim(ctx context.Context, claimInput *api.ClaimInput) (string, error) {
	return svc.claimRepo.PostClaim(ctx, claimInput)
}
