import React from 'react';
import styles from './ChatArea.module.scss';
import { Message } from './Message';
import { useChatMessageStore } from '@stores';

export type MessageType = {
  sender: 'received' | 'sent';
  content: string;
  timestamp: Date;
  img?: string; // Optional image URL for the sender's profile (for 'received' messages)
};

type ChatAreaProps = {
  messages: MessageType[];
};

export function ChatArea({ messages: propMessages }: ChatAreaProps) {
  const { messages: storeMessages } = useChatMessageStore((state) => ({
    messages: state.messages,
  }));

  return (
    <div className={styles.chatArea}>
      {[...storeMessages, ...propMessages].map((message, timestamp) => (
        <Message key={timestamp} {...message} text={message.content} />
      ))}
    </div>
  );
}
