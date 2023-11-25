import { ChatArea, InputArea, MessageType } from '@features/chat';
import { fetchUserFn } from '@features/user';

import { useQuery } from '@tanstack/react-query';
import { useState } from 'react';
import { useNavigate } from 'react-router-dom';

export function Lobby() {
  const [messages, setMessages] = useState<MessageType[]>([]);
  const { data } = useQuery({
    queryKey: ['users'],
    queryFn: fetchUserFn,
  });
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
    <>
      <ChatArea messages={messages} />
      <InputArea onMessageSend={handleNewMessage} />
    </>
  );
}
