package main

import (
	"fmt"
	
	"github.com/gin-gonic/gin"

	"source-score/pkg/api"
	"source-score/pkg/handlers"
)

//go:generate go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen --config=configs/config.yaml api/user-app.yaml
func main() {
	server := gin.Default()
	api.RegisterHandlers(server, handlers.NewPong())

	// err := server.Run()
	// if err != nil {
	// 	log.Fatalf("failed to start the server: %s\n", err)
	// } else {
	// 	log.Println("Server is up and running!")
	// }

	fmt.Println("unimplemented")
}
