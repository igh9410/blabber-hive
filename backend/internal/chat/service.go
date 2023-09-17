package chat

import (
	"backend/internal/user"
	"context"
	"errors"
	"log"
	"time"

	confluentKafka "github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type Service interface {
	CreateChatRoom(ctx context.Context) (*CreateChatRoomRes, error)
	GetChatRoomByID(ctx context.Context, chatRoomID uuid.UUID) (*ChatRoom, error)
	JoinChatRoomByID(ctx context.Context, chatRoomID uuid.UUID, userID uuid.UUID) (*ChatRoom, error)
	RegisterClient(ctx context.Context, hub *Hub, conn *websocket.Conn, chatroomID uuid.UUID, senderEmail string, kafkaProducer *confluentKafka.Producer) (*Client, error)
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

	// Check if the chat room is full
	if chatRoom.UserID1 != uuid.Nil && chatRoom.UserID2 != uuid.Nil {
		return nil, errors.New("chat room is full") // Handle this as you see fit
	}

	res, err := s.Repository.JoinChatRoomByID(ctx, chatroomId, userID)
	if err != nil {
		log.Printf("Error occured with joining the chatroom with ID %v and user id with %v: %v", chatroomId, userID, err)
		return chatRoom, err
	}
	return res, nil

}

func (s *service) RegisterClient(ctx context.Context, hub *Hub, conn *websocket.Conn, chatroomID uuid.UUID, senderEmail string, kafkaProducer *confluentKafka.Producer) (*Client, error) {

	sender, err := s.UserRepository.FindUserByEmail(ctx, senderEmail)
	if err != nil {
		log.Printf("Error occured with finding the user with email %v", senderEmail)
		return nil, err
	}

	senderID := sender.ID

	client := NewClient(hub, conn, chatroomID, senderID, kafkaProducer)

	client.hub.register <- client
	return client, nil
}
