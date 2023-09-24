package chat

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"time"

	"github.com/google/uuid"
)

type Repository interface {
	CreateChatRoom(ctx context.Context, chatRoom *ChatRoom) (*ChatRoom, error)
	FindChatRoomByID(ctx context.Context, chatRoomID uuid.UUID) (*ChatRoom, error)
	FetchRecentMessages(ctx context.Context, chatRoomID uuid.UUID, limit int) ([]Message, error)
	JoinChatRoomByID(ctx context.Context, chatRoomID uuid.UUID, userID uuid.UUID) (*ChatRoom, error)

	SaveMessage(ctx context.Context, message *Message) error
}

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

	query := "INSERT INTO chat_rooms(id, created_at) VALUES ($1, $2) RETURNING id"

	err := r.db.QueryRowContext(ctx, query, chatRoom.ID, chatRoom.CreatedAt).Scan(&chatRoom.ID)
	if err != nil {
		log.Printf("Error creating chat room: %v", err)

		return nil, errors.New("failed to create chat room")
	}

	return chatRoom, nil
}

func (r *repository) FindChatRoomByID(ctx context.Context, id uuid.UUID) (*ChatRoom, error) {
	chatRoom := &ChatRoom{}
	query := "SELECT id,created_at FROM chat_rooms WHERE id = $1"
	err := r.db.QueryRowContext(ctx, query, id).Scan(&chatRoom.ID, &chatRoom.CreatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("chat room not found") // or a custom error indicating not found
		}
		return nil, err
	}
	return chatRoom, nil
}

func (r *repository) JoinChatRoomByID(ctx context.Context, chatRoomID uuid.UUID, userID uuid.UUID) (*ChatRoom, error) {
	// Find the existing chat room first
	chatRoom, err := r.FindChatRoomByID(ctx, chatRoomID)
	if err != nil {
		return nil, err
	}

	// INSERT into chat rooms, UNIQUE constraint will prevent duplicates
	query := `INSERT INTO users_in_chat_rooms (user_id, chat_room_id) VALUES ($1, $2)`

	_, err2 := r.db.ExecContext(ctx, query, userID, chatRoomID)

	if err2 != nil {
		return nil, err2
	}
	return chatRoom, nil

}

func (r *repository) FetchRecentMessages(ctx context.Context, chatRoomID uuid.UUID, limit int) ([]Message, error) {
	query := `
        SELECT id, chat_room_id, sender_id, content, media_url, created_at, read_at, deleted_by_user_id
        FROM messages
        WHERE chat_room_id = $1
        ORDER BY created_at DESC
        LIMIT $2
    `

	rows, err := r.db.QueryContext(ctx, query, chatRoomID, limit)
	if err != nil {
		log.Printf("Failed to fetch recent messages for this chat room, err: %v", err)
		return nil, err
	}
	defer rows.Close()

	var messages []Message
	for rows.Next() {
		var msg Message
		if err := rows.Scan(&msg.ID, &msg.ChatRoomID, &msg.SenderID, &msg.Content, &msg.MediaURL, &msg.CreatedAt, &msg.ReadAt, &msg.DeletedByUserID); err != nil {
			return nil, err
		}
		messages = append(messages, msg)
	}

	return messages, nil
}

func (r *repository) SaveMessage(ctx context.Context, message *Message) error {
	query := `
        INSERT INTO messages (id, chat_room_id, sender_id, content, media_url, created_at, read_at, deleted_by_user_id)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
    `

	_, err := r.db.ExecContext(ctx, query, message.ID, message.ChatRoomID, message.SenderID, message.Content, message.MediaURL, message.CreatedAt, message.ReadAt, message.DeletedByUserID)
	if err != nil {
		log.Printf("Problem occured related to saving the message into the db, err: %v", err)
		return err
	}

	return nil
}
