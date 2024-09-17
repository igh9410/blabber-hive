package match

import (
	"context"
	"log"
	"log/slog"
	"strings"
	"time"

	"github.com/igh9410/blabber-hive/backend/internal/app/infrastructure/redis"
	"github.com/igh9410/blabber-hive/backend/internal/chat"

	"github.com/google/uuid"
)

type Service interface {
	EnqueueUser(ctx context.Context, userID uuid.UUID) (*EnqueUserRes, error)
	DequeueUser(ctx context.Context, userID string) error
	StartMatchmakingSubscriber(ctx context.Context)
	performMatchmaking(userID string) error
}

type service struct {
	Repository
	ChatRepository chat.Repository // Add this field
	timeout        time.Duration
	redisClient    *redis.RedisClient
}

func NewService(repository Repository, chatRepository chat.Repository, redisClient *redis.RedisClient) Service {
	// Code goes here
	return &service{
		Repository:     repository,
		ChatRepository: chatRepository,
		timeout:        time.Duration(2) * time.Second,
		redisClient:    redisClient,
	}
}

func (s *service) EnqueueUser(c context.Context, userID uuid.UUID) (*EnqueUserRes, error) {

	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	r, err := s.Repository.EnqueueUser(ctx, userID)

	if err != nil {
		slog.Error("Error enqueuing user:", err.Error(), " in EnqueueUser")
		return nil, err
	}

	// Publish an event to the Redis channel
	err = s.redisClient.Publish(ctx, "matchmaking_channel", userID.String()).Err()
	if err != nil {
		slog.Error("Error publishing to matchmaking channel:", err.Error(), " in EnqueueUser")
		return nil, err
	}
	slog.Info("User enqueued successfully")

	res := &EnqueUserRes{
		ID:        r.ID,
		Timestamp: r.Timestamp,
	}

	return res, nil

}

func (s *service) DequeueUser(ctx context.Context, userID string) error {
	// Logic to remove user from Redis set
	// Use SREM command

	return s.Repository.DequeueUser(ctx, userID)
}

func (s *service) StartMatchmakingSubscriber(ctx context.Context) {
	pubsub := s.redisClient.Subscribe(context.Background(), "matchmaking_channel")
	defer pubsub.Close()

	for {
		select {
		case <-ctx.Done():
			slog.Info("Matchmaking subscriber shutting down")
			return
		case msg := <-pubsub.Channel():
			if msg == nil {
				continue // or handle the nil message scenario
			}

			// Trigger matchmaking logic
			if err := s.performMatchmaking(msg.Payload); err != nil {
				slog.Error("Error performing matchmaking for user:", msg.Payload, err)
			}
		}
	}
}

func (s *service) performMatchmaking(userID string) error {
	// Fetch potential match candidates
	candidates, err := s.Repository.FetchCandidates(context.Background())
	if err != nil {
		slog.Error("Error fetching match candidates:", err.Error(), " in performMatchmaking")
		return err
	}

	// Filter out the current user and select a match
	var matchID string
	for _, candidate := range candidates {
		if candidate != userID {
			matchID = candidate
			break
		}
	}

	if matchID == "" {
		slog.Error("EnqueueUser error", "userID", userID, "error", err.Error())
		return nil // or custom error indicating no match found
	}

	// Handle the match
	err = s.handleMatch(userID, matchID)
	if err != nil {
		slog.Error("Error handling match:", err.Error(), " in performMatchmaking")
		return err
	}
	// Rest of the matchmaking logic...
	return nil
}

func (s *service) handleMatch(userID, matchID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), s.timeout)
	defer cancel()

	removedIDs := []string{userID, matchID}                                 // Remove two users from the queue to form a match
	log.Printf("Match found for users: %s", strings.Join(removedIDs, ", ")) // Fix the format specifier from %d to %s
	err := s.Repository.DequeueUsers(ctx, removedIDs...)
	// Code goes here
	if err != nil {
		slog.Error("Error dequeuing users:", err.Error(), " in handleMatch")
		return err
	}
	chatRoom := &chat.ChatRoom{}

	_, err = s.ChatRepository.CreateChatRoom(ctx, chatRoom)
	if err != nil {
		log.Print("Error creating chat room:", err)
		return err
	}
	return nil
}
