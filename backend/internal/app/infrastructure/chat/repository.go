package chat

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/igh9410/blabber-hive/backend/internal/app/domain/chat"
	db "github.com/igh9410/blabber-hive/backend/internal/app/infrastructure/database"
	"github.com/igh9410/blabber-hive/backend/internal/pkg/sqlc"
	"github.com/igh9410/blabber-hive/backend/pkg/pointerutil"
	"github.com/igh9410/blabber-hive/backend/pkg/typeutil"
)

type repository struct {
	db *db.Database
}

func NewRepository(db *db.Database) chat.Repository {
	return &repository{db: db}
}

func (r *repository) CreateChatRoom(ctx context.Context, chatRoom chat.ChatRoom) (*chat.ChatRoom, error) {

	err := r.db.Querier.CreateChatRoom(ctx, sqlc.CreateChatRoomParams{
		ID:        uuid.New(),
		Name:      typeutil.StringToPgtype(chatRoom.Name),
		CreatedAt: typeutil.TimeToPgtype(time.Now()),
	})

	if err != nil {
		slog.Error("Error creating chat room: ", err.Error(), " in CreateChatRoom")
		return nil, err
	}

	return &chatRoom, nil
}

func (r *repository) FindChatRoomByID(ctx context.Context, id uuid.UUID) (*chat.ChatRoom, error) {

	chatRoom, err := r.db.Querier.FindChatRoomByID(ctx, id)

	if err != nil {
		slog.Error("Error finding chat room by ID: ", err.Error(), " in FindChatRoomByID")
		return nil, err
	}

	return &chat.ChatRoom{
		ID:        chatRoom.ID,
		Name:      typeutil.PgtypeToString(chatRoom.Name),
		CreatedAt: typeutil.PgtypeToTime(chatRoom.CreatedAt),
	}, nil

}

// FindChatRoomInfoByID implements Repository.
func (r *repository) FindChatRoomInfoByID(ctx context.Context, chatRoomID uuid.UUID) (*chat.ChatRoomInfo, error) {

	sqlcRows, err := r.db.Querier.FindChatRoomInfoByID(ctx, chatRoomID)
	if err != nil {
		slog.Error("Error finding chat room info by ID", "error", err, "chatRoomID", chatRoomID)
		return nil, fmt.Errorf("failed to find chat room info: %w", err)
	}

	if len(sqlcRows) == 0 {
		return nil, sql.ErrNoRows
	}

	usersInChatRoom := make([]chat.UserInChatRoom, len(sqlcRows))
	for i, row := range sqlcRows {
		usersInChatRoom[i] = chat.UserInChatRoom{
			ID:         row.ID,
			UserID:     typeutil.PgtypeToUUID(row.UserID),
			ChatRoomID: typeutil.PgtypeToUUID(row.ChatRoomID),
		}
	}

	return &chat.ChatRoomInfo{
		ID:        chatRoomID,
		UserList:  usersInChatRoom,
		CreatedAt: typeutil.PgtypeToTime(sqlcRows[0].CreatedAt),
	}, nil

}

func (r *repository) JoinChatRoomByID(ctx context.Context, chatRoomID uuid.UUID, userID uuid.UUID) (*chat.ChatRoom, error) {
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
	chatRoom := &chat.ChatRoom{
		ID: chatRoomID,
	}

	return chatRoom, nil

}

func (r *repository) FindChatRoomList(ctx context.Context) ([]*chat.ChatRoom, error) {
	query := "SELECT id, name, created_at FROM chat_rooms"

	rows, err := r.db.Pool.Query(ctx, query)
	if err != nil {
		slog.Error("Error occured with finding chat room list: ", err.Error(), " in FindChatRoomList")
		return nil, err
	}
	defer rows.Close()

	chatRoomList := make([]*chat.ChatRoom, 0)

	for rows.Next() {
		var name sql.NullString
		chatRoom := &chat.ChatRoom{}

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

func (r *repository) GetPaginatedMessages(ctx context.Context, chatRoomID uuid.UUID, cursor *time.Time, pageSize int) ([]chat.ChatMessage, error) {
	// Use the generated sqlc method
	sqlcRows, err := r.db.Querier.GetPaginatedMessages(ctx, sqlc.GetPaginatedMessagesParams{
		ChatRoomID: typeutil.UUIDToPgtype(chatRoomID),
		CreatedAt:  typeutil.TimeToPgtype(*cursor),
		Limit:      int32(pageSize),
	})
	if err != nil {
		return nil, fmt.Errorf("querying for paginated messages: %w", err)
	}

	// Convert sqlcRows to domain-specific messages
	messages := make([]chat.ChatMessage, len(sqlcRows))
	for i, row := range sqlcRows {
		message := chat.ChatMessage{
			ID:         row.ID,
			ChatRoomID: typeutil.PgtypeToUUID(row.ChatRoomID),
			SenderID:   typeutil.PgtypeToUUID(row.SenderID),
			Content:    typeutil.PgtypeToString(row.Content),
			MediaURL:   typeutil.PgtypeToString(row.MediaUrl),
			CreatedAt:  typeutil.PgtypeToTime(row.CreatedAt),
		}

		// Handle nullable fields
		if row.ReadAt.Valid {
			message.ReadAt = &row.ReadAt.Time
		}

		messages[i] = message
	}

	return messages, nil
}

func (r *repository) GetFirstPageMessages(ctx context.Context, chatRoomID uuid.UUID, pageSize int) ([]chat.ChatMessage, error) {
	// Use the generated sqlc method
	sqlcRows, err := r.db.Querier.GetFirstPageMessages(ctx, sqlc.GetFirstPageMessagesParams{
		ChatRoomID: typeutil.UUIDToPgtype(chatRoomID),
		Limit:      int32(pageSize),
	})
	if err != nil {
		return nil, fmt.Errorf("querying for first page messages: %w", err)
	}

	// Convert sqlcRows to domain-specific messages
	messages := make([]chat.ChatMessage, len(sqlcRows))
	for i, row := range sqlcRows {
		message := chat.ChatMessage{
			ID:         row.ID,
			ChatRoomID: typeutil.PgtypeToUUID(row.ChatRoomID),
			SenderID:   typeutil.PgtypeToUUID(row.SenderID),
			Content:    typeutil.PgtypeToString(row.Content),
			MediaURL:   typeutil.PgtypeToString(row.MediaUrl),
			CreatedAt:  typeutil.PgtypeToTime(row.CreatedAt),
		}

		// Handle nullable fields
		if row.ReadAt.Valid {
			message.ReadAt = pointerutil.TimeToPtr(row.ReadAt.Time)
		}

		if row.DeletedByUserID.Valid {
			uuid := typeutil.PgtypeToUUID(row.DeletedByUserID)
			message.DeletedByUserID = pointerutil.UUIDToPtr(uuid)
		}

		messages[i] = message
	}

	return messages, nil
}
