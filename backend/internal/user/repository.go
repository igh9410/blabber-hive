package user

import (
	"context"
	"database/sql"
	"log"
	"time"
)

type Repository interface {
	CreateUser(ctx context.Context, user *User) (*User, error)
	FindUserByEmail(ctx context.Context, email string) (*User, error)
}

type DBTX interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	PrepareContext(context.Context, string) (*sql.Stmt, error)
	QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...interface{}) *sql.Row
	BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error)
}

type repository struct {
	db DBTX
}

func NewRepository(db DBTX) Repository {
	return &repository{db: db}
}

func (r *repository) CreateUser(ctx context.Context, user *User) (*User, error) {

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		log.Println("Creating User transaction failed")
		return nil, err // handle error appropriately
	}

	// Ensure rollback in case of error
	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				log.Printf("Transaction rollback failed: %v", rbErr)
			}
		}
	}()

	// Set the current timestamp for CreatedAt
	user.CreatedAt = time.Now()

	query := "INSERT INTO users(id, username, email, profile_image_url, created_at) VALUES ($1, $2, $3, $4, $5) RETURNING id"

	err = r.db.QueryRowContext(ctx, query, user.ID, user.Username, user.Email, user.ProfileImageURL, user.CreatedAt).Scan(&user.ID)

	if err != nil {
		log.Printf("Error creating user, db execcontext: %d", err)
		return nil, err
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		log.Printf("Transaction commit failed: %d", err)
		return nil, err
	}

	return user, nil
}

func (r *repository) FindUserByEmail(ctx context.Context, email string) (*User, error) {
	u := User{}
	query := "SELECT id, email, username FROM users WHERE email = $1"
	err := r.db.QueryRowContext(ctx, query, email).Scan(&u.ID, &u.Email, &u.Username)
	if err != nil {
		return nil, err
	}

	return &u, nil
}
