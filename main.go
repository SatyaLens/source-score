package main

import "fmt"

//go:generate go run oapi-codegen -package=api -generate "types,server,spec" user-app.yaml > api/user.gen.go
func main() {
	fmt.Println("unimplemented")
}