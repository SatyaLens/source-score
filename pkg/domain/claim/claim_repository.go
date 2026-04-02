package claim

import (
	"context"
	"fmt"
	"log/slog"

	"source-score/pkg/api"
	"source-score/pkg/db/pgsql"
)

//go:generate go tool counterfeiter . ClaimRepository
type ClaimRepository interface {
	GetClaims(ctx context.Context) ([]api.Claim, error)
}

type claimRepository struct {
	client *pgsql.Client
}

func NewClaimRepository(ctx context.Context, client *pgsql.Client) ClaimRepository {
	return &claimRepository{client: client}
}

// GetClaims returns all claims from the DB
func (cr *claimRepository) GetClaims(ctx context.Context) ([]api.Claim, error) {
	var claims []api.Claim
	result := cr.client.FindAll(ctx, &claims)

	if result.Error != nil {
		return nil, result.Error
	}
	slog.InfoContext(ctx, fmt.Sprintf("returned %d claims", len(claims)))

	return claims, nil
}
