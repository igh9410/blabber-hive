package user

import (
	"context"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestCreateUserRepository(t *testing.T) {
	// Create a new mock database connection
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %s", err)
	}
	defer db.Close()

	r := NewRepository(db)

	mockUser := &User{
		ID:              uuid.New(), // Generate a new UUID
		Username:        "test1",
		Email:           "test1@gmail.com",
		ProfileImageURL: nil, // This can be omitted since it's the zero value for a pointer
		CreatedAt:       time.Now(),
	}

	rows := sqlmock.NewRows([]string{"id"}).AddRow(uuid.New().String())

	mock.ExpectQuery("INSERT INTO users").WithArgs(sqlmock.AnyArg(), mockUser.Username, mockUser.Email, mockUser.ProfileImageURL, sqlmock.AnyArg()).WillReturnRows(rows)

	user, err := r.CreateUser(context.Background(), mockUser)
	assert.NoError(t, err)
	assert.NotNil(t, user)
	// Add more assertions as needed
}
