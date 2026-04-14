package claim

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log/slog"
	"strings"

	"source-score/pkg/api"
	"source-score/pkg/db/pgsql"
)

//go:generate go tool counterfeiter . ClaimRepository
type ClaimRepository interface {
	GetClaims(ctx context.Context) ([]api.Claim, error)
	PostClaim(ctx context.Context, claimInput *api.ClaimInput) (string, error)
	GetClaimByUriDigest(ctx context.Context, uriDigest string) (*api.Claim, error)
	DeleteClaimByUriDigest(ctx context.Context, claim *api.Claim) error
	PatchClaimByUriDigest(ctx context.Context, claimInput *api.ClaimPatchInput, uriDigest string) error
	VerifyClaimByUriDigest(ctx context.Context, claimVerification *api.ClaimVerification, uriDigest string) error
	VerifyClaims(ctx context.Context, updatedClaims []api.Claim) error
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

// PostClaim creates a new claim record and returns the computed uriDigest
func (cr *claimRepository) PostClaim(ctx context.Context, claimInput *api.ClaimInput) (string, error) {
	hash := sha256.New()
	_, err := hash.Write([]byte(claimInput.Uri))
	if err != nil {
		return "", err
	}
	uriDigest := hex.EncodeToString(hash.Sum(nil))

	claim := &api.Claim{
		SourceUriDigest: claimInput.SourceUriDigest,
		Summary:         claimInput.Summary,
		Title:           claimInput.Title,
		Uri:             claimInput.Uri,
		UriDigest:       uriDigest,
		Checked:         false,
		Validity:        false,
	}

	result := cr.client.Create(ctx, claim)
	slog.InfoContext(ctx, fmt.Sprintf("%d rows affected\n", result.RowsAffected))

	if result.Error != nil {
		return "", result.Error
	}

	return uriDigest, nil
}

// GetClaimByUriDigest returns a single claim by its uri digest
func (cr *claimRepository) GetClaimByUriDigest(ctx context.Context, uriDigest string) (*api.Claim, error) {
	claim := &api.Claim{}
	claim.UriDigest = uriDigest
	result := cr.client.FindFirst(ctx, claim)

	if result.Error != nil {
		return nil, result.Error
	}

	return claim, nil
}

// DeleteClaimByUriDigest deletes the provided claim record
func (cr *claimRepository) DeleteClaimByUriDigest(ctx context.Context, claim *api.Claim) error {
	result := cr.client.Delete(ctx, claim)
	slog.InfoContext(ctx, fmt.Sprintf("%d rows affected\n", result.RowsAffected))
	return result.Error
}

// PatchClaimByUriDigest updates claim fields
func (cr *claimRepository) PatchClaimByUriDigest(ctx context.Context, claimInput *api.ClaimPatchInput, uriDigest string) error {
	claim := &api.Claim{}
	claim.UriDigest = uriDigest

	result := cr.client.FindFirst(ctx, claim)
	if result.Error != nil {
		return result.Error
	}

	if claimInput.Summary != nil {
		claim.Summary = *claimInput.Summary
	}
	if claimInput.Title != nil {
		claim.Title = *claimInput.Title
	}

	result = cr.client.Update(ctx, claim)
	slog.InfoContext(ctx, fmt.Sprintf("%d rows affected\n", result.RowsAffected))

	return result.Error
}

func (cr *claimRepository) VerifyClaimByUriDigest(ctx context.Context, claimVerification *api.ClaimVerification, uriDigest string) error {
	claim := &api.Claim{}
	claim.UriDigest = uriDigest

	result := cr.client.FindFirst(ctx, claim)
	if result.Error != nil {
		return result.Error
	}

	claim.Checked = true
	claim.Validity = *claimVerification.Validity

	result = cr.client.Update(ctx, claim)
	slog.InfoContext(ctx, fmt.Sprintf("%d rows affected\n", result.RowsAffected))

	return result.Error
}

func (cr *claimRepository) VerifyClaims(ctx context.Context, updatedClaims []api.Claim) error {
	var args []any
	var query strings.Builder
	claimDigests := []string{}
	query.WriteString("UPDATE claims SET checked = true, validity = CASE uri_digest")

	for _, claim := range updatedClaims {
		query.WriteString(" WHEN ? THEN ?")
		args = append(args, claim.UriDigest, claim.Validity)
		claimDigests = append(claimDigests, claim.UriDigest)
	}

	query.WriteString(" END WHERE uri_digest IN ?")
	args = append(args, claimDigests)

	result := cr.client.DB.Exec(query.String(), args...)
	if result.Error != nil {
		return result.Error
	}

	fmt.Printf("rows updated: %d\n", result.RowsAffected)
	return nil
}
