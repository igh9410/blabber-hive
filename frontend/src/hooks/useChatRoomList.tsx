import { ChatRoom, fetchChatRoomListFn } from '@features/chat';
import { useQuery } from '@tanstack/react-query';

export const useChatRoomList = () => {
  const {
    data: chatRooms,
    isLoading,
    error,
  } = useQuery<ChatRoom[], Error>(['chatRooms'], fetchChatRoomListFn);
  return { chatRooms, isLoading, error };
};
