package redis

import (
	"context"
	"log"
	"os"

	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	*redis.Client
}

func NewRedisClient() (*RedisClient, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     os.Getenv(("REDIS_URL")),
		Password: os.Getenv("REDIS_PASSWORD"), // no password set for local
		DB:       0,                           // use default DB
	})

	// Create a background context for the Ping operation
	ctx := context.Background()

	// Ping the Redis server to check for a successful connection
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Printf("Failed to ping Redis server: %s\n", err)
		return nil, err
	}
	log.Println("Redis initialized")
	return &RedisClient{Client: rdb}, nil
}

func (r *RedisClient) Close() error {
	return r.Client.Close()
}
