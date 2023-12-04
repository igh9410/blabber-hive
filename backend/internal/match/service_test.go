package match

import (
	"backend/internal/chat"
	"context"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type MockRepository struct {
	mock.Mock
}

// Implement each method of the Repository interface
func (m *MockRepository) EnqueueUser(ctx context.Context, userID uuid.UUID) (*EnqueUserRes, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(*EnqueUserRes), args.Error(1)
}

func (m *MockRepository) DequeueUser(ctx context.Context, userID string) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func (m *MockRepository) FetchCandidates(ctx context.Context) ([]string, error) {
	args := m.Called(ctx)
	return args.Get(0).([]string), args.Error(1)
}

type MockChatRepository struct {
	mock.Mock
}

// Assuming CreateChatRoom is a method in your ChatRepository interface
func (m *MockChatRepository) CreateChatRoom(ctx context.Context, chatRoom *chat.ChatRoom) (*chat.ChatRoom, error) {
	args := m.Called(ctx, chatRoom)
	return args.Get(0).(*chat.ChatRoom), args.Error(1)
}
