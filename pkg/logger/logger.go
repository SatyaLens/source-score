package logger

import (
	"context"
	"log/slog"

	"github.com/gin-gonic/gin"
)

const (
	slogFields = "slog_fields"
)

type ContextHandler struct {
	slog.Handler
}

// Handle adds contextual attributes to the Record before calling the underlying handler
func (h ContextHandler) Handle(ctx context.Context, r slog.Record) error {
	if attrs, ok := ctx.Value(slogFields).([]slog.Attr); ok {
		for _, v := range attrs {
			r.AddAttrs(v)
		}
	}

	return h.Handler.Handle(ctx, r)
}

// AppendGinCtx adds an slog attribute to the provided context so that it will be included
// in any Record created with such context
func AppendGinCtx(parent *gin.Context, attr slog.Attr) {
	if v, ok := parent.Value(slogFields).([]slog.Attr); ok {
		v = append(v, attr)
		parent.Set(slogFields, v)
	}

	v := []slog.Attr{}
	v = append(v, attr)
	parent.Set(slogFields, v)
}
