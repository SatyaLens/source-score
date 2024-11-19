package main

import (
	"log"
	"source-score/pkg/api"

	"github.com/gin-gonic/gin"
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