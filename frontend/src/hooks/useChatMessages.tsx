import { useInfiniteQuery } from '@tanstack/react-query';
import { ChatMessagesResponse, fetchChatMessagesFn } from '@features/chat';

export const useChatMessages = (chatRoomId: string) => {
  return useInfiniteQuery<ChatMessagesResponse, Error>(
    ['chatMessages', chatRoomId], // Query key includes chatRoomId to ensure uniqueness per chat room
    ({ pageParam }) => fetchChatMessagesFn(chatRoomId, pageParam), // Fetch function
    {
      getNextPageParam: (lastPage) => lastPage.nextCursor ?? undefined, // Determine the cursor for the next page
    }
  );
};
