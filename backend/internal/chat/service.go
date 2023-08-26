package chat

import (
	"context"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
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

func (s *service) CreateChatRoom(c context.Context, req *CreateChatRoomReq) (*CreateChatRoomRes, error) {

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

func (s *service) InitializeWebSocketConnection(conn *websocket.Conn, chatRoomID uuid.UUID) error {
	// Initialize connection
	messages, err := s.Repository.FetchRecentMessages(context.Background(), chatRoomID, 50)
	if err != nil {
		return err
	}

	// Send recent messages to client
	for _, msg := range messages {
		if err := conn.WriteJSON(msg); err != nil {
			return err
		}
	}

	return nil
}

func (s *service) HandleChatMessages(conn *websocket.Conn, chatRoomID uuid.UUID) {
	for {
		// Read message from browser
		var msg Message
		err := conn.ReadJSON(&msg)
		if err != nil {
			log.Printf("Error while reading message: %v", err)
			break
		}

		// Save message to database
		msg.ID = uuid.New()
		msg.ChatRoomID = chatRoomID
		if err := s.Repository.SaveMessage(context.Background(), &msg); err != nil {
			log.Printf("Error while saving message: %v", err)
			break
		}

		// Broadcast message to all connected clients (for simplicity, only sending back to the same client)
		if err := conn.WriteJSON(msg); err != nil {
			log.Printf("Error while writing message: %v", err)
			break
		}
	}
}
