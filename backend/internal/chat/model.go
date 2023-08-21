package chat

import (
	"backend/internal/user"
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type Room struct {
	ID      string             `json:"id"`
	Name    string             `json:"name"`
	Clients map[string]*Client `json:"clients"`
}

type ChatRoom struct {
	ID        uuid.UUID          `json:"id"`
	Clients   map[string]*Client `json:"clients"`
	CreatedAt time.Time          `json:"created_at"`
}

type Message struct {
	ID              uuid.UUID `json:"id"`
	ChatRoomID      uuid.UUID `json:"chat_room_id"`
	SenderID        uuid.UUID `json:"sender_id"`
	Content         string    `json:"content"`
	MediaURL        string    `json:"media_url"`
	CreatedAt       time.Time `json:"created_at"`
	ReadAt          time.Time `json:"read_at"`
	DeletedByUserID uuid.UUID `json:"deleted_by_user_id"`
}

type Client struct {
	// The WebSocket connection for this client.
	Conn *websocket.Conn
	// Channel to send messages to the client.
	Send chan *Message
	// Details about the user.
	User *user.User
	// The chat room this client is in (if any).
	ChatRoomID uuid.UUID
}

type CreateChatRoomReq struct {
	Clients map[string]*Client `json:"clients"`
}

type RoomRes struct {
	Clients map[string]*Client `json:"clients"`
}

type CreateChatRoomRes struct {
	Clients map[string]*Client `json:"clients"`
}

type Service interface {
	CreateChatRoom(ctx context.Context, req *CreateChatRoomReq) (*CreateChatRoomRes, error)
}

type Repository interface {
	CreateChatRoom(ctx context.Context, chatRoom *ChatRoom) (*ChatRoom, error)
}
