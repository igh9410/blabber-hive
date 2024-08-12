package service

import (
	"context"

	chat "github.com/igh9410/blabber-hive/backend/internal/app/domain/chat"
)

type ChatService interface {
	CreateChatRoom(ctx context.Context, chatRoomName string) (*chat.ChatRoom, error)
}

type chatService struct {
	chatRepo chat.Repository
}

// CreateChatRoom implements ChatService.
func (s *chatService) CreateChatRoom(ctx context.Context, chatRoomName string) (*chat.ChatRoom, error) {

	chatRoom := chat.ChatRoom{
		Name: chatRoomName,
	}

	res, err := s.chatRepo.CreateChatRoom(ctx, chatRoom)
	if err != nil {
		return nil, err
	}

	return res, nil

}

func NewChatService(chatRepo chat.Repository) ChatService {
	return &chatService{
		chatRepo: chatRepo,
	}
}
