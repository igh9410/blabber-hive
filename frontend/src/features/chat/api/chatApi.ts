import { axiosInstance } from '@lib/axios';
import { AxiosError } from 'axios';
import { ChatMessagesResponse } from '../types';

// Function to fetch chat messages with proper types
export async function fetchChatMessagesFn(
  chatRoomId: string,
  pageParam?: string // Make pageParam optional as it may not be provided for the initial request
): Promise<ChatMessagesResponse> {
  const chatMessagesURL = `/api/chats/${chatRoomId}/messages?cursor=${
    pageParam ?? ''
  }`;

  try {
    const { data } = await axiosInstance.get<ChatMessagesResponse>(
      chatMessagesURL
    );
    return data;
  } catch (error) {
    const axiosError = error as AxiosError;
    // Handle the error as you see fit for your application context
    // This could involve logging the error, returning a default response, throwing a custom error, etc.
    throw new Error(axiosError.message);
  }
}
