package source

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"log"

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

	if result.Error != nil {
		return result.Error
	}

	log.Printf("%d rows affected\n", result.RowsAffected)

	return nil
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

	if result.Error != nil {
		return result.Error
	}

	log.Printf("%d rows affected\n", result.RowsAffected)

	return nil
}

func (sr *SourceRepository) UpdateSourceByUriDigest(ctx context.Context, uriDigest string) error {
	source := &api.Source{}
	source.UriDigest = uriDigest
	result := sr.client.FindFirst(ctx, source)

	if result.Statement.RaiseErrorOnNotFound {
		log.Printf("no matching record found with uri digest:%s\n", uriDigest)
	} else {
		result = sr.client.Update(ctx, source)
	}

	if result.Error != nil {
		return result.Error
	}

	return nil
}
