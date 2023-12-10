import { ChatArea, InputArea, MessageType } from '@features/chat';
import { fetchUserFn } from '@features/user';

import { useQuery } from '@tanstack/react-query';
import { useState } from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import styles from './ChatRoom.module.scss';

export function ChatRoom() {
  const [messages, setMessages] = useState<MessageType[]>([]);
  const { data } = useQuery({
    queryKey: ['users'],
    queryFn: fetchUserFn,
  });
  const { id } = useParams(); // id is the chat room ID from the URL
  const chatRoomId = id || '';
  const navigate = useNavigate();

  if (data === null) {
    // If data is null, then the user is first logging in, redirecting to sign up page
    navigate('/signup');
  }
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
