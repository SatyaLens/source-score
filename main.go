package main

import (
	"log"

	"github.com/gin-gonic/gin"

	"source-score/pkg/api"
)

func main() {
	server := gin.Default()
	api.RegisterHandlers(server, api.NewRouter())

	err := server.Run()
	if err != nil {
		log.Fatalf("failed to start the server: %s\n", err)
	} else {
		log.Println("Server is up and running!")
	}
}