package handlers

import (
	"context"
	"log/slog"
)

type PingHandler struct {
	message string
}

func NewPingHandler() *PingHandler {
	return &PingHandler{
		message: "Pong",
	}
}

func (ph PingHandler) GetPing(ctx context.Context) string {
	slog.InfoContext(ctx, "returning ping response message")
	return ph.message
}
