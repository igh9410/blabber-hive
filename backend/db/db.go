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
	var username string
	var password string
	var domain string
	var connectionString string

	if os.Getenv("IS_PRODUCTION") == "YES" { // Production Environment
		username = os.Getenv("POSTGRES_USERNAME")
		password = os.Getenv("SUPABASE_POSTGRES_PASSWORD")
		domain = "db." + os.Getenv("SUPABASE_DOMAIN")
		connectionString = fmt.Sprintf("user=%s password=%s host=%s port=5432 dbname=postgres", username, password, domain)
	} else {
		username = os.Getenv("POSTGRES_USERNAME")
		password = os.Getenv("POSTGRES_PASSWORD")
		//domain = "db." + os.Getenv("SUPABASE_DOMAIN")
		connectionString = fmt.Sprintf("postgresql://%s:%s@postgres:5432/blabber_hive?sslmode=disable", username, password)
	}

	db, err := sql.Open("postgres", connectionString)
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
