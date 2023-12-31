package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/jackc/pgx/v4/stdlib"
	_ "github.com/lib/pq"
)

type Database struct {
	db *sql.DB
}

func NewDatabase() (*Database, error) {

	// Create the connection string using the retrieved values
	// Retrieve the values from environment variables
	var username string
	var password string
	var host string
	var connectionString string

	if os.Getenv("IS_PRODUCTION") == "YES" { // Production Environment
		connectionString = os.Getenv("DATABASE_URL")
	} else {
		username = os.Getenv("POSTGRES_USERNAME")
		password = os.Getenv("POSTGRES_PASSWORD")
		host = os.Getenv("POSTGRES_HOST") // Running in local Docker container

		if host == "" { // Running in local environment
			host = "localhost"
		}

		connectionString = fmt.Sprintf("postgresql://%s:%s@%s:5432/postgres?sslmode=disable", username, password, host)
	}

	db, err := sql.Open("pgx", connectionString)
	if err != nil {
		return nil, err
	}

	// Force a connection to verify it works.
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	log.Println("Database initialized")

	return &Database{db: db}, nil
}

func (d *Database) Close() error {
	return d.db.Close()
}

func (d *Database) GetDB() *sql.DB {
	return d.db
}
