-- +goose Up
-- +goose StatementBegin
-- For the `public.friends` table
ALTER TABLE public.friends
    ALTER COLUMN id SET DATA TYPE integer USING id::integer;

-- For the `public.friend_requests` table
ALTER TABLE public.friend_requests
    ALTER COLUMN id SET DATA TYPE integer USING id::integer;

-- For the `public.messages` table
ALTER TABLE public.messages
    ALTER COLUMN id SET DATA TYPE integer USING id::integer;

-- For the `public.users_in_chat_rooms` table
ALTER TABLE public.users_in_chat_rooms
    ALTER COLUMN id SET DATA TYPE integer USING id::integer;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- For the `public.friends` table
ALTER TABLE public.friends
    ALTER COLUMN id SET DATA TYPE integer USING id::integer;

-- For the `public.friend_requests` table
ALTER TABLE public.friend_requests
    ALTER COLUMN id SET DATA TYPE integer USING id::integer;

-- For the `public.messages` table
ALTER TABLE public.messages
    ALTER COLUMN id SET DATA TYPE integer USING id::integer;

-- For the `public.users_in_chat_rooms` table
ALTER TABLE public.users_in_chat_rooms
    ALTER COLUMN id SET DATA TYPE integer USING id::integer;
-- +goose StatementEnd
