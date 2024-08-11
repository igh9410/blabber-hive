package db

import (
	"context"
	"fmt"
	"log"
	"log/slog"

	"os"

	"github.com/igh9410/blabber-hive/backend/internal/pkg/sqlc"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Database struct {
	Pool    *pgxpool.Pool
	Querier sqlc.Querier
}

func NewDatabase() (*Database, error) {

	// Create the connection string using the retrieved values
	// Retrieve the values from environment variables
	var username string
	var password string
	var host string
	var connectionString string

	if os.Getenv("IS_PRODUCTION") == "YES" || os.Getenv("IS_TEST") == "YES" { // Production Environment or Test Environment
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

	// Connect to the database
	ctx := context.Background()

	pool, err := pgxpool.New(ctx, connectionString)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
		return nil, err
	}

	slog.Info("Database initialized")

	querier := sqlc.New(pool)

	return &Database{Pool: pool, Querier: querier}, nil
}

// Close closes the database connection pool.
func (d *Database) Close() error {
	if d.Pool != nil {
		d.Pool.Close()
	}
	return nil
}

// GetDB returns the underlying *pgxpool.Pool instance for advanced database operations.
func (d *Database) GetDB() *pgxpool.Pool {
	return d.Pool
}
