package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

type Database struct {
	db *sql.DB
}

func NewDatabase() (*Database, error) {

	// Create the connection string using the retrieved values
	// Retrieve the values from environment variables
	password := os.Getenv("SUPABASE_POSTGRES_PASSWORD")
	domain := "db." + os.Getenv("SUPABASE_DOMAIN")

	connectionString := fmt.Sprintf("user=postgres password=%s host=%s port=5432 dbname=postgres", password, domain)
	db, err := sql.Open("postgres", connectionString)
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
