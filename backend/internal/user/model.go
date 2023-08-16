package user

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID              uuid.UUID `json:"id"`
	Username        string    `json:"username"`
	Email           string    `json:"email"`
	ProfileImageURL *string   `json:"profile_image_url,omitempty"`
	CreatedAt       time.Time `json:"created_at"`
}

type CreateUserReq struct {
	Username        string  `json:"username"`
	ProfileImageURL *string `json:"profile_image_url,omitempty"`
}

type CreateUserRes struct {
	Username        string  `json:"username"`
	ProfileImageURL *string `json:"profile_image_url,omitempty"`
}

type Repository interface {
	CreateUser(ctx context.Context, user *User) (*User, error)
	FindUserByEmail(ctx context.Context, email string) (*User, error)
}

type Service interface {
	CreateUser(c context.Context, req *CreateUserReq, email string) (*CreateUserRes, error)
	IsUserRegistered(c context.Context, email string) (bool, error)
}
