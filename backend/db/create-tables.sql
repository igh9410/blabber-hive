-- Enable the UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- public.users definition

-- Drop table

-- DROP TABLE users;

CREATE TABLE users (
	id uuid NOT NULL DEFAULT uuid_generate_v4(),
	username varchar(255) NOT NULL,
	email varchar(255) NOT NULL,
	profile_image_url varchar(2048) NULL,
	created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
	CONSTRAINT users_email_key UNIQUE (email),
	CONSTRAINT users_nickname_key UNIQUE (username),
	CONSTRAINT users_pkey PRIMARY KEY (id)
);


-- public.blocks definition

-- Drop table

-- DROP TABLE blocks;

CREATE TABLE blocks (
	user_id uuid NOT NULL,
	blocked_user_id uuid NOT NULL,
	CONSTRAINT blocks_pkey PRIMARY KEY (user_id, blocked_user_id),
	CONSTRAINT blocks_blocked_user_id_fkey FOREIGN KEY (blocked_user_id) REFERENCES users(id),
	CONSTRAINT blocks_user_id_fkey FOREIGN KEY (user_id) REFERENCES users(id)
);


-- public.chat_rooms definition

-- Drop table

-- DROP TABLE chat_rooms;

CREATE TABLE chat_rooms (
	id uuid NOT NULL DEFAULT uuid_generate_v4(),
	user_id_1 uuid NULL,
	user_id_2 uuid NULL,
	created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
	CONSTRAINT chat_rooms_pkey PRIMARY KEY (id),
	CONSTRAINT chat_rooms_user_id_1_user_id_2_key UNIQUE (user_id_1, user_id_2),
	CONSTRAINT chat_rooms_user_id_1_fkey FOREIGN KEY (user_id_1) REFERENCES users(id),
	CONSTRAINT chat_rooms_user_id_2_fkey FOREIGN KEY (user_id_2) REFERENCES users(id)
);


-- public.friends definition

-- Drop table

-- DROP TABLE friends;

CREATE TABLE friends (
	user_id uuid NOT NULL,
	friend_id uuid NOT NULL,
	CONSTRAINT friends_pkey PRIMARY KEY (user_id, friend_id),
	CONSTRAINT friends_friend_id_fkey FOREIGN KEY (friend_id) REFERENCES users(id),
	CONSTRAINT friends_user_id_fkey FOREIGN KEY (user_id) REFERENCES users(id)
);


-- public.messages definition

-- Drop table

-- DROP TABLE messages;

CREATE TABLE messages (
	id uuid NOT NULL DEFAULT uuid_generate_v4(),
	chat_room_id uuid NULL,
	sender_id uuid NULL,
	"content" varchar(1000) NULL,
	media_url varchar(2048) NULL,
	created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
	read_at timestamp NULL,
	deleted_by_user_id uuid NULL,
	CONSTRAINT messages_pkey PRIMARY KEY (id),
	CONSTRAINT messages_chat_room_id_fkey FOREIGN KEY (chat_room_id) REFERENCES chat_rooms(id),
	CONSTRAINT messages_deleted_by_user_id_fkey FOREIGN KEY (deleted_by_user_id) REFERENCES users(id),
	CONSTRAINT messages_sender_id_fkey FOREIGN KEY (sender_id) REFERENCES users(id)
);
CREATE INDEX idx_messages_on_chat_room_id ON public.messages USING btree (chat_room_id);

CREATE TABLE friend_requests (
    id uuid NOT NULL DEFAULT uuid_generate_v4(),
    sender_id uuid NOT NULL,
    recipient_id uuid NOT NULL,
    status varchar(10) NOT NULL DEFAULT 'PENDING',
    created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT friend_requests_pkey PRIMARY KEY (id),
    CONSTRAINT friend_requests_recipient_id_fkey FOREIGN KEY (recipient_id) REFERENCES users(id),
    CONSTRAINT friend_requests_sender_id_fkey FOREIGN KEY (sender_id) REFERENCES users(id)
);