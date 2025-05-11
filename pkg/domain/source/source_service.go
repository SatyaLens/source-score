package source

import (
	"context"
	"source-score/pkg/api"
	"source-score/pkg/db/cnpg"
	"source-score/pkg/helpers"
)

type SourceService struct {
	repo *SourceRepository
}

func NewSourceService(ctx context.Context, client *cnpg.Client) *SourceService {
	return &SourceService{
		repo: NewSourceRepository(ctx, client),
	}
}

func (svc *SourceService) DeleteSourceByUriDigest(ctx context.Context, uri string) error {
	uriDigest := helpers.GetSHA256Hash(uri)

	return svc.repo.DeleteSource(ctx, &api.Source{
		Uri: uri,
		UriDigest: uriDigest,
	})
}

func (svc *SourceService) GetSourceByUriDigest(ctx context.Context, uri string) (*api.Source, error) {
	source := &api.Source{
		Uri: uri,
		UriDigest: helpers.GetSHA256Hash(uri),
	}

	result := svc.repo.client.FindFirst(ctx, source)
	if result.Error != nil {
		return nil, result.Error
	}

	return source, nil
}
