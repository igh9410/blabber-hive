package db

import (
	"database/sql"
	"fmt"
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

	var connectionString string

	if os.Getenv("IS_PRODUCTION") == "YES" { // Production Environment
		connectionString = os.Getenv("DATABASE_URL")
	} else {
		username = os.Getenv("POSTGRES_USERNAME")
		password = os.Getenv("POSTGRES_PASSWORD")
		connectionString = fmt.Sprintf("postgresql://%s:%s@localhost:5432/blabber_hive?sslmode=disable", username, password)
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

	return &Database{db: db}, nil
}

func (d *Database) Close() error {
	return d.db.Close()
}

func (d *Database) GetDB() *sql.DB {
	return d.db
}
