import React from 'react';
import styles from './ChatArea.module.scss';
import { Message } from './Message';

type MessageType = {
  sender: 'received' | 'sent';
  text: string;
  img?: string; // Optional image URL for the sender's profile (for 'received' messages)
};

type ChatAreaProps = {
  messages: MessageType[];
};

export function ChatArea({ messages }: ChatAreaProps) {
  // Add the messages prop})
  return (
    <div className={styles.chatArea}>
      {messages.map((message, index) => (
        <Message key={index} {...message} />
      ))}
    </div>
  );
}
