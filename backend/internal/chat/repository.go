package chat

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
)

type DBTX interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	PrepareContext(context.Context, string) (*sql.Stmt, error)
	QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...interface{}) *sql.Row
}

type repository struct {
	db DBTX
}

func NewRepository(db DBTX) Repository {
	return &repository{db: db}
}

func (r *repository) CreateChatRoom(ctx context.Context, chatRoom *ChatRoom) (*ChatRoom, error) {
	// Generate a UUID for user ID
	chatRoom.ID = uuid.New()

	// Set the current timestamp for CreatedAt
	chatRoom.CreatedAt = time.Now()

	// Extract the user IDs from the Clients map
	var userIDs []uuid.UUID

	for _, client := range chatRoom.Clients {
		userIDs = append(userIDs, client.User.ID)
	}
	if len(userIDs) != 2 {
		return nil, errors.New("chat room must have exactly two clients")
	}

	query := "INSERT INTO chat_rooms(id, user_id_1, user_id_2, created_at) VALUES ($1, $2, $3, $4) RETURNING id"

	err := r.db.QueryRowContext(ctx, query, chatRoom.ID, userIDs[0], userIDs[1], chatRoom.CreatedAt).Scan(&chatRoom.ID)
	if err != nil {
		return nil, err
	}

	return chatRoom, nil
}
