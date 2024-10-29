package handlers

type PingHandler struct {
	message string
}

func NewPingHandler() *PingHandler {
	return &PingHandler{
		message: "Pong",
	}
}

func (ph PingHandler) GetPing() string {
	return ph.message
}
