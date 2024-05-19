package main

import (
	"github.com/RinnAnd/ww-backend/server"
)

func main() {
	conn := server.Pool()
	server := server.NewServer(":8080", conn)
	server.RegisterRoutes()
	server.Start()
}
