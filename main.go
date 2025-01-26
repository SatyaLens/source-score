package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"source-score/pkg/api"
	"source-score/pkg/conf"
)

var (
	Logger *zap.Logger
	server *gin.Engine
)

func init() {
	// initialize the config
	conf.LoadConfig()

	// initialize the logger
	Logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("failed to initialize the logger: %s\n", err)
	}
	defer Logger.Sync()

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
