package proof

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log/slog"

	"source-score/pkg/api"
	"source-score/pkg/db/pgsql"
)

//go:generate go tool counterfeiter . ProofRepository
type ProofRepository interface {
	GetProofs(ctx context.Context) ([]api.Proof, error)
	PostProof(ctx context.Context, proofInput *api.ProofInput) (string, error)
	GetProofByUriDigest(ctx context.Context, uriDigest string) (*api.Proof, error)
	DeleteProofByUriDigest(ctx context.Context, proof *api.Proof) error
	PatchProofByUriDigest(ctx context.Context, proofInput *api.ProofPatchInput, uriDigest string) error
}

type proofRepository struct {
	client *pgsql.Client
}

func NewProofRepository(ctx context.Context, client *pgsql.Client) ProofRepository {
	return &proofRepository{client: client}
}

// GetProofs returns all proofs from the DB
func (pr *proofRepository) GetProofs(ctx context.Context) ([]api.Proof, error) {
	var proofs []api.Proof
	result := pr.client.FindAll(ctx, &proofs)

	if result.Error != nil {
		return nil, result.Error
	}

	slog.InfoContext(ctx, fmt.Sprintf("returned %d proofs", len(proofs)))

	return proofs, nil
}

// PostProof creates a new proof record and returns the computed uriDigest
func (pr *proofRepository) PostProof(ctx context.Context, proofInput *api.ProofInput) (string, error) {
	hash := sha256.New()
	_, err := hash.Write([]byte(proofInput.Uri))
	if err != nil {
		return "", err
	}
	uriDigest := hex.EncodeToString(hash.Sum(nil))

	proof := &api.Proof{
		ClaimUriDigest: proofInput.ClaimUriDigest,
		ReviewedBy:     proofInput.ReviewedBy,
		SupportsClaim:  *proofInput.SupportsClaim,
		Uri:            proofInput.Uri,
		UriDigest:      uriDigest,
	}

	result := pr.client.Create(ctx, proof)
	slog.InfoContext(ctx, fmt.Sprintf("%d rows affected\n", result.RowsAffected))

	if result.Error != nil {
		return "", result.Error
	}

	return uriDigest, nil
}

// GetProofByUriDigest returns a single proof by its uri digest
func (pr *proofRepository) GetProofByUriDigest(ctx context.Context, uriDigest string) (*api.Proof, error) {
	proof := &api.Proof{}
	proof.UriDigest = uriDigest
	result := pr.client.FindFirst(ctx, proof)

	if result.Error != nil {
		return nil, result.Error
	}

	return proof, nil
}

// DeleteProofByUriDigest deletes the provided proof record
func (pr *proofRepository) DeleteProofByUriDigest(ctx context.Context, proof *api.Proof) error {
	result := pr.client.Delete(ctx, proof)
	slog.InfoContext(ctx, fmt.Sprintf("%d rows affected\n", result.RowsAffected))
	return result.Error
}

// PatchProofByUriDigest updates proof fields
func (pr *proofRepository) PatchProofByUriDigest(ctx context.Context, proofInput *api.ProofPatchInput, uriDigest string) error {
	proof := &api.Proof{}
	proof.UriDigest = uriDigest

	result := pr.client.FindFirst(ctx, proof)
	if result.Error != nil {
		return result.Error
	}

	// ProofPatchInput.ReviewedBy is required in the API spec, so update directly
	proof.ReviewedBy = proofInput.ReviewedBy

	result = pr.client.Update(ctx, proof)
	slog.InfoContext(ctx, fmt.Sprintf("%d rows affected\n", result.RowsAffected))

	return result.Error
}
