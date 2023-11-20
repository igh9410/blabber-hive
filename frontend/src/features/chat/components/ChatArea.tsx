import styles from './ChatArea.module.scss';
import { Message } from './Message';
import { useChatMessageStore } from '@stores';
import { MessageType, ServerMessageType } from '../types';
import { useChatMessages } from '@hooks';
import { useEffect, useRef } from 'react';
import { authTokenKey } from '@config';

type ChatAreaProps = {
  messages: MessageType[];
};

export function ChatArea({ messages: propMessages }: Readonly<ChatAreaProps>) {
  const { messages: storeMessages } = useChatMessageStore((state) => ({
    messages: state.messages,
  }));
  // Call the custom hook and pass the chatRoomId to it
  const chatRoomId = '25e4eb83-5210-448d-be58-3a4c355113be';
  const {
    data,
    error,
    fetchNextPage,
    hasNextPage,
    isFetchingNextPage,
    status,
  } = useChatMessages(chatRoomId);
  const chatAreaRef = useRef<HTMLDivElement>(null);

  const handleScroll: React.UIEventHandler<HTMLDivElement> = (event) => {
    const { scrollTop } = event.currentTarget;

    if (scrollTop === 0 && hasNextPage && !isFetchingNextPage) {
      fetchNextPage();
    }
  };

  if (status === 'loading') return <p>Loading...</p>;
  if (status === 'error') return <p>Error: {error.message}</p>;

  const supabaseData = localStorage.getItem(authTokenKey);

  let currentUserId = '';

  if (supabaseData) {
    const parsedData = JSON.parse(supabaseData);
    currentUserId = parsedData.user.id;
  } else {
    console.error('No Supabase data found in Local Storage');
    return null;
  }

  const restAPIMessages: ServerMessageType[] = [];

  if (data && data.pages) {
    data.pages.forEach((page) => {
      if (page.messages) {
        const pageMessages = page.messages.map((message: unknown) => {
          const serverMessage = message as ServerMessageType;
          // ... rest of your mapping logic
          const adaptedMessage: ServerMessageType = {
            id: serverMessage.id,
            chat_room_id: serverMessage.chat_room_id,
            sender:
              serverMessage.sender_id === currentUserId ? 'sent' : 'received', // Set the sender based on the current user's ID
            media_url: serverMessage.media_url,
            created_at: serverMessage.created_at,
            content: serverMessage.content,
            sender_id: serverMessage.sender_id,
          };

          return adaptedMessage;
        });
        restAPIMessages.push(...pageMessages);
      }
    });
  }

  const sortedRestAPIMessages = restAPIMessages.sort(
    (a, b) =>
      new Date(a.created_at).getTime() - new Date(b.created_at).getTime()
  );

  const combinedMessages = [...storeMessages, ...propMessages];
  combinedMessages.sort(
    (a, b) => new Date(a.createdAt).getTime() - new Date(b.createdAt).getTime()
  );
  return (
    <div className={styles.chatArea} onScroll={handleScroll} ref={chatAreaRef}>
      {sortedRestAPIMessages.map((message) => (
        <Message
          key={message.id}
          {...message}
          sender={message.sender}
          text={message.content}
        />
      ))}
      {combinedMessages.map((message) => (
        <Message {...message} text={message.content} />
      ))}
      <div />
    </div>
  );
}
