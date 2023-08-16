package user

import (
	"context"
	"database/sql"
	"time"
)

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

func (s *service) CreateUser(c context.Context, req *CreateUserReq, email string) (*CreateUserRes, error) {
	// Your implementation here
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	u := &User{
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
