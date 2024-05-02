package main

import (
	"github.com/RinnAnd/ww-backend/server"
)

func main() {
	server := server.NewServer(":8080")
	server.Start()
}
