package main

import (
	"github.com/ambientis-org/hefesto/internal/http/routes"
)

func main() {
	server := routes.GetRouter().Server

	server.Logger.Fatal(server.Start(":8080"))
}
