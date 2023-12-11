import { ChatArea, InputArea, MessageType } from '@features/chat';
import { useState } from 'react';
import { useParams } from 'react-router-dom';
import styles from './ChatRoom.module.scss';

export function ChatRoom() {
  const [messages, setMessages] = useState<MessageType[]>([]);

  const { id } = useParams(); // id is the chat room ID from the URL
  const chatRoomId = id || '';

  const handleNewMessage = (text: string) => {
    const newMessage: MessageType = {
      sender: 'sent',
      senderID: '1',
      content: text,
      createdAt: new Date(),
    };
    setMessages((prevMessages) => [...prevMessages, newMessage]);
  };

  return (
    <div className={styles.wrapper}>
      <ChatArea messages={messages} chatRoomId={chatRoomId} />

      <InputArea onMessageSend={handleNewMessage} chatRoomId={chatRoomId} />
    </div>
  );
}
