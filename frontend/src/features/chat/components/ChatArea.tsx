import React, { useRef, useState } from 'react';
import styles from './ChatArea.module.scss';
import { Message } from './Message';
import { useChatMessageStore } from '@stores';
import { MessageType } from '../types';

type ChatAreaProps = {
  messages: MessageType[];
};

export function ChatArea({ messages: propMessages }: Readonly<ChatAreaProps>) {
  const { messages: storeMessages } = useChatMessageStore((state) => ({
    messages: state.messages,
  }));
  const [loading, setLoading] = useState(false);
  const [cursor, setCursor] = useState<string | null>(null);
  const chatAreaRef = useRef<HTMLDivElement>(null); // Ref for the chat area container */

  return (
    <div className={styles.chatArea}>
      {[...storeMessages, ...propMessages].map((message, timestamp) => (
        <Message key={timestamp} {...message} text={message.content} />
      ))}
      {loading && <div>Loading more messages...</div>}
    </div>
  );
}
