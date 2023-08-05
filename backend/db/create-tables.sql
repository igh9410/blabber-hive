-- Enable the UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Create the users table
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    auth0_id VARCHAR(255) NOT NULL,
    username VARCHAR(255) NOT NULL UNIQUE,
    profile_image_url TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Create the friends table
CREATE TABLE friends (
    user_id UUID REFERENCES users(id),
    friend_id UUID REFERENCES users(id),
    PRIMARY KEY (user_id, friend_id)
);

-- Create the blocks table
CREATE TABLE blocks (
    user_id UUID REFERENCES users(id),
    blocked_user_id UUID REFERENCES users(id),
    PRIMARY KEY (user_id, blocked_user_id)
);

-- Create the chat_rooms table
CREATE TABLE chat_rooms (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id_1 UUID REFERENCES users(id),
    user_id_2 UUID REFERENCES users(id),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id_1, user_id_2)
);

-- Create the messages table
CREATE TABLE messages (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    chat_room_id UUID REFERENCES chat_rooms(id),
    sender_id UUID REFERENCES users(id),
    content TEXT,
    media_url TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

