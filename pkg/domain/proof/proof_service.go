package proof

import (
	"context"
	"errors"
	"fmt"
	"source-score/pkg/api"
	"source-score/pkg/apperrors"
	"source-score/pkg/helpers"
	"strings"

	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

//go:generate go tool counterfeiter . ProofService
type ProofService interface {
	GetProofs(ctx context.Context) ([]api.Proof, error)
	PostProof(ctx context.Context, proofInput *api.ProofInput) (string, error)
	GetProofByUriDigest(ctx context.Context, uriDigest string) (*api.Proof, error)
	DeleteProofByUriDigest(ctx context.Context, uriDigest string) error
	PatchProofByUriDigest(ctx context.Context, proofInput *api.ProofPatchInput, uriDigest string) error
	GetProofsByClaims(ctx context.Context) (map[string][]api.Proof, error)
	GetProofsByClaimDigest(ctx context.Context, digest string) ([]api.Proof, error)
}

type proofService struct {
	proofRepo ProofRepository
}

var (
	proofValidate = validator.New()
)

func init() {
	if err := proofValidate.RegisterValidation("nonempty", helpers.ValidateNonEmpty); err != nil {
		panic(fmt.Sprintf("failed to register nonempty validator with error: %v", err))
	}
	if err := proofValidate.RegisterValidation("httpsurl", helpers.ValidateHttpsURL); err != nil {
		panic(fmt.Sprintf("failed to register httpsurl validator with error: %v", err))
	}
	if err := proofValidate.RegisterValidation("nospace", helpers.ValidateNoSpace); err != nil {
		panic(fmt.Sprintf("failed to register nospace validator with error: %v", err))
	}
}

func NewProofService(ctx context.Context, proofRepo ProofRepository) ProofService {
	return &proofService{proofRepo: proofRepo}
}

func (svc *proofService) GetProofs(ctx context.Context) ([]api.Proof, error) {
	return svc.proofRepo.GetProofs(ctx)
}

func (svc *proofService) GetProofByUriDigest(ctx context.Context, uriDigest string) (*api.Proof, error) {
	p, err := svc.proofRepo.GetProofByUriDigest(ctx, uriDigest)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("%w: %s", apperrors.ErrProofNotFound, err.Error())
		}
		return nil, err
	}
	return p, nil
}

func (svc *proofService) DeleteProofByUriDigest(ctx context.Context, uriDigest string) error {
	p, err := svc.GetProofByUriDigest(ctx, uriDigest)
	if err != nil {
		return err
	}
	return svc.proofRepo.DeleteProofByUriDigest(ctx, p)
}

func (svc *proofService) PatchProofByUriDigest(ctx context.Context, proofInput *api.ProofPatchInput, uriDigest string) error {
	err := proofValidate.Struct(proofInput)
	if err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			return fmt.Errorf("%w: %s", apperrors.ErrValidationLogic, err.Error())
		}
		combinedErrs := ""
		for _, e := range err.(validator.ValidationErrors) {
			combinedErrs = fmt.Sprintf("%s\n%s validation failed for value %v with error %s", combinedErrs, e.Field(), e.Value(), e.Tag())
		}
		combinedErrs = strings.TrimSpace(combinedErrs)
		return fmt.Errorf("%w: %s", apperrors.ErrInvalidProof, combinedErrs)
	}

	err = svc.proofRepo.PatchProofByUriDigest(ctx, proofInput, uriDigest)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("%w: %s", apperrors.ErrProofNotFound, err.Error())
	}
	return err
}

func (svc *proofService) PostProof(ctx context.Context, proofInput *api.ProofInput) (string, error) {
	err := proofValidate.Struct(proofInput)
	if err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			return "", fmt.Errorf("%w: %s", apperrors.ErrValidationLogic, err.Error())
		}
		combinedErrs := ""
		for _, e := range err.(validator.ValidationErrors) {
			combinedErrs = fmt.Sprintf("%s\n%s validation failed for value %v with error %s", combinedErrs, e.Field(), e.Value(), e.Tag())
		}
		combinedErrs = strings.TrimSpace(combinedErrs)
		return "", fmt.Errorf("%w: %s", apperrors.ErrInvalidProof, combinedErrs)
	}

	return svc.proofRepo.PostProof(ctx, proofInput)
}

func (svc *proofService) GetProofsByClaims(ctx context.Context) (map[string][]api.Proof, error) {
	allProofs, err := svc.GetProofs(ctx)
	if err != nil {
		return nil, err
	}

	claimsProofs := make(map[string][]api.Proof)
	for _, proof := range allProofs {
		if proofs, ok := claimsProofs[proof.ClaimUriDigest]; ok {
			claimsProofs[proof.ClaimUriDigest] = append(proofs, proof)
		} else {
			claimsProofs[proof.ClaimUriDigest] = []api.Proof{proof}
		}
	}

	return claimsProofs, nil
}

func (svc *proofService) GetProofsByClaimDigest(ctx context.Context, digest string) ([]api.Proof, error) {
	return svc.proofRepo.GetProofsByClaimDigest(ctx, digest)
}
