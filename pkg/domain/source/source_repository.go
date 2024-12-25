package source

import (
	"context"
	"log"

	"source-score/pkg/api"
	"source-score/pkg/db/cnpg"
)

type sourceRepository struct {
	client *cnpg.Client
}

func NewSourceRepository(ctx context.Context, client *cnpg.Client) *sourceRepository {
	return &sourceRepository{
		client: client,
	}
}

func (sr *sourceRepository) DeleteSourceByUriDigest(ctx context.Context, source *api.Source) error {
	result := sr.client.Delete(ctx, source)

	if result.Error != nil {
		return result.Error
	}

	log.Printf("%d rows affected\n", result.RowsAffected)

	return nil
}

func (sr *sourceRepository) GetSourceByUriDigest(ctx context.Context, uriDigest string) error {
	var source *api.Source
	result := sr.client.FindByPrimaryKey(ctx, source, uriDigest)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (sr *sourceRepository) PutSource(ctx context.Context, source *api.SourceInput) error {
	result := sr.client.Create(ctx, source)

	if result.Error != nil {
		return result.Error
	}

	log.Printf("%d rows affected\n", result.RowsAffected)

	return nil
}

func (sr *sourceRepository) UpdateSourceByUriDigest(ctx context.Context, uriDigest string) error {
	var source *api.Source
	result := sr.client.FindByPrimaryKey(ctx, source, uriDigest)

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
