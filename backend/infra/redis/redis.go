package redis

import (
	"context"
	"log"
	"log/slog"
	"os"

	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	*redis.Client
}

func NewRedisClient() (*RedisClient, error) {
	// Construct the Redis URL from environment variables
	redisURL := os.Getenv("REDIS_URL") //

	// Parse the Redis URL
	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		slog.Error("Failed to parse Redis URL: %s\n", err.Error(), " in NewRedisClient")
	}

	client := redis.NewClient(opt)

	// Create a background context for the Ping operation
	ctx := context.Background()

	// Ping the Redis server to check for a successful connection
	_, err = client.Ping(ctx).Result()
	if err != nil {
		log.Printf("Failed to ping Redis server: %s\n", err)
		return nil, err
	}
	log.Println("Redis initialized")
	return &RedisClient{Client: client}, nil
}
func (r *RedisClient) Close() error {
	return r.Client.Close()
}
