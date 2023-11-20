import { ChatMessagesResponse, fetchChatMessagesFn } from '@features/chat';
import { useInfiniteQuery } from '@tanstack/react-query';

export function useChatMessages(chatRoomId: string) {
  return useInfiniteQuery<ChatMessagesResponse, Error>(
    ['chatMessages', chatRoomId],
    ({ pageParam = '' }) => fetchChatMessagesFn(chatRoomId, pageParam),
    {
      getNextPageParam: (lastPage) => {
        console.log('Last page data: ', lastPage);
        const nextPage = lastPage.nextCursor;
        console.log('Next page cursor: ', nextPage);
        return nextPage;
      },
    }
  );
}
