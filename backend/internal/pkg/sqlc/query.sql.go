// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: query.sql

package sqlc

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

const createChatRoom = `-- name: CreateChatRoom :exec

INSERT INTO chat_rooms (id, name, created_at)
VALUES ($1, $2, $3)
`

type CreateChatRoomParams struct {
	ID        uuid.UUID        `json:"id"`
	Name      pgtype.Text      `json:"name"`
	CreatedAt pgtype.Timestamp `json:"created_at"`
}

// File: internal/pkg/sqlc/queries.sql
// CreateChatRoom.sql
func (q *Queries) CreateChatRoom(ctx context.Context, arg CreateChatRoomParams) error {
	_, err := q.db.Exec(ctx, createChatRoom, arg.ID, arg.Name, arg.CreatedAt)
	return err
}

const findChatRoomByID = `-- name: FindChatRoomByID :one
SELECT id, name, created_at
FROM chat_rooms
WHERE id = $1
`

type FindChatRoomByIDRow struct {
	ID        uuid.UUID        `json:"id"`
	Name      pgtype.Text      `json:"name"`
	CreatedAt pgtype.Timestamp `json:"created_at"`
}

// FindChatRoomByID.sql
func (q *Queries) FindChatRoomByID(ctx context.Context, id uuid.UUID) (FindChatRoomByIDRow, error) {
	row := q.db.QueryRow(ctx, findChatRoomByID, id)
	var i FindChatRoomByIDRow
	err := row.Scan(&i.ID, &i.Name, &i.CreatedAt)
	return i, err
}

const findChatRoomInfoByID = `-- name: FindChatRoomInfoByID :many
SELECT
    uicr.id,
    uicr.user_id,
    uicr.chat_room_id,
    cr.created_at
FROM users_in_chat_rooms AS uicr
INNER JOIN chat_rooms AS cr ON uicr.chat_room_id = cr.id
WHERE cr.id = $1
GROUP BY uicr.id, uicr.user_id, uicr.chat_room_id, cr.created_at
`

type FindChatRoomInfoByIDRow struct {
	ID         int32            `json:"id"`
	UserID     pgtype.UUID      `json:"user_id"`
	ChatRoomID pgtype.UUID      `json:"chat_room_id"`
	CreatedAt  pgtype.Timestamp `json:"created_at"`
}

// FindChatRoomInfoByID.sql
func (q *Queries) FindChatRoomInfoByID(ctx context.Context, id uuid.UUID) ([]FindChatRoomInfoByIDRow, error) {
	rows, err := q.db.Query(ctx, findChatRoomInfoByID, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []FindChatRoomInfoByIDRow{}
	for rows.Next() {
		var i FindChatRoomInfoByIDRow
		if err := rows.Scan(
			&i.ID,
			&i.UserID,
			&i.ChatRoomID,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const findChatRoomList = `-- name: FindChatRoomList :many
SELECT id, name, created_at
FROM chat_rooms
`

type FindChatRoomListRow struct {
	ID        uuid.UUID        `json:"id"`
	Name      pgtype.Text      `json:"name"`
	CreatedAt pgtype.Timestamp `json:"created_at"`
}

// FindChatRoomList.sql
func (q *Queries) FindChatRoomList(ctx context.Context) ([]FindChatRoomListRow, error) {
	rows, err := q.db.Query(ctx, findChatRoomList)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []FindChatRoomListRow{}
	for rows.Next() {
		var i FindChatRoomListRow
		if err := rows.Scan(&i.ID, &i.Name, &i.CreatedAt); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getFirstPageMessages = `-- name: GetFirstPageMessages :many
SELECT
    id, chat_room_id, sender_id, content, media_url, created_at, read_at, deleted_by_user_id
FROM messages
WHERE chat_room_id = $1
ORDER BY created_at DESC
LIMIT $2
`

type GetFirstPageMessagesParams struct {
	ChatRoomID pgtype.UUID `json:"chat_room_id"`
	Limit      int32       `json:"limit"`
}

type GetFirstPageMessagesRow struct {
	ID              int32            `json:"id"`
	ChatRoomID      pgtype.UUID      `json:"chat_room_id"`
	SenderID        pgtype.UUID      `json:"sender_id"`
	Content         pgtype.Text      `json:"content"`
	MediaUrl        pgtype.Text      `json:"media_url"`
	CreatedAt       pgtype.Timestamp `json:"created_at"`
	ReadAt          pgtype.Timestamp `json:"read_at"`
	DeletedByUserID pgtype.UUID      `json:"deleted_by_user_id"`
}

// GetFirstPageMessages.sql
func (q *Queries) GetFirstPageMessages(ctx context.Context, arg GetFirstPageMessagesParams) ([]GetFirstPageMessagesRow, error) {
	rows, err := q.db.Query(ctx, getFirstPageMessages, arg.ChatRoomID, arg.Limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetFirstPageMessagesRow{}
	for rows.Next() {
		var i GetFirstPageMessagesRow
		if err := rows.Scan(
			&i.ID,
			&i.ChatRoomID,
			&i.SenderID,
			&i.Content,
			&i.MediaUrl,
			&i.CreatedAt,
			&i.ReadAt,
			&i.DeletedByUserID,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getPaginatedMessages = `-- name: GetPaginatedMessages :many
SELECT
    id, chat_room_id, sender_id, content, media_url, created_at, read_at, deleted_by_user_id
FROM messages
WHERE chat_room_id = $1 AND created_at < $2
ORDER BY created_at DESC
LIMIT $3
`

type GetPaginatedMessagesParams struct {
	ChatRoomID pgtype.UUID      `json:"chat_room_id"`
	CreatedAt  pgtype.Timestamp `json:"created_at"`
	Limit      int32            `json:"limit"`
}

type GetPaginatedMessagesRow struct {
	ID              int32            `json:"id"`
	ChatRoomID      pgtype.UUID      `json:"chat_room_id"`
	SenderID        pgtype.UUID      `json:"sender_id"`
	Content         pgtype.Text      `json:"content"`
	MediaUrl        pgtype.Text      `json:"media_url"`
	CreatedAt       pgtype.Timestamp `json:"created_at"`
	ReadAt          pgtype.Timestamp `json:"read_at"`
	DeletedByUserID pgtype.UUID      `json:"deleted_by_user_id"`
}

// GetPaginatedMessages.sql
func (q *Queries) GetPaginatedMessages(ctx context.Context, arg GetPaginatedMessagesParams) ([]GetPaginatedMessagesRow, error) {
	rows, err := q.db.Query(ctx, getPaginatedMessages, arg.ChatRoomID, arg.CreatedAt, arg.Limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetPaginatedMessagesRow{}
	for rows.Next() {
		var i GetPaginatedMessagesRow
		if err := rows.Scan(
			&i.ID,
			&i.ChatRoomID,
			&i.SenderID,
			&i.Content,
			&i.MediaUrl,
			&i.CreatedAt,
			&i.ReadAt,
			&i.DeletedByUserID,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const joinChatRoomByID = `-- name: JoinChatRoomByID :exec
INSERT INTO users_in_chat_rooms (user_id, chat_room_id)
VALUES ($1, $2)
`

type JoinChatRoomByIDParams struct {
	UserID     pgtype.UUID `json:"user_id"`
	ChatRoomID pgtype.UUID `json:"chat_room_id"`
}

// JoinChatRoomByID.sql
func (q *Queries) JoinChatRoomByID(ctx context.Context, arg JoinChatRoomByIDParams) error {
	_, err := q.db.Exec(ctx, joinChatRoomByID, arg.UserID, arg.ChatRoomID)
	return err
}