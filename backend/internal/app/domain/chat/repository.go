package chat

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Repository interface {
	CreateChatRoom(ctx context.Context, chatRoom ChatRoom) (*ChatRoom, error)
	FindChatRoomByID(ctx context.Context, chatRoomID uuid.UUID) (*ChatRoom, error)
	FindChatRoomInfoByID(ctx context.Context, chatRoomID uuid.UUID) (*ChatRoomInfo, error)
	JoinChatRoomByID(ctx context.Context, chatRoomID uuid.UUID, userID uuid.UUID) (*ChatRoom, error)
	FindChatRoomList(ctx context.Context) ([]*ChatRoom, error)

	GetPaginatedMessages(ctx context.Context, chatRoomID uuid.UUID, cursor *time.Time, pageSize int) ([]ChatMessage, error)
	GetFirstPageMessages(ctx context.Context, chatRoomID uuid.UUID, pageSize int) ([]ChatMessage, error)
}
