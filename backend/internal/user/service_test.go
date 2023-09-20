package user

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

// MockRepository is a mock implementation of the Repository interface
type MockRepository struct {
	user  *User
	error error
}

func (m *MockRepository) CreateUser(ctx context.Context, user *User) (*User, error) {
	return m.user, m.error
}

func (m *MockRepository) FindUserByEmail(ctx context.Context, email string) (*User, error) {
	return m.user, m.error
}

func TestCreateUserService(t *testing.T) {
	mockRepo := &MockRepository{
		user: &User{
			Username:        testUsername,
			Email:           testEmail,
			ProfileImageURL: nil,
		},
		error: nil,
	}

	s := NewService(mockRepo)
	testID := mockRepo.user.ID
	res, err := s.CreateUser(context.Background(), &CreateUserReq{
		Username:        testUsername,
		ProfileImageURL: nil,
	}, testID, testEmail)

	assert.NoError(t, err)
	assert.Equal(t, testUsername, res.Username)
	assert.Nil(t, res.ProfileImageURL)

	// Test case when the repository returns an error
	mockRepo.error = errors.New("some error")
	_, err = s.CreateUser(context.Background(), &CreateUserReq{
		Username:        testUsername,
		ProfileImageURL: nil,
	}, testID, testEmail)

	assert.Error(t, err)
}

func TestIsUserRegistered(t *testing.T) {
	mockRepo := &MockRepository{
		user: &User{
			Username: testUsername,
			Email:    testEmail,
		},
		error: nil,
	}

	s := NewService(mockRepo)

	isRegistered, err := s.IsUserRegistered(context.Background(), testEmail)

	assert.NoError(t, err)
	assert.True(t, isRegistered)

	// Test case when the user is not registered
	mockRepo.error = sql.ErrNoRows
	isRegistered, err = s.IsUserRegistered(context.Background(), "notfound@example.com")

	assert.NoError(t, err)
	assert.False(t, isRegistered)

	// Test case when the repository returns an error other than sql.ErrNoRows
	mockRepo.error = errors.New("some error")
	_, err = s.IsUserRegistered(context.Background(), testEmail)

	assert.Error(t, err)
}
