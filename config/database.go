package config

import "strings"

type DBDriver string

const (
	PostgresDBDriver DBDriver = "postgres"
	MongoDBDriver    DBDriver = "mongodb"
)

type Postgres struct {
	Host     string
	Port     string
	User     string
	Database string
	Password string
	SSLMode  string
}

func (c Postgres) ConnectionString() string {
	components := make([]string, 0)
	if c.Host != "" {
		components = append(components, "host="+c.Host)
	}
	if c.Port != "" {
		components = append(components, "port="+c.Port)
	}
	if c.User != "" {
		components = append(components, "user="+c.User)
	}
	if c.Database != "" {
		components = append(components, "dbname="+c.Database)
	}
	if c.Password != "" {
		components = append(components, "password="+c.Password)
	}
	if c.SSLMode != "" {
		components = append(components, "sslmode="+c.SSLMode)
	}
	return strings.Join(components, " ")
}

type Mongo struct{}

type Database struct {
	Driver   DBDriver
	Postgres Postgres
	Mongo    Mongo
}
