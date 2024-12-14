package main

import (
	"log"

	"github.com/gin-gonic/gin"

	"source-score/pkg/api"
	"source-score/pkg/conf"
)

func main() {
	conf.LoadConfig()

	server := gin.Default()
	api.RegisterHandlers(server, api.NewRouter())

	err := server.Run()
	if err != nil {
		log.Fatalf("failed to start the server: %s\n", err)
	} else {
		log.Println("Server is up and running!")
	}
}