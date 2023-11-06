package chat

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestCreateChatRoom(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewRepository(db)

	chatRoom := &ChatRoom{}
	expectedID := uuid.New()

	// Set up the expectations
	mock.ExpectQuery("INSERT INTO chat_rooms").
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(expectedID))

	// Call the method under test
	result, err := repo.CreateChatRoom(context.Background(), chatRoom)

	// Assert the expectations
	assert.NoError(t, err)
	assert.Equal(t, expectedID, result.ID)

	// Ensure all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
