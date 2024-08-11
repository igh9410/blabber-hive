-- File: internal/pkg/sqlc/queries.sql

-- CreateChatRoom.sql
-- name: CreateChatRoom :exec
INSERT INTO chat_rooms (id, name, created_at)
VALUES ($1, $2, $3);

-- FindChatRoomByID.sql
-- name: FindChatRoomByID :one
SELECT id, name, created_at
FROM chat_rooms
WHERE id = $1;

-- FindChatRoomInfoByID.sql
-- name: FindChatRoomInfoByID :many
SELECT
    uicr.id,
    uicr.user_id,
    uicr.chat_room_id,
    cr.created_at
FROM users_in_chat_rooms AS uicr
INNER JOIN chat_rooms AS cr ON uicr.chat_room_id = cr.id
WHERE cr.id = $1
GROUP BY uicr.id, uicr.user_id, uicr.chat_room_id, cr.created_at;

-- JoinChatRoomByID.sql
-- name: JoinChatRoomByID :exec
INSERT INTO users_in_chat_rooms (user_id, chat_room_id)
VALUES ($1, $2);

-- FindChatRoomList.sql
-- name: FindChatRoomList :many
SELECT id, name, created_at
FROM chat_rooms;

-- GetPaginatedMessages.sql
-- name: GetPaginatedMessages :many
SELECT
    id, chat_room_id, sender_id, content, media_url, created_at, read_at, deleted_by_user_id
FROM messages
WHERE chat_room_id = $1 AND created_at < $2
ORDER BY created_at DESC
LIMIT $3;

-- GetFirstPageMessages.sql
-- name: GetFirstPageMessages :many
SELECT
    id, chat_room_id, sender_id, content, media_url, created_at, read_at, deleted_by_user_id
FROM messages
WHERE chat_room_id = $1
ORDER BY created_at DESC
LIMIT $2;
