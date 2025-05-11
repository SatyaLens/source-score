package source

import (
	"context"
	"fmt"
	"log/slog"

	"source-score/pkg/api"
	"source-score/pkg/db/cnpg"
	"source-score/pkg/helpers"
)

type SourceRepository struct {
	client *cnpg.Client
}

func NewSourceRepository(ctx context.Context, client *cnpg.Client) *SourceRepository {
	return &SourceRepository{
		client: client,
	}
}

func (sr *SourceRepository) DeleteSource(ctx context.Context, source *api.Source) error {
	result := sr.client.Delete(ctx, source)

	if result.Error == nil {
		slog.InfoContext(
			ctx,
			fmt.Sprintf("%d rows affected\n", result.RowsAffected),
		)
	}

	return result.Error
}

func (sr *SourceRepository) GetSource(ctx context.Context, source *api.Source) (*api.Source, error) {
	result := sr.client.FindFirst(ctx, source)

	if result.Error != nil {
		return nil, result.Error
	}

	return source, nil
}

func (sr *SourceRepository) PutSource(ctx context.Context, sourceInput *api.SourceInput) error {
	uriDigest := helpers.GetSHA256Hash(sourceInput.Uri)
	source := &api.Source{
		Name:      sourceInput.Name,
		Summary:   sourceInput.Summary,
		Tags:      sourceInput.Tags,
		Uri:       sourceInput.Uri,
		UriDigest: uriDigest,
	}

	result := sr.client.Create(ctx, source)
	if result.Error == nil {
		slog.InfoContext(
			ctx,
			fmt.Sprintf("%d rows affected\n", result.RowsAffected),
		)
	}

	return result.Error
}

// Updates source model fields except for `uri` and `uriDigest`
func (sr *SourceRepository) UpdateSourceByUriDigest(ctx context.Context, sourceInput *api.SourceInput, uriDigest string) error {
	source := &api.Source{}
	source.UriDigest = uriDigest

	result := sr.client.FindFirst(ctx, source)
	if result.Error != nil {
		return result.Error
	}

	source.Name = sourceInput.Name
	source.Summary = sourceInput.Summary
	source.Tags = sourceInput.Tags

	result = sr.client.Update(ctx, source)
	if result.Error == nil {
		slog.InfoContext(
			ctx,
			fmt.Sprintf("%d rows affected\n", result.RowsAffected),
		)
	}

	return result.Error
}
