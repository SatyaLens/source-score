package source

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log/slog"

	"source-score/pkg/api"
	"source-score/pkg/db/cnpg"
)

type SourceRepository struct {
	client *cnpg.Client
}

func NewSourceRepository(ctx context.Context, client *cnpg.Client) *SourceRepository {
	return &SourceRepository{
		client: client,
	}
}

func (sr *SourceRepository) DeleteSourceByUriDigest(ctx context.Context, source *api.Source) error {
	result := sr.client.Delete(ctx, source)
	slog.InfoContext(
		ctx,
		fmt.Sprintf("%d rows affected\n", result.RowsAffected),
	)

	return result.Error
}

func (sr *SourceRepository) GetSourceByUriDigest(ctx context.Context, uriDigest string) (*api.Source, error) {
	source := &api.Source{}
	source.UriDigest = uriDigest
	result := sr.client.FindFirst(ctx, source)

	if result.Error != nil {
		return nil, result.Error
	}

	return source, nil
}

func (sr *SourceRepository) PutSource(ctx context.Context, sourceInput *api.SourceInput) error {
	hash := sha256.New()
	_, err := hash.Write([]byte(sourceInput.Uri))
	if err != nil {
		return err
	}

	uriDigest := hex.EncodeToString(hash.Sum(nil))
	source := &api.Source{
		Name:      sourceInput.Name,
		Summary:   sourceInput.Summary,
		Tags:      sourceInput.Tags,
		Uri:       sourceInput.Uri,
		UriDigest: uriDigest,
	}

	result := sr.client.Create(ctx, source)
	slog.InfoContext(
		ctx,
		fmt.Sprintf("%d rows affected\n", result.RowsAffected),
	)

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
	slog.InfoContext(
		ctx,
		fmt.Sprintf("%d rows affected\n", result.RowsAffected),
	)

	return result.Error
}
