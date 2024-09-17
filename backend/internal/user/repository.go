package user

import (
	"context"

	"log"
	"time"

	"github.com/igh9410/blabber-hive/backend/internal/app/infrastructure/database"
)

type Repository interface {
	CreateUser(ctx context.Context, user *User) (*User, error)
	FindUserByEmail(ctx context.Context, email string) (*User, error)
}

type repository struct {
	db *db.Database
}

func NewRepository(db *db.Database) Repository {
	return &repository{db: db}
}

func (r *repository) CreateUser(ctx context.Context, user *User) (*User, error) {

	tx, err := r.db.Pool.Begin(ctx)
	if err != nil {
		log.Println("Creating User transaction failed")
		return nil, err // handle error appropriately
	}

	// Ensure rollback in case of error
	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(ctx); rbErr != nil {
				log.Printf("Transaction rollback failed: %v", rbErr)
			}
		}
	}()

	// Set the current timestamp for CreatedAt
	user.CreatedAt = time.Now()

	query := "INSERT INTO users(id, username, email, profile_image_url, created_at) VALUES ($1, $2, $3, $4, $5) RETURNING id"

	err = r.db.Pool.QueryRow(ctx, query, user.ID, user.Username, user.Email, user.ProfileImageURL, user.CreatedAt).Scan(&user.ID)

	if err != nil {
		log.Printf("Error creating user, db execcontext: %v", err)
		return nil, err
	}

	// Commit transaction
	if err = tx.Commit(ctx); err != nil {
		log.Printf("Transaction commit failed: %v", err)
		return nil, err
	}

	return user, nil
}

func (r *repository) FindUserByEmail(ctx context.Context, email string) (*User, error) {
	u := User{}
	query := "SELECT id, email, username FROM users WHERE email = $1"
	err := r.db.Pool.QueryRow(ctx, query, email).Scan(&u.ID, &u.Email, &u.Username)
	if err != nil {
		return nil, err
	}

	return &u, nil
}
