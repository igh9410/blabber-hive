import styles from './ChatArea.module.scss';
import { Message } from './Message';
import { useChatMessageStore } from '@stores';
import { MessageType, ServerMessageType } from '../types';
import { useChatMessages } from '@hooks';
import { useEffect, useRef, useState } from 'react';
import { authTokenKey } from '@config';

type ChatAreaProps = {
  messages: MessageType[];
  chatRoomId: string;
};

export function ChatArea({
  messages: propMessages,
  chatRoomId,
}: ChatAreaProps) {
  const { messages: storeMessages } = useChatMessageStore((state) => ({
    messages: state.messages,
  }));
  // Call the custom hook and pass the chatRoomId to it

  const {
    data,
    error,
    fetchNextPage,
    hasNextPage,
    isFetchingNextPage,
    status,
  } = useChatMessages(chatRoomId);
  const [initialLoad, setInitialLoad] = useState(true);
  const chatAreaRef = useRef<HTMLDivElement>(null);

  const handleScroll: React.UIEventHandler<HTMLDivElement> = (event) => {
    const { scrollTop } = event.currentTarget;

    if (scrollTop === 0 && hasNextPage && !isFetchingNextPage) {
      fetchNextPage();
    }
  };

  // Scroll to bottom function
  const scrollToBottom = () => {
    if (chatAreaRef.current) {
      chatAreaRef.current.scrollTop = chatAreaRef.current.scrollHeight;
    }
  };

  const supabaseData = localStorage.getItem(authTokenKey);

  let currentUserId = '';
  if (supabaseData) {
    const parsedData = JSON.parse(supabaseData);
    currentUserId = parsedData.user.id;
  } else {
    console.error('No Supabase data found in Local Storage');
  }

  const restAPIMessages: ServerMessageType[] = [];
  const existingMessageIds = new Set(restAPIMessages.map((msg) => msg.id));

  if (data && data.pages) {
    data.pages.forEach((page) => {
      if (page.messages) {
        const pageMessages = page.messages
          .map((message: unknown) => {
            const serverMessage = message as ServerMessageType;
            const adaptedMessage: ServerMessageType = {
              id: serverMessage.id,
              chat_room_id: serverMessage.chat_room_id,
              sender:
                serverMessage.sender_id === currentUserId ? 'sent' : 'received',
              media_url: serverMessage.media_url,
              created_at: serverMessage.created_at,
              content: serverMessage.content,
              sender_id: serverMessage.sender_id,
            };

            return adaptedMessage;
          })
          .filter((message) => {
            // Filter out messages that are already in restAPIMessages
            return !existingMessageIds.has(message.id);
          });

        // Update existingMessageIds to include new messages
        pageMessages.forEach((message) => existingMessageIds.add(message.id));

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

  useEffect(() => {
    if (
      initialLoad &&
      (sortedRestAPIMessages.length > 0 || combinedMessages.length > 0)
    ) {
      scrollToBottom();
      setInitialLoad(false); // Set initial load to false after scrolling
    }
  }, [sortedRestAPIMessages.length, combinedMessages.length, initialLoad]);

  if (status === 'loading') return <p>Loading...</p>;
  if (status === 'error') return <p>Error: {error.message}</p>;

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
