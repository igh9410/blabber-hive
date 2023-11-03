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
	FindChatRoomInfoByID(ctx context.Context, chatRoomID uuid.UUID) (*ChatRoomInfo, error)
	FetchRecentMessages(ctx context.Context, chatRoomID uuid.UUID, limit int) ([]Message, error)
	JoinChatRoomByID(ctx context.Context, chatRoomID uuid.UUID, userID uuid.UUID) (*ChatRoom, error)
	SaveMessage(ctx context.Context, message *Message) error
	//FetchMessagesFromRedis(ctx context.Context, chatRoomID uuid.UUID, limit int) ([]Message, error)
	FetchMessagesFromDatabase(ctx context.Context, chatRoomID uuid.UUID, page int, limit int) ([]Message, error)
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

// FindChatRoomInfoByID implements Repository.
func (r *repository) FindChatRoomInfoByID(ctx context.Context, chatRoomID uuid.UUID) (*ChatRoomInfo, error) {
	chatRoomInfo := &ChatRoomInfo{}

	query := "SELECT uicr.user_id, uicr.chat_room_id FROM users_in_chat_rooms AS uicr INNER JOIN chat_rooms AS cr ON uicr.chat_room_id = cr.id WHERE cr.id = $1"

	rows, err := r.db.QueryContext(ctx, query, chatRoomID)
	if err != nil {
		log.Printf("Failed to fetch chat room info, err: %v", err)
		return nil, err
	}

	defer rows.Close()

	// Create a slice to store user IDs associated with the chat room.
	usersInChatRoom := make([]UserInChatRoom, 0)

	for rows.Next() {
		var userInChatRoom UserInChatRoom
		err := rows.Scan(&userInChatRoom.UserID, &userInChatRoom.ChatRoomID)
		if err != nil {
			log.Printf("Failed to scan chat room info, err: %v", err)
			return nil, err
		}
		log.Printf("User ID: %v", userInChatRoom.UserID)
		usersInChatRoom = append(usersInChatRoom, userInChatRoom)
	}

	// Set the chat room ID in the chat room info.
	chatRoomInfo.ID = chatRoomID

	// Set the user IDs in the chat room info.
	chatRoomInfo.UserList = make([]UserInChatRoom, 0, len(usersInChatRoom))
	for _, userInChatRoom := range usersInChatRoom {
		chatRoomInfo.UserList = append(chatRoomInfo.UserList, UserInChatRoom{ID: userInChatRoom.ID, UserID: userInChatRoom.UserID, ChatRoomID: userInChatRoom.ChatRoomID})
	}

	return chatRoomInfo, nil

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

/*
func (r *repository) FetchMessagesFromRedis(ctx context.Context, chatRoomID uuid.UUID, limit int) ([]Message, error) {
	redisKey := "chatroom:" + chatRoomID.String() + ":messages"
	messagesJSON, err := r.redisClient.LRange(ctx, redisKey, -limit, -1).Result()
	if err != nil {
		if err == redis.Nil {
			// Redis key does not exist
			return nil, nil
		}
		// Some other error occurred
		return nil, err
	}

	var messages []Message
	for _, mJSON := range messagesJSON {
		var msg Message
		err := json.Unmarshal([]byte(mJSON), &msg)
		if err != nil {
			log.Printf("Failed to unmarshal message: %v", err)
			// Decide how you want to handle partial failure
			continue
		}
		messages = append(messages, msg)
	}

	return messages, nil
} */

// FetchMessagesFromDatabase retrieves chat messages from PostgreSQL
func (r *repository) FetchMessagesFromDatabase(ctx context.Context, chatRoomID uuid.UUID, page int, limit int) ([]Message, error) {
	offset := page * limit
	query := `
		SELECT id, chat_room_id, sender_id, content, created_at
		FROM messages
		WHERE chat_room_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`
	rows, err := r.db.QueryContext(ctx, query, chatRoomID, limit, offset)
	if err != nil {
		log.Printf("Failed to fetch messages from database: %v", err)
		return nil, err
	}
	defer rows.Close()

	var messages []Message
	for rows.Next() {
		var msg Message
		if err := rows.Scan(&msg.ID, &msg.ChatRoomID, &msg.SenderID, &msg.Content, &msg.CreatedAt); err != nil {
			log.Printf("Failed to scan message: %v", err)
			// Decide how you want to handle partial failure
			continue
		}
		messages = append(messages, msg)
	}

	// Check for any error that occurred during iteration
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return messages, nil
}
