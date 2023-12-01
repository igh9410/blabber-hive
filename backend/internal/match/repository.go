package match

import (
	"backend/infra/redis"
	"context"
	"time"

	"github.com/google/uuid"
)

type Repository interface {
	EnqueUser(ctx context.Context, userID uuid.UUID) (*EnqueUserRes, error)
}

type repository struct {
	redisClient *redis.RedisClient
}

func NewRepository(redisClient *redis.RedisClient) Repository {
	return &repository{
		redisClient: redisClient,
	}
}

func (r *repository) EnqueUser(ctx context.Context, userID uuid.UUID) (*EnqueUserRes, error) {
	// Assuming "match_queue" is the Redis key for the matchmaking queue
	err := r.redisClient.LPush(ctx, "match_queue", userID.String()).Err()
	if err != nil {
		return nil, err
	}

	return &EnqueUserRes{
		ID:        userID,
		Timestamp: time.Now(),
	}, nil
}
