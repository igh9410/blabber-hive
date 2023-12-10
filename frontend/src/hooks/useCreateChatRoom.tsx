import { createChatRoomFn } from '@features/chat';
import { useMutation } from '@tanstack/react-query';
import { useNavigate } from 'react-router-dom';

export const useCreateChatRoom = () => {
  const navigate = useNavigate();

  const {
    mutate: createChatRoom,
    isLoading,
    error,
  } = useMutation(createChatRoomFn, {
    onSuccess: (chatRoom) => {
      // Success actions
      if (chatRoom.id) {
        // If the data has an id, it means it was successfully created
        navigate(`/chats/${chatRoom.id}`);
      }
    },
    onError: (error) => {
      // Error actions
      console.error('Error creating chat room', error);
    },
  });

  return {
    createChatRoom, // Changed from addSubscriptions to createChatRoom
    isLoading,
    error,
  };
};
