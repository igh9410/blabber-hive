package match

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/igh9410/blabber-hive/backend/internal/app/infrastructure/redis"

	"github.com/google/uuid"
)

type Repository interface {
	EnqueueUser(ctx context.Context, userID uuid.UUID) (*EnqueUserRes, error)
	DequeueUser(ctx context.Context, userID string) error
	FetchCandidates(ctx context.Context) ([]string, error) // Fetches all users in the matchmaking queue
	DequeueUsers(ctx context.Context, userIDs ...string) error
}

type repository struct {
	redisClient *redis.RedisClient
}

func NewRepository(redisClient *redis.RedisClient) Repository {
	return &repository{
		redisClient: redisClient,
	}
}

func (r *repository) EnqueueUser(ctx context.Context, userID uuid.UUID) (*EnqueUserRes, error) {
	// Assuming "match_queue" is the Redis key for the matchmaking queue
	added, err := r.redisClient.SAdd(ctx, "match_queue", userID.String()).Result()

	if err != nil {
		slog.Error("Error adding user to the matchmaking queue:", err.Error(), " in EnqueueUser")
		return nil, err
	}
	if added == 0 { // User already in the queue
		slog.Info("User already in the matchmaking queue")
	}

	return &EnqueUserRes{
		ID:        userID,
		Timestamp: time.Now(),
	}, nil
}

func (r *repository) DequeueUser(ctx context.Context, userID string) error {
	// Logic to remove user from Redis set
	// Use SREM command
	_, err := r.redisClient.SRem(ctx, "match_queue", userID).Result()
	if err != nil {
		slog.Error("Error removing user from the matchmaking queue:", err.Error(), " in DequeueUser")
		return err
	}
	slog.Info(fmt.Sprintf("User with ID %s removed from the matchmaking queue", userID))
	return nil
}

func (r *repository) FetchCandidates(ctx context.Context) ([]string, error) {
	// Fetch user IDs from the match_queue
	// For example, if you're using a Redis SET, you might use SMEMBERS to get all members
	candidates, err := r.redisClient.SMembers(ctx, "match_queue").Result()
	if err != nil {
		slog.Error("Error fetching candidates from the matchmaking queue:", err.Error(), " in FetchCandidates")
		return nil, err
	}

	return candidates, nil
}

func (r *repository) DequeueUsers(ctx context.Context, userIDs ...string) error {
	if len(userIDs) == 0 {
		return nil // No users to remove
	}

	_, err := r.redisClient.SRem(ctx, "match_queue", userIDs).Result()
	if err != nil {
		slog.Error("Error removing users from the matchmaking queue:", err.Error(), " in DequeueUsers")
		return err
	}

	slog.Info(fmt.Sprintf("Users %v removed from the matchmaking queue", userIDs))
	return nil
}
