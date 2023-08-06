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
	// Load .env file from the /backend directory
	/*	err := godotenv.Load("../.env")
		if err != nil {
			log.Fatal("Error loading .env file")
		} */

	/* Retrieve the username and password from environment variables

	username := os.Getenv("POSTGRES_USERNAME")
	password := os.Getenv("POSTGRES_PASSWORD")

	// Create the connection string using the retrieved values
	// local connectionString := fmt.Sprintf("postgresql://%s:%s@localhost:5432/blabber_hive?sslmode=disable", username, password) */

	// Retrieve the values from environment variables
	password := os.Getenv("SUPABASE_POSTGRES_PASSWORD")
	domain := "db." + os.Getenv("SUPABASE_DOMAIN")

	connectionString := fmt.Sprintf("postgres:%s@%s/postgres", password, domain)

	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, err
	}

	return &Database{db: db}, nil
}

func (d *Database) Close() {
	d.db.Close()
}

func (d *Database) GetDB() *sql.DB {
	return d.db
}
