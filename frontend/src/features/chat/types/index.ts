import { UUID } from '@types';

export type ChatRoom = {
  id: UUID;
  name: string;
  createdAt: Date;
};

export type MessageType = {
  sender: 'received' | 'sent';
  senderID: string;
  content: string;
  createdAt: Date;
  img?: string; // Optional image URL for the sender's profile (for 'received' messages)
};

// A type for the incoming message that reflects the server's response
export type ServerMessageType = {
  id: number;
  chat_room_id: string;
  sender_id: string;
  content: string;
  media_url: string;
  created_at: Date;
  sender: 'received' | 'sent';
};

export type CombinedMessageType = ChatMessage & ServerMessageType;

// TypeScript interfaces to represent the API response structure
export interface ChatMessage {
  id: number;
  chatRoomId: string;
  senderId: string;
  content: string;
  mediaUrl: string;
  createdAt: string; // Use string if you plan to manually handle date parsing, otherwise consider using Date type
}

export interface ChatMessagesResponse {
  messages: ChatMessage[];
  nextCursor?: string; // Assuming your API provides a cursor for the next page
}
