package api

import "source-score/pkg/handlers"

type router struct {
	pingHandler   *handlers.PongHandler
	sourceHandler *handlers.SourceHandler
}

func NewRouter() *router {

}