package match

import (
	"context"
	"log/slog"
	"time"

	"github.com/google/uuid"
)

type Service interface {
	EnqueUser(ctx context.Context, userID uuid.UUID) (*EnqueUserRes, error)
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

func (s *service) EnqueUser(c context.Context, userID uuid.UUID) (*EnqueUserRes, error) {

	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	r, err := s.Repository.EnqueUser(ctx, userID)

	if err != nil {
		slog.Error("Error enqueuing user:", err)
		return nil, err
	}

	res := &EnqueUserRes{
		ID:        r.ID,
		Timestamp: r.Timestamp,
	}

	return res, nil

}
