package main

import (
	"log"
	"source-score/pkg/api"

	"github.com/gin-gonic/gin"
)

//go:generate go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen --config=configs/config.yaml api/source-score.yaml
func main() {
	server := gin.Default()
	err := server.Run()

	if err != nil {
		log.Fatalf("failed to start the server: %s\n", err)
	} else {
		log.Println("Server is up and running!")
	}

	api.RegisterHandlers(server, api.NewRouter())
}
