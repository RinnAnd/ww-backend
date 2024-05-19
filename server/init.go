package server

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func Pool() *sql.DB {
	godotenv.Load()
	user := os.Getenv("PG_USER")
	dbname := os.Getenv("DB_NAME")
	password := os.Getenv("PG_PASSWORD")
	sslmode := os.Getenv("SSL")
	connString := fmt.Sprintf("user=%s dbname=%s password=%s sslmode=%s", user, dbname, password, sslmode)

	db, err := sql.Open("postgres", connString)
	if err != nil {
		panic("Could not connect to the database")
	}

	db.Exec(`
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY DEFAULT uuid_generate_v4(),
		username TEXT NOT NULL,
    name TEXT NOT NULL,
    email TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL
)`)

	db.Exec(`
CREATE TABLE IF NOT EXISTS friendships (
    id SERIAL PRIMARY KEY DEFAULT uuid_generate_v4(),
    FOREIGN KEY (user_id1) REFERENCES users (id),
    FOREIGN KEY (user_id2) REFERENCES users (id)
)`)

	db.Exec(`
CREATE TABLE IF NOT EXISTS finances (
    id SERIAL PRIMARY KEY DEFAULT uuid_generate_v4(),
    FOREIGN KEY (user_id) REFERENCES users (id),
		month INTEGER,
		year INTEGER,
		salary INTEGER,
    FOREIGN KEY (saving) REFERENCES savings (id),
)`)

	db.Exec(`
CREATE TABLE IF NOT EXISTS expenses (
    id SERIAL PRIMARY KEY DEFAULT uuid_generate_v4(),
    FOREIGN KEY (user_id) REFERENCES users (id),
    amount DECIMAL,
)`)

	db.Exec(`
CREATE TABLE IF NOT EXISTS savings (
    id SERIAL PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id INTEGER,
    amount DECIMAL,
    FOREIGN KEY (user_id) REFERENCES users (id)
)`)
	return db
}
