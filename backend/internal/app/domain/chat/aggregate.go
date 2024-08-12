package chat

import (
	"time"

	"github.com/google/uuid"
)

type ChatRoomInfo struct {
	ID        uuid.UUID        `json:"id"`
	UserList  []UserInChatRoom `json:"user_list"`
	CreatedAt time.Time        `json:"created_at"`
}
