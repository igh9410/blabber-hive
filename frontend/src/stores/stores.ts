import { MessageType } from '@features/chat';
import { create } from 'zustand';

interface ChatMessageState {
  messages: MessageType[];
  addMessage: (message: MessageType) => void;
}

export const useChatMessageStore = create<ChatMessageState>((set) => ({
  messages: [],
  addMessage: (message) =>
    set((state) => ({
      messages: [...state.messages, message],
    })),
}));
