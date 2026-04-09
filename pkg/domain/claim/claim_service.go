package claim

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

//go:generate go tool counterfeiter . ClaimService
type ClaimService interface {
	GetClaims(ctx context.Context) ([]api.Claim, error)
	PostClaim(ctx context.Context, claimInput *api.ClaimInput) (string, error)
	GetClaimByUriDigest(ctx context.Context, uriDigest string) (*api.Claim, error)
	DeleteClaimByUriDigest(ctx context.Context, uriDigest string) error
	PatchClaimByUriDigest(ctx context.Context, claimInput *api.ClaimPatchInput, uriDigest string) error
}

type claimService struct {
	claimRepo ClaimRepository
}

var (
	claimValidate = validator.New()
)

func init() {
	if err := claimValidate.RegisterValidation("nonempty", helpers.ValidateNonEmpty); err != nil {
		panic(fmt.Sprintf("failed to register nonempty validator with error: %v", err))
	}
	if err := claimValidate.RegisterValidation("httpsurl", helpers.ValidateHttpsURL); err != nil {
		panic(fmt.Sprintf("failed to register httpsurl validator with error: %v", err))
	}
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

func (svc *claimService) PatchClaimByUriDigest(ctx context.Context, claimInput *api.ClaimPatchInput, uriDigest string) error {
	err := claimValidate.Struct(claimInput)
	if err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			return fmt.Errorf("%w: %s", apperrors.ErrValidationLogic, err.Error())
		}
		combinedErrs := ""
		for _, e := range err.(validator.ValidationErrors) {
			combinedErrs = fmt.Sprintf(
				"%s\n%s validation failed for value %v with error %s", combinedErrs, e.Field(), e.Value(), e.Tag(),
			)
		}
		combinedErrs = strings.TrimSpace(combinedErrs)
		return fmt.Errorf("%w: %s", apperrors.ErrInvalidClaim, combinedErrs)
	}

	err = svc.claimRepo.PatchClaimByUriDigest(ctx, claimInput, uriDigest)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("%w: %s", apperrors.ErrClaimNotFound, err.Error())
	}
	return err
}

func (svc *claimService) PostClaim(ctx context.Context, claimInput *api.ClaimInput) (string, error) {
	err := claimValidate.Struct(claimInput)
	if err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			return "", fmt.Errorf("%w: %s", apperrors.ErrValidationLogic, err.Error())
		}
		combinedErrs := ""
		for _, e := range err.(validator.ValidationErrors) {
			combinedErrs = fmt.Sprintf(
				"%s\n%s validation failed for value %v with error %s", combinedErrs, e.Field(), e.Value(), e.Tag(),
			)
		}
		combinedErrs = strings.TrimSpace(combinedErrs)
		return "", fmt.Errorf("%w: %s", apperrors.ErrInvalidClaim, combinedErrs)
	}

	return svc.claimRepo.PostClaim(ctx, claimInput)
}
