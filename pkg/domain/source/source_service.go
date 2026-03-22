package source

import (
	"context"
	"source-score/pkg/api"
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
	return svc.sourceRepo.PostSource(ctx, sourceInput)
}

func (svc *sourceService) PatchSourceByUriDigest(ctx context.Context, sourceInput *api.SourceInput, uriDigest string) error {
	return svc.sourceRepo.PatchSourceByUriDigest(ctx, sourceInput, uriDigest)
}
