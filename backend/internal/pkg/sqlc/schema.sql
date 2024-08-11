-- public.chat_rooms definition

-- Drop table

-- DROP TABLE public.chat_rooms;

CREATE TABLE public.chat_rooms (
	id uuid DEFAULT uuid_generate_v1() NOT NULL,
	created_at timestamp DEFAULT CURRENT_TIMESTAMP NOT NULL,
	"name" varchar(255) NULL,
	CONSTRAINT chat_rooms_pkey PRIMARY KEY (id)
);
CREATE INDEX idx_chat_rooms_name ON public.chat_rooms USING btree (name);


-- public.friends definition

-- Drop table

-- DROP TABLE public.friends;

CREATE TABLE public.friends (
	user_id_1 uuid NULL,
	user_id_2 uuid NULL,
	id int4 DEFAULT nextval('friends_new_id_seq1'::regclass) NOT NULL,
	CONSTRAINT friends_pkey PRIMARY KEY (id)
);
CREATE UNIQUE INDEX unique_user_pair_idx ON public.friends USING btree (LEAST(user_id_1, user_id_2), GREATEST(user_id_1, user_id_2));


-- public.goose_db_version definition

-- Drop table

-- DROP TABLE public.goose_db_version;

CREATE TABLE public.goose_db_version (
	id serial4 NOT NULL,
	version_id int8 NOT NULL,
	is_applied bool NOT NULL,
	tstamp timestamp DEFAULT now() NULL,
	CONSTRAINT goose_db_version_pkey PRIMARY KEY (id)
);


-- public.users definition

-- Drop table

-- DROP TABLE public.users;

CREATE TABLE public.users (
	id uuid DEFAULT uuid_generate_v1() NOT NULL,
	username varchar(255) NOT NULL,
	email varchar(255) NOT NULL,
	profile_image_url varchar(2048) NULL,
	created_at timestamp DEFAULT CURRENT_TIMESTAMP NOT NULL,
	CONSTRAINT users_email_key UNIQUE (email),
	CONSTRAINT users_nickname_key UNIQUE (username),
	CONSTRAINT users_pkey PRIMARY KEY (id)
);


-- public.blocks definition

-- Drop table

-- DROP TABLE public.blocks;

CREATE TABLE public.blocks (
	id uuid DEFAULT uuid_generate_v1() NOT NULL,
	user_id uuid NOT NULL,
	blocked_user_id uuid NOT NULL,
	CONSTRAINT blocks_pkey PRIMARY KEY (user_id, blocked_user_id),
	CONSTRAINT blocks_blocked_user_id_fkey FOREIGN KEY (blocked_user_id) REFERENCES public.users(id),
	CONSTRAINT blocks_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id)
);
CREATE UNIQUE INDEX blocks_unique ON public.blocks USING btree (user_id, blocked_user_id);


-- public.friend_requests definition

-- Drop table

-- DROP TABLE public.friend_requests;

CREATE TABLE public.friend_requests (
	sender_id uuid NOT NULL,
	recipient_id uuid NOT NULL,
	status varchar(10) DEFAULT 'PENDING'::character varying NOT NULL,
	created_at timestamp DEFAULT CURRENT_TIMESTAMP NOT NULL,
	id int4 DEFAULT nextval('friend_requests_new_id_seq1'::regclass) NOT NULL,
	CONSTRAINT friend_requests_pkey PRIMARY KEY (id),
	CONSTRAINT friend_requests_unique UNIQUE (sender_id, recipient_id),
	CONSTRAINT friend_requests_recipient_id_fkey FOREIGN KEY (recipient_id) REFERENCES public.users(id),
	CONSTRAINT friend_requests_sender_id_fkey FOREIGN KEY (sender_id) REFERENCES public.users(id)
);
CREATE INDEX idx_friend_requests_recipient_id ON public.friend_requests USING btree (recipient_id);
CREATE INDEX idx_friend_requests_sender_id ON public.friend_requests USING btree (sender_id);


-- public.messages definition

-- Drop table

-- DROP TABLE public.messages;

CREATE TABLE public.messages (
	chat_room_id uuid NULL,
	sender_id uuid NULL,
	"content" varchar(1000) NULL,
	media_url varchar(2048) NULL,
	created_at timestamp DEFAULT CURRENT_TIMESTAMP NOT NULL,
	deleted_by_user_id uuid NULL,
	read_at timestamp NULL,
	id int4 DEFAULT nextval('messages_new_id_seq1'::regclass) NOT NULL,
	CONSTRAINT messages_pkey PRIMARY KEY (id),
	CONSTRAINT messages_chat_room_id_fkey FOREIGN KEY (chat_room_id) REFERENCES public.chat_rooms(id) ON DELETE CASCADE,
	CONSTRAINT messages_deleted_by_user_id_fkey FOREIGN KEY (deleted_by_user_id) REFERENCES public.users(id)
);
CREATE INDEX idx_messages_chat_room_id ON public.messages USING btree (chat_room_id);
CREATE INDEX idx_messages_deleted_by_user_id ON public.messages USING btree (deleted_by_user_id);
CREATE INDEX idx_messages_on_chat_room_id ON public.messages USING btree (chat_room_id);
CREATE INDEX idx_messages_sender_id ON public.messages USING btree (sender_id);


-- public.users_in_chat_rooms definition

-- Drop table

-- DROP TABLE public.users_in_chat_rooms;

CREATE TABLE public.users_in_chat_rooms (
	id serial4 NOT NULL,
	user_id uuid NULL,
	chat_room_id uuid NULL,
	CONSTRAINT users_in_chat_rooms_pkey PRIMARY KEY (id),
	CONSTRAINT users_in_chat_rooms_unique UNIQUE (user_id, chat_room_id),
	CONSTRAINT users_in_chat_rooms_fk FOREIGN KEY (chat_room_id) REFERENCES public.chat_rooms(id) ON DELETE CASCADE
);