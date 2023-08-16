package user

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type DBTX interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	PrepareContext(context.Context, string) (*sql.Stmt, error)
	QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...interface{}) *sql.Row
}

type repository struct {
	db DBTX
}

func NewRepository(db DBTX) Repository {
	return &repository{db: db}
}

func (r *repository) CreateUser(ctx context.Context, user *User) (*User, error) {
	// Generate a UUID for user ID
	user.ID = uuid.New()

	// Set the current timestamp for CreatedAt
	user.CreatedAt = time.Now()

	query := "INSERT INTO users(id, username, email, profile_image_url, created_at) VALUES ($1, $2, $3, $4, $5) RETURNING id"

	err := r.db.QueryRowContext(ctx, query, user.ID, user.Username, user.Email, user.ProfileImageURL, user.CreatedAt).Scan(&user.ID)
	if err != nil {
		return &User{}, err
	}

	return user, nil
}

func (r *repository) FindUserByEmail(ctx context.Context, email string) (*User, error) {
	u := User{}
	query := "SELECT email, username FROM users WHERE email = $1"
	err := r.db.QueryRowContext(ctx, query, email).Scan(&u.ID, &u.Email)
	if err != nil {
		return &User{}, nil
	}

	return &u, nil
}
