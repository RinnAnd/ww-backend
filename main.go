package main

import (
	"github.com/RinnAnd/ww-backend/config"
	"github.com/RinnAnd/ww-backend/database"
	"github.com/RinnAnd/ww-backend/server"
)

func main() {
	cfg := config.Get()
	db, err := database.New(cfg.Database)
	if err != nil {
		panic(err)
	}
	server := server.New(":8080", db)
	server.RegisterRoutes()
	server.Start()
}
