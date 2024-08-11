package service

import "context"

type ChatService interface {
	CreateChatRoom(ctx context.Context) error
}

type chatService struct{}

// CreateChatRoom implements ChatService.
func (c *chatService) CreateChatRoom(ctx context.Context) error {
	panic("unimplemented")
}

func NewChatService() ChatService {
	return &chatService{}
}
