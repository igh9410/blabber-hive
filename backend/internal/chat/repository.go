package chat

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
	db "github.com/igh9410/blabber-hive/backend/internal/app/infrastructure/database"
	"github.com/igh9410/blabber-hive/backend/internal/pkg/sqlc"
)

type Repository interface {
	CreateChatRoom(ctx context.Context, chatRoom *ChatRoom) (*ChatRoom, error)
	FindChatRoomByID(ctx context.Context, chatRoomID uuid.UUID) (*ChatRoom, error)
	FindChatRoomInfoByID(ctx context.Context, chatRoomID uuid.UUID) (*ChatRoomInfo, error)
	JoinChatRoomByID(ctx context.Context, chatRoomID uuid.UUID, userID uuid.UUID) (*ChatRoom, error)
	FindChatRoomList(ctx context.Context) ([]*ChatRoom, error)

	GetPaginatedMessages(ctx context.Context, chatRoomID uuid.UUID, cursor *time.Time, pageSize int) ([]Message, error)
	GetFirstPageMessages(ctx context.Context, chatRoomID uuid.UUID, pageSize int) ([]Message, error)
}

type repository struct {
	db *db.Database
}

func NewRepository(db *db.Database) Repository {
	return &repository{db: db}
}

func (r *repository) CreateChatRoom(ctx context.Context, chatRoom *ChatRoom) (*ChatRoom, error) {

	err := r.db.Querier.CreateChatRoom(ctx, sqlc.CreateChatRoomParams{
		ID:        uuid.New(),
		Name:      sqlc.StringToPgtype(chatRoom.Name),
		CreatedAt: sqlc.TimeToPgtype(time.Now()),
	})

	if err != nil {
		slog.Error("Error creating chat room: ", err.Error(), " in CreateChatRoom")
		return nil, err
	}

	return chatRoom, nil
}

func (r *repository) FindChatRoomByID(ctx context.Context, id uuid.UUID) (*ChatRoom, error) {

	chatRoom, err := r.db.Querier.FindChatRoomByID(ctx, id)
	if err != nil {
		slog.Error("Error finding chat room by ID: ", err.Error(), " in FindChatRoomByID")
		return nil, err
	}

	return &ChatRoom{
		ID:        chatRoom.ID,
		Name:      sqlc.PgtypeToString(chatRoom.Name),
		CreatedAt: sqlc.PgtypeToTime(chatRoom.CreatedAt),
	}, nil

}

// FindChatRoomInfoByID implements Repository.
func (r *repository) FindChatRoomInfoByID(ctx context.Context, chatRoomID uuid.UUID) (*ChatRoomInfo, error) {

	sqlcRows, err := r.db.Querier.FindChatRoomInfoByID(ctx, chatRoomID)
	if err != nil {
		slog.Error("Error finding chat room info by ID", "error", err, "chatRoomID", chatRoomID)
		return nil, fmt.Errorf("failed to find chat room info: %w", err)
	}

	if len(sqlcRows) == 0 {
		return nil, sql.ErrNoRows
	}

	usersInChatRoom := make([]UserInChatRoom, len(sqlcRows))
	for i, row := range sqlcRows {
		usersInChatRoom[i] = UserInChatRoom{
			ID:         row.ID,
			UserID:     sqlc.PgtypeToUUID(row.UserID),
			ChatRoomID: sqlc.PgtypeToUUID(row.ChatRoomID),
		}
	}

	return &ChatRoomInfo{
		ID:        chatRoomID,
		UserList:  usersInChatRoom,
		CreatedAt: sqlc.PgtypeToTime(sqlcRows[0].CreatedAt),
	}, nil

}

func (r *repository) JoinChatRoomByID(ctx context.Context, chatRoomID uuid.UUID, userID uuid.UUID) (*ChatRoom, error) {
	tx, err := r.db.Pool.Begin(ctx)
	if err != nil {
		slog.Error("Joining chat room transaction failed")
		return nil, err
	}

	// Ensure rollback in case of error
	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(ctx); rbErr != nil {
				slog.Error("Transaction rollback failed: ", rbErr.Error(), " in JoinChatRoomByID")
			}
		}
	}()

	// INSERT into chat rooms, UNIQUE and FOREIGN KEY constraints will handle duplicates and non-existing chat room
	query := `INSERT INTO users_in_chat_rooms (user_id, chat_room_id) VALUES ($1, $2)`
	_, err = tx.Exec(ctx, query, userID, chatRoomID)
	if err != nil {
		slog.Error("Error joining chat room, db execcontext: ", err.Error(), " in JoinChatRoomByID")
		return nil, err
	}

	// Commit transaction
	if err = tx.Commit(ctx); err != nil {
		slog.Error("Transaction commit failed: ", err.Error(), " in JoinChatRoomByID")
		return nil, err
	}
	chatRoom := &ChatRoom{
		ID: chatRoomID,
	}

	return chatRoom, nil

}

func (r *repository) FindChatRoomList(ctx context.Context) ([]*ChatRoom, error) {
	query := "SELECT id, name, created_at FROM chat_rooms"

	rows, err := r.db.Pool.Query(ctx, query)
	if err != nil {
		slog.Error("Error occured with finding chat room list: ", err.Error(), " in FindChatRoomList")
		return nil, err
	}
	defer rows.Close()

	chatRoomList := make([]*ChatRoom, 0)

	for rows.Next() {
		var name sql.NullString
		chatRoom := &ChatRoom{}

		err := rows.Scan(&chatRoom.ID, &name, &chatRoom.CreatedAt)
		if err != nil {
			slog.Error("Error occurred while scanning chat room list: ", err.Error(), " in FindChatRoomList")
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

func (r *repository) GetPaginatedMessages(ctx context.Context, chatRoomID uuid.UUID, cursor *time.Time, pageSize int) ([]Message, error) {
	query := `
			SELECT id, chat_room_id, sender_id, content, media_url, created_at, read_at, deleted_by_user_id
			FROM messages
			WHERE chat_room_id = $1 AND created_at < $2
			ORDER BY created_at DESC
			LIMIT $3
		`
	rows, err := r.db.Pool.Query(ctx, query, chatRoomID, cursor, pageSize)
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

	rows, err := r.db.Pool.Query(ctx, query, chatRoomID, pageSize)
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
