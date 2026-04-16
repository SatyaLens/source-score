package source

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

//go:generate go tool counterfeiter . SourceRepository
type SourceRepository interface {
	DeleteSourceByUriDigest(ctx context.Context, source *api.Source) error
	GetSources(ctx context.Context) ([]api.Source, error)
	GetSourceByUriDigest(ctx context.Context, uriDigest string) (*api.Source, error)
	PostSource(ctx context.Context, sourceInput *api.SourceInput) (string, error)
	PatchSourceByUriDigest(ctx context.Context, sourceInput *api.SourcePatchInput, uriDigest string) error
	UpdateAllScores(ctx context.Context, updatedSources *[]api.Source) error
}

type sourceRepository struct {
	client *pgsql.Client
}

func NewSourceRepository(ctx context.Context, client *pgsql.Client) SourceRepository {
	return &sourceRepository{
		client: client,
	}
}

func (sr *sourceRepository) DeleteSourceByUriDigest(ctx context.Context, source *api.Source) error {
	result := sr.client.Delete(ctx, source)
	slog.InfoContext(
		ctx,
		fmt.Sprintf("%d rows affected\n", result.RowsAffected),
	)

	return result.Error
}

func (sr *sourceRepository) GetSources(ctx context.Context) ([]api.Source, error) {
	var sources []api.Source
	result := sr.client.FindAll(ctx, &sources)

	if result.Error != nil {
		return nil, result.Error
	}

	return sources, nil
}

func (sr *sourceRepository) GetSourceByUriDigest(ctx context.Context, uriDigest string) (*api.Source, error) {
	source := &api.Source{}
	source.UriDigest = uriDigest
	result := sr.client.FindFirst(ctx, source)

	if result.Error != nil {
		return nil, result.Error
	}

	return source, nil
}

func (sr *sourceRepository) PostSource(ctx context.Context, sourceInput *api.SourceInput) (string, error) {
	hash := sha256.New()
	_, err := hash.Write([]byte(sourceInput.Uri))
	if err != nil {
		return "", err
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

	if result.Error != nil {
		return "", result.Error
	}

	return uriDigest, nil
}

// Updates source model fields except for `uri` and `uriDigest`
func (sr *sourceRepository) PatchSourceByUriDigest(ctx context.Context, sourceInput *api.SourcePatchInput, uriDigest string) error {
	source := &api.Source{}
	source.UriDigest = uriDigest

	result := sr.client.FindFirst(ctx, source)
	if result.Error != nil {
		return result.Error
	}

	if sourceInput.Name != nil {
		source.Name = *sourceInput.Name
	}
	if sourceInput.Summary != nil {
		source.Summary = *sourceInput.Summary
	}
	if sourceInput.Tags != nil {
		source.Tags = *sourceInput.Tags
	}

	result = sr.client.Update(ctx, source)
	slog.InfoContext(
		ctx,
		fmt.Sprintf("%d rows affected\n", result.RowsAffected),
	)

	return result.Error
}

func (sr *sourceRepository) UpdateAllScores(ctx context.Context, updatedSources *[]api.Source) error {
	var args []any
	var query strings.Builder
	srcDigests := []string{}
	query.WriteString("UPDATE sources SET score = CASE uri_digest")

	for _, src := range *updatedSources {
		query.WriteString(" WHEN ? THEN CAST(? AS FLOAT)")
		args = append(args, src.UriDigest, fmt.Sprintf("%f", src.Score))
		srcDigests = append(srcDigests, src.UriDigest)
	}

	query.WriteString(" END WHERE uri_digest IN ?")
	args = append(args, srcDigests)

	result := sr.client.DB.Exec(query.String(), args...)
	if result.Error != nil {
		return result.Error
	}

	fmt.Printf("rows updated: %d\n", result.RowsAffected)
	return nil
}
