-- +goose Up
-- +goose StatementBegin
ALTER TABLE messages RENAME TO chat_messages;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE chat_messages RENAME TO messages;
-- +goose StatementEnd
