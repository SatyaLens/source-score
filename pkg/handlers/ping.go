package handlers

import (
	"context"
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
	return ph.message
}
