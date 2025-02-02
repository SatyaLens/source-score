package main

import (
	"log/slog"
	"os"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"source-score/pkg/api"
	"source-score/pkg/conf"
	"source-score/pkg/logger"
)

var (
	Logger *zap.Logger
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
	server := gin.Default()
	api.RegisterHandlers(server, api.NewRouter())
}

func main() {
	err := server.Run()
	if err != nil {
		Logger.Fatal(
			"failed to start the server",
			zap.String("err", err.Error()),
		)
	} else {
		Logger.Info("Server is up and running!")
	}
}
