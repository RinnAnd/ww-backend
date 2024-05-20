package server

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func Pool() *sql.DB {
	fmt.Println("Initializing database connection...")
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

	where, err := TableMaker(db)
	if err != nil {
		fmt.Println(where)
		fmt.Println(err)
		panic("There was an error creating the tables")
	}

	return db
}
