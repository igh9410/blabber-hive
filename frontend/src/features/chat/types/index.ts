export type MessageType = {
  sender: 'received' | 'sent';
  content: string;
  timestamp: Date;
  img?: string; // Optional image URL for the sender's profile (for 'received' messages)
};

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
