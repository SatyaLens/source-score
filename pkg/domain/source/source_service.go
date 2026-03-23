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
)

var (
	validate = validator.New()
)

type SourceService interface {
	DeleteSourceByUriDigest(ctx context.Context, uriDigest string) error
	GetSourceByUriDigest(ctx context.Context, uriDigest string) (*api.Source, error)
	PostSource(ctx context.Context, sourceInput *api.SourceInput) (string, error)
	PatchSourceByUriDigest(ctx context.Context, sourceInput *api.SourceInput, uriDigest string) error
}

type sourceService struct {
	sourceRepo SourceRepository
}

func init() {
	validate.RegisterValidation("nonempty", helpers.ValidateNonEmpty)
	validate.RegisterValidation("httpsurl", helpers.ValidateHttpsURL)
	validate.RegisterValidation("nospace", helpers.ValidateNoSpace)
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

func (svc *sourceService) GetSourceByUriDigest(ctx context.Context, uriDigest string) (*api.Source, error) {
	source, err := svc.sourceRepo.GetSourceByUriDigest(ctx, uriDigest)
	if err != nil {
		return nil, err
	}

	return source, nil
}

func (svc *sourceService) PostSource(ctx context.Context, sourceInput *api.SourceInput) (string, error) {
	err := validate.Struct(sourceInput)
	if err != nil {
		if errors.Is(err, &validator.InvalidValidationError{}) {
			return "", fmt.Errorf("%w: %s", apperrors.ValidationLogic, err.Error())
		}
		combinedErrs := ""
		for _, e := range err.(validator.ValidationErrors) {
			combinedErrs = fmt.Sprintf(
				"%s\n%s validation failed for value %s with error %s",
				combinedErrs, e.Field(), e.Value(), e.Tag(),
			)
		}
		combinedErrs = strings.TrimSpace(combinedErrs)
		return "", fmt.Errorf("%w: %s", apperrors.InvalidSource, combinedErrs)
	}
	return svc.sourceRepo.PostSource(ctx, sourceInput)
}

func (svc *sourceService) PatchSourceByUriDigest(ctx context.Context, sourceInput *api.SourceInput, uriDigest string) error {
	return svc.sourceRepo.PatchSourceByUriDigest(ctx, sourceInput, uriDigest)
}
