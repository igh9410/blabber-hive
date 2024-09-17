package chat

import (
	"time"

	"github.com/google/uuid"
)

type ChatRoom struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

type Message struct {
	ID              int32      `json:"id"`
	ChatRoomID      uuid.UUID  `json:"chat_room_id"`
	SenderID        uuid.UUID  `json:"sender_id"`
	Content         string     `json:"content"`
	MediaURL        string     `json:"media_url"`
	CreatedAt       time.Time  `json:"created_at"`
	ReadAt          *time.Time `json:"read_at,omitempty"`
	DeletedByUserID *uuid.UUID `json:"deleted_by_user_id,omitempty"`
}

type UserInChatRoom struct {
	ID         int32     `json:"id"`
	UserID     uuid.UUID `json:"user_id"`
	ChatRoomID uuid.UUID `json:"chat_room_id"`
}

type ChatRoomInfo struct {
	ID        uuid.UUID        `json:"id"`
	UserList  []UserInChatRoom `json:"user_list"`
	CreatedAt time.Time        `json:"created_at"`
}

type CreateChatRoomReq struct {
	Name string `json:"name"`
}

type CreateChatRoomRes struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}
