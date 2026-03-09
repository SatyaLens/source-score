package handlers

import (
	"context"
	"source-score/pkg/domain/source"
)

type SourceHandler struct {
	sourceSvc source.SourceService
}

func NewSourceHandler(ctx context.Context, sourceSvc source.SourceService) *SourceHandler {
	return &SourceHandler{
		sourceSvc: sourceSvc,
	}
}
