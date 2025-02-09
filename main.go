package main

import (
	"log"
	"log/slog"
	"os"

	"github.com/gin-gonic/gin"

	"source-score/pkg/api"
	"source-score/pkg/conf"
	"source-score/pkg/helpers"
	"source-score/pkg/logger"
)

var (
	server *gin.Engine
)

func init() {
	// initialize the config
	conf.LoadConfig()

	// initialize the logger
	slog.SetDefault(
		slog.New(&logger.ContextHandler{
			Handler: slog.NewJSONHandler(os.Stdout, nil),
		}),
	)

	// initialize the server
	loggerOpts := api.GinServerOptions{
		Middlewares: []api.MiddlewareFunc{
			// function to add request headers to log fields
			func(c *gin.Context) {
				for _, fieldKey := range helpers.ApiReqLogFields {
					fieldVal := c.Request.Header.Get(fieldKey)
					logger.AppendGinCtx(c, slog.String(fieldKey, fieldVal))
				}
			},
		},
	}

	server := gin.Default()
	api.RegisterHandlersWithOptions(server, api.NewRouter(), loggerOpts)
}

func main() {
	err := server.Run()
	if err != nil {
		log.Fatalf("failed to start the server : %s\n", err.Error())
	} else {
		log.Println("Server is up and running!")
	}
}
