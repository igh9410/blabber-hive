package user

import "context"

type User struct {
	ID              int64  `json:"id"`
	Username        string `json:"username"`
	Email           string `json:"email"`
	ProfileImageUrl string `json:"profile_image_url"`
}

type CreateUserReq struct {
	Username string `json:"username"`
	Email    string `json:"email"`
}

type CreateUserRes struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

type Repository interface {
	CreateUser(ctx context.Context, user *User) (*User, error)
	GetUserByEmail(ctx context.Context, email string) (*User, error)
}

type Service interface {
	CreateUser(c context.Context, req *CreateUserReq) (*CreateUserRes, error)
}
