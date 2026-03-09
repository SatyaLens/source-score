package handlers

import (
	"source-score/pkg/domain/source"
)

type SourceHandler struct {
	sourceSvc *source.SourceService
}

func NewSourceHandler(sourceSvc *source.SourceService) *SourceHandler {
	return &SourceHandler{
		sourceSvc: sourceSvc,
	}
}
