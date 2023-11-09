import styles from './ChatArea.module.scss';
import { Message } from './Message';
import { useChatMessageStore } from '@stores';
import { MessageType, ServerMessageType } from '../types';
import { useChatMessages } from '@hooks';
import { useEffect } from 'react';
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
  useEffect(() => {
    console.log('Data loaded', data?.pages[0].messages);
  }, [data]);

  if (status === 'loading') return <p>Loading...</p>;
  if (status === 'error') return <p>Error: {error.message}</p>;

  const supabaseData = localStorage.getItem(authTokenKey);

  let currentUserId = '';

  if (supabaseData) {
    const parsedData = JSON.parse(supabaseData);
    currentUserId = parsedData.user.id;
    console.log('Current user id', currentUserId);
  } else {
    console.error('No Supabase data found in Local Storage');
  }

  const restAPIMessages: ServerMessageType[] =
    data?.pages.flatMap((page) =>
      page.messages.map((message: unknown) => {
        // Assert the message to be of type ServerMessageType
        const serverMessage = message as unknown as ServerMessageType;

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
      })
    ) ?? [];

  const combinedMessages = [...storeMessages, ...propMessages].sort(
    (a, b) => new Date(a.createdAt).getTime() - new Date(b.createdAt).getTime()
  );
  return (
    <div className={styles.chatArea}>
      {restAPIMessages.map((message) => (
        <Message {...message} sender={message.sender} text={message.content} />
      ))}
      {combinedMessages.map((message) => (
        <Message {...message} text={message.content} />
      ))}
    </div>
  );
}
