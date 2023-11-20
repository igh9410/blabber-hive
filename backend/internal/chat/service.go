package chat

import (
	"backend/internal/user"
	"context"
	"log"
	"time"

	confluentKafka "github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type Service interface {
	CreateChatRoom(ctx context.Context) (*CreateChatRoomRes, error)
	GetChatRoomByID(ctx context.Context, chatRoomID uuid.UUID) (*ChatRoom, error)
	GetChatRoomInfoByID(ctx context.Context, chatRoomID uuid.UUID) (*ChatRoomInfo, error)
	JoinChatRoomByID(ctx context.Context, chatRoomID uuid.UUID, userID uuid.UUID) (*ChatRoom, error)
	RegisterClient(ctx context.Context, hub *Hub, conn *websocket.Conn, chatroomID uuid.UUID, userID uuid.UUID, kafkaProducer *confluentKafka.Producer) (*Client, error)
	GetPaginatedMessages(ctx context.Context, chatRoomID uuid.UUID, cursorTime *time.Time, pageSize int) ([]Message, error)
}

type service struct {
	Repository
	UserRepository user.Repository
	timeout        time.Duration
	KafkaProducer  *confluentKafka.Producer // Add KafkaProducer as a field
}

func NewService(repository Repository, userRepository user.Repository, kafkaProducer *confluentKafka.Producer) Service {
	return &service{
		repository,
		userRepository,
		time.Duration(2) * time.Second,
		kafkaProducer,
	}
}

func (s *service) CreateChatRoom(c context.Context) (*CreateChatRoomRes, error) {

	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	chatRoom := &ChatRoom{}

	r, err := s.Repository.CreateChatRoom(ctx, chatRoom)
	if err != nil {
		return nil, err
	}

	res := &CreateChatRoomRes{
		ID:        r.ID,
		CreatedAt: r.CreatedAt,
	}

	return res, nil

}

func (s *service) GetChatRoomByID(ctx context.Context, id uuid.UUID) (*ChatRoom, error) {
	return s.Repository.FindChatRoomByID(ctx, id)
}

func (s *service) JoinChatRoomByID(ctx context.Context, chatroomId uuid.UUID, userID uuid.UUID) (*ChatRoom, error) {
	// First, find the chat room to make sure it exists
	chatRoom, err := s.Repository.FindChatRoomByID(ctx, chatroomId)
	if err != nil {
		return nil, err // Handle error appropriately
	}

	// Next, find the user to make sure it exists

	res, err := s.Repository.JoinChatRoomByID(ctx, chatroomId, userID)
	if err != nil {
		log.Printf("Error occured with joining the chatroom with ID %v and user id with %v: %v", chatroomId, userID, err)
		return chatRoom, err
	}
	return res, nil

}

// GetChatRoomInfoByID implements Service.
func (s *service) GetChatRoomInfoByID(ctx context.Context, chatRoomID uuid.UUID) (*ChatRoomInfo, error) {
	chatRoomInfo, err := s.Repository.FindChatRoomInfoByID(ctx, chatRoomID)
	if err != nil {
		log.Printf("Error occured with finding chat room info by ID %v: %v", chatRoomID, err)
		return nil, err
	}
	return chatRoomInfo, nil
}

func (s *service) RegisterClient(ctx context.Context, hub *Hub, conn *websocket.Conn, chatroomID uuid.UUID, userID uuid.UUID, kafkaProducer *confluentKafka.Producer) (*Client, error) {

	senderID := userID

	// Create a new client
	client := NewClient(hub, conn, chatroomID, senderID, kafkaProducer)

	client.hub.register <- client
	return client, nil
}

func (s *service) GetPaginatedMessages(ctx context.Context, chatRoomID uuid.UUID, cursorTime *time.Time, pageSize int) ([]Message, error) {

	if cursorTime == nil { // Fetch first 50 messages if cursor is nil (the user is on the first page)
		res, err := s.Repository.GetFirstPageMessages(ctx, chatRoomID, pageSize)

		if err != nil {
			log.Printf("Error occured with getting first page messages for chat room ID %v: %v", chatRoomID, err)
			return nil, err
		}
		log.Printf("Fetching First page messages for chat room ID %v", chatRoomID)
		return res, nil
	}

	res, err := s.Repository.GetPaginatedMessages(ctx, chatRoomID, cursorTime, pageSize)
	if err != nil {
		log.Printf("Error occured with getting paginated messages for chat room ID %v: %v", chatRoomID, err)
		return nil, err
	}
	log.Printf("Fetching paginated messages for chat room ID %v", chatRoomID)
	return res, nil
}
