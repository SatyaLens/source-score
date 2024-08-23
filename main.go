package main

import "fmt"

//go:generate go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen -package api -o api/user.gen.go -generate gin user-app.yaml
//go:generate go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen -package api -o api/server.gen.go -generate gin minimal-api.yaml
func main() {
	fmt.Println("unimplemented")
}