package claim

import (
	"context"
	"errors"
	"fmt"
	"source-score/pkg/api"
	"source-score/pkg/apperrors"

	"gorm.io/gorm"
)

//go:generate go tool counterfeiter . ClaimService
type ClaimService interface {
	GetClaims(ctx context.Context) ([]api.Claim, error)
	PostClaim(ctx context.Context, claimInput *api.ClaimInput) (string, error)
	GetClaimByUriDigest(ctx context.Context, uriDigest string) (*api.Claim, error)
	DeleteClaimByUriDigest(ctx context.Context, uriDigest string) error
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

func (svc *claimService) GetClaimByUriDigest(ctx context.Context, uriDigest string) (*api.Claim, error) {
	claim, err := svc.claimRepo.GetClaimByUriDigest(ctx, uriDigest)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("%w: %s", apperrors.ErrClaimNotFound, err.Error())
		}
		return nil, err
	}
	return claim, nil
}

func (svc *claimService) DeleteClaimByUriDigest(ctx context.Context, uriDigest string) error {
	claim, err := svc.GetClaimByUriDigest(ctx, uriDigest)
	if err != nil {
		return err
	}

	return svc.claimRepo.DeleteClaimByUriDigest(ctx, claim)
}

func (svc *claimService) PostClaim(ctx context.Context, claimInput *api.ClaimInput) (string, error) {
	return svc.claimRepo.PostClaim(ctx, claimInput)
}
