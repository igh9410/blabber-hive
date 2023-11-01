package user

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/google/uuid"
)

type Service interface {
	CreateUser(c context.Context, req *CreateUserReq, userID uuid.UUID, email string) (*CreateUserRes, error)
	IsUserRegistered(c context.Context, email string) (bool, error)
	FindUserByEmail(ctx context.Context, email string) (*UserDTO, error)
}

type service struct {
	Repository
	timeout time.Duration
}

func NewService(repository Repository) Service {
	return &service{
		repository,
		time.Duration(2) * time.Second,
	}
}

func (s *service) CreateUser(c context.Context, req *CreateUserReq, userID uuid.UUID, email string) (*CreateUserRes, error) {

	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	u := &User{
		ID:              userID,
		Username:        req.Username,
		Email:           email,
		ProfileImageURL: req.ProfileImageURL,
	}

	r, err := s.Repository.CreateUser(ctx, u)
	if err != nil {
		return nil, err
	}

	res := &CreateUserRes{
		Username:        r.Username,
		ProfileImageURL: r.ProfileImageURL,
	}

	return res, nil

}

func (s *service) IsUserRegistered(c context.Context, email string) (bool, error) {
	_, err := s.Repository.FindUserByEmail(c, email)
	if err == sql.ErrNoRows {
		return false, nil
	} else if err != nil {
		return false, err
	}
	return true, nil
}

func (s *service) FindUserByEmail(c context.Context, email string) (*UserDTO, error) {
	u, err := s.Repository.FindUserByEmail(c, email)
	if err != nil {
		return nil, err
	}
	r := &UserDTO{
		ID:       u.ID,
		Username: u.Username,
	}
	log.Printf("Email: %v", email)
	log.Printf("User ID: %v, Username: %v", r.ID, r.Username)
	return r, nil
}
