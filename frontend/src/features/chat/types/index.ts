export type MessageType = {
  sender: 'received' | 'sent';
  text: string;
  img?: string; // Optional image URL for the sender's profile (for 'received' messages)
};
