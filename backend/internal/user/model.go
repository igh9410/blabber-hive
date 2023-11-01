package user

import (
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

type UserDTO struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"username"`
}

type CreateUserReq struct {
	Username        string  `json:"username"`
	ProfileImageURL *string `json:"profile_image_url,omitempty"`
}

type CreateUserRes struct {
	Username        string  `json:"username"`
	ProfileImageURL *string `json:"profile_image_url,omitempty"`
}
