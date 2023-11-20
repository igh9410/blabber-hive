import { ChatMessagesResponse, fetchChatMessagesFn } from '@features/chat';
import { useInfiniteQuery } from '@tanstack/react-query';

export function useChatMessages(chatRoomId: string) {
  return useInfiniteQuery<ChatMessagesResponse, Error>(
    ['chatMessages', chatRoomId],
    ({ pageParam = '' }) => fetchChatMessagesFn(chatRoomId, pageParam),
    {
      getNextPageParam: (lastPage, pages) => {
        lastPage.nextCursor ?? '';
      },
    }
  );
}
