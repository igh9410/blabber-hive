import { ChatArea, InputArea, MessageType } from '@features/chat';
import React, { useState } from 'react';

export function Lobby() {
  const [messages, setMessages] = useState<MessageType[]>([]);

  const handleNewMessage = (text: string) => {
    const newMessage: MessageType = {
      sender: 'sent',
      content: text,
      timestamp: new Date(),
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
