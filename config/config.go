package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Database Database
}

func Get() *Config {
	godotenv.Load()

	return &Config{
		Database: Database{
			Driver: DBDriver(os.Getenv("DB_DRIVER")),
			Postgres: Postgres{
				Host:     os.Getenv("DB_HOST"),
				Port:     os.Getenv("PG_PORT"),
				User:     os.Getenv("PG_USER"),
				Database: os.Getenv("PG_DATABASE"),
				Password: os.Getenv("PG_PASSWORD"),
				SSLMode:  os.Getenv("PG_SSLMODE"),
			},
		},
	}
}
