package chat

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"time"

	"github.com/google/uuid"
)

type Repository interface {
	CreateChatRoom(ctx context.Context, chatRoom *ChatRoom) (*ChatRoom, error)
	FindChatRoomByID(ctx context.Context, chatRoomID uuid.UUID) (*ChatRoom, error)
	FindChatRoomInfoByID(ctx context.Context, chatRoomID uuid.UUID) (*ChatRoomInfo, error)
	JoinChatRoomByID(ctx context.Context, chatRoomID uuid.UUID, userID uuid.UUID) (*ChatRoom, error)
	FindChatRoomList(ctx context.Context) ([]*ChatRoom, error)
	SaveMessage(ctx context.Context, message *Message) error
	GetPaginatedMessages(ctx context.Context, chatRoomID uuid.UUID, cursor *time.Time, pageSize int) ([]Message, error)
	GetFirstPageMessages(ctx context.Context, chatRoomID uuid.UUID, pageSize int) ([]Message, error)
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

	query := "INSERT INTO chat_rooms(id, name, created_at) VALUES ($1, $2, $3) RETURNING id"

	err := r.db.QueryRowContext(ctx, query, chatRoom.ID, chatRoom.Name, chatRoom.CreatedAt).Scan(&chatRoom.ID)
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

func (r *repository) FindChatRoomList(ctx context.Context) ([]*ChatRoom, error) {
	query := "SELECT id, name, created_at FROM chat_rooms"

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		slog.Error("Error occured with finding chat room list: ", err)
		return nil, err
	}
	defer rows.Close()

	chatRoomList := make([]*ChatRoom, 0)

	for rows.Next() {
		var name sql.NullString
		chatRoom := &ChatRoom{}

		err := rows.Scan(&chatRoom.ID, &name, &chatRoom.CreatedAt)
		if err != nil {
			slog.Error("Error occurred while scanning chat room list: ", err)
			return nil, err
		}

		// Convert sql.NullString to string, defaulting to empty string if NULL
		if name.Valid {
			chatRoom.Name = name.String
		} else {
			chatRoom.Name = ""
		}

		chatRoomList = append(chatRoomList, chatRoom)
	}

	return chatRoomList, nil

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

func (r *repository) GetPaginatedMessages(ctx context.Context, chatRoomID uuid.UUID, cursor *time.Time, pageSize int) ([]Message, error) {
	query := `
			SELECT id, chat_room_id, sender_id, content, media_url, created_at, read_at, deleted_by_user_id
			FROM messages
			WHERE chat_room_id = $1 AND created_at < $2
			ORDER BY created_at DESC
			LIMIT $3
		`
	rows, err := r.db.QueryContext(ctx, query, chatRoomID, cursor, pageSize)
	if err != nil {
		return nil, fmt.Errorf("querying for paginated messages: %w", err)
	}
	defer rows.Close()

	var messages []Message
	for rows.Next() {
		var msg Message
		var readAt sql.NullTime
		var deletedByUserID sql.NullString // Use sql.NullString for UUID fields that can be NULL

		if err := rows.Scan(
			&msg.ID,
			&msg.ChatRoomID,
			&msg.SenderID,
			&msg.Content,
			&msg.MediaURL,
			&msg.CreatedAt,
			&readAt,
			&deletedByUserID,
		); err != nil {
			return nil, fmt.Errorf("scanning message: %w", err)
		}

		// Check if readAt is valid, if so, assign it to the struct
		if readAt.Valid {
			msg.ReadAt = &readAt.Time
		}

		// Check if deletedByUserID is valid, if so, convert to uuid.UUID and assign it to the struct
		if deletedByUserID.Valid {
			uid, err := uuid.Parse(deletedByUserID.String)
			if err != nil {
				return nil, fmt.Errorf("parsing UUID: %w", err)
			}
			msg.DeletedByUserID = &uid
		}

		messages = append(messages, msg)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating rows: %w", err)
	}

	return messages, nil
}

func (r *repository) GetFirstPageMessages(ctx context.Context, chatRoomID uuid.UUID, pageSize int) ([]Message, error) {
	query := `
		SELECT id, chat_room_id, sender_id, content, media_url, created_at, read_at, deleted_by_user_id
		FROM messages
		WHERE chat_room_id = $1
		ORDER BY created_at DESC
		LIMIT $2
	`

	rows, err := r.db.QueryContext(ctx, query, chatRoomID, pageSize)
	if err != nil {
		return nil, fmt.Errorf("querying for paginated messages: %w", err)
	}
	defer rows.Close()

	var messages []Message
	for rows.Next() {
		var msg Message
		var readAt sql.NullTime
		var deletedByUserID sql.NullString // Use sql.NullString for UUID fields that can be NULL

		if err := rows.Scan(
			&msg.ID,
			&msg.ChatRoomID,
			&msg.SenderID,
			&msg.Content,
			&msg.MediaURL,
			&msg.CreatedAt,
			&readAt,
			&deletedByUserID,
		); err != nil {
			return nil, fmt.Errorf("scanning message: %w", err)
		}

		// Check if readAt is valid, if so, assign it to the struct
		if readAt.Valid {
			msg.ReadAt = &readAt.Time
		}

		// Check if deletedByUserID is valid, if so, convert to uuid.UUID and assign it to the struct
		if deletedByUserID.Valid {
			uid, err := uuid.Parse(deletedByUserID.String)
			if err != nil {
				return nil, fmt.Errorf("parsing UUID: %w", err)
			}
			msg.DeletedByUserID = &uid
		}

		messages = append(messages, msg)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating rows: %w", err)
	}

	return messages, nil
}
