package main

import (
	"ww-backend/server"
)

func main() {
	server := server.NewServer(":8080")
	server.RegisterRoutes()
	server.Start()
}
