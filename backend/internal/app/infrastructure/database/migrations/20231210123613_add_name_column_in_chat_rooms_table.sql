-- +goose Up
-- +goose StatementBegin
ALTER TABLE chat_rooms ADD COLUMN name VARCHAR(255);
CREATE INDEX idx_chat_rooms_name ON chat_rooms(name);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_chat_rooms_name;
ALTER TABLE chat_rooms DROP COLUMN name VARCHAR(255);
-- +goose StatementEnd
