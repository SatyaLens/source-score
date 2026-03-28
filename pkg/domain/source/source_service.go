package source

import (
	"context"
	"errors"
	"fmt"
	"source-score/pkg/api"
	"source-score/pkg/apperrors"
	"strings"

	"source-score/pkg/helpers"

	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

var (
	validate = validator.New()
)

type SourceService interface {
	DeleteSourceByUriDigest(ctx context.Context, uriDigest string) error
	GetSources(ctx context.Context) ([]api.Source, error)
	GetSourceByUriDigest(ctx context.Context, uriDigest string) (*api.Source, error)
	PostSource(ctx context.Context, sourceInput *api.SourceInput) (string, error)
	PatchSourceByUriDigest(ctx context.Context, sourceInput *api.SourcePatchInput, uriDigest string) error
}

type sourceService struct {
	sourceRepo SourceRepository
}

func init() {
	if err := validate.RegisterValidation("nonempty", helpers.ValidateNonEmpty); err != nil {
		panic(fmt.Sprintf("failed to register nonempty validator with error: %v", err))
	}
	if err := validate.RegisterValidation("httpsurl", helpers.ValidateHttpsURL); err != nil {
		panic(fmt.Sprintf("failed to register httpsurl validator with error: %v", err))
	}
	if err := validate.RegisterValidation("nospace", helpers.ValidateNoSpace); err != nil {
		panic(fmt.Sprintf("failed to register nospace validator with error: %v", err))
	}
}

func NewSourceService(ctx context.Context, sourceRepo SourceRepository) SourceService {
	return &sourceService{
		sourceRepo: sourceRepo,
	}
}

func (svc *sourceService) DeleteSourceByUriDigest(ctx context.Context, uriDigest string) error {
	source, err := svc.GetSourceByUriDigest(ctx, uriDigest)
	if err != nil {
		return err
	}

	return svc.sourceRepo.DeleteSourceByUriDigest(ctx, source)
}

func (svc *sourceService) GetSources(ctx context.Context) ([]api.Source, error) {
	return svc.sourceRepo.GetSources(ctx)
}

func (svc *sourceService) GetSourceByUriDigest(ctx context.Context, uriDigest string) (*api.Source, error) {
	return svc.sourceRepo.GetSourceByUriDigest(ctx, uriDigest)
}

func (svc *sourceService) PostSource(ctx context.Context, sourceInput *api.SourceInput) (string, error) {
	err := validate.Struct(sourceInput)
	if err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			return "", fmt.Errorf("%w: %s", apperrors.ErrValidationLogic, err.Error())
		}
		combinedErrs := ""
		for _, e := range err.(validator.ValidationErrors) {
			combinedErrs = fmt.Sprintf(
				"%s\n%s validation failed for value %v with error %s",
				combinedErrs, e.Field(), e.Value(), e.Tag(),
			)
		}
		combinedErrs = strings.TrimSpace(combinedErrs)
		return "", fmt.Errorf("%w: %s", apperrors.ErrInvalidSource, combinedErrs)
	}
	return svc.sourceRepo.PostSource(ctx, sourceInput)
}

func (svc *sourceService) PatchSourceByUriDigest(ctx context.Context, sourceInput *api.SourcePatchInput, uriDigest string) error {
	err := validate.Struct(sourceInput)
	if err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			return fmt.Errorf("%w: %s", apperrors.ErrValidationLogic, err.Error())
		}
		combinedErrs := ""
		for _, e := range err.(validator.ValidationErrors) {
			combinedErrs = fmt.Sprintf(
				"%s\n%s validation failed for value %v with error %s",
				combinedErrs, e.Field(), e.Value(), e.Tag(),
			)
		}
		combinedErrs = strings.TrimSpace(combinedErrs)
		return fmt.Errorf("%w: %s", apperrors.ErrInvalidSource, combinedErrs)
	}

	err = svc.sourceRepo.PatchSourceByUriDigest(ctx, sourceInput, uriDigest)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("%w: %s", apperrors.ErrSourceNotFound, err.Error())
	}
	return err
}
