import { axiosInstance } from '@lib/axios';
import { AxiosError } from 'axios';
import { ChatMessagesResponse } from '../types';
import { authTokenKey } from '@config';

// Function to fetch chat messages with proper types
export async function fetchChatMessagesFn(
  chatRoomId: string,
  cursor = '' // Make pageParam optional as it may not be provided for the initial request
): Promise<ChatMessagesResponse> {
  // Create the URL to fetch the chat messages from
  const chatMessagesURL = `/api/chats/${chatRoomId}/messages?cursor=${cursor}`;

  const supabaseData = localStorage.getItem(authTokenKey);

  let accessToken = '';

  if (supabaseData) {
    const parsedData = JSON.parse(supabaseData);
    accessToken = parsedData.access_token;
  } else {
    console.error('No Supabase data found in Local Storage');
  }
  try {
    // Make the request to fetch the chat messages
    // console.log('Fetching messages for cursor:', cursor);
    const { data } = await axiosInstance.get<ChatMessagesResponse>(
      chatMessagesURL,
      {
        headers: {
          'Content-Type': 'application/json',
          Authorization: `Bearer ${accessToken}`,
        },
      }
    );

    return data;
  } catch (error) {
    const axiosError = error as AxiosError;
    // Handle the error as you see fit for your application context
    // This could involve logging the error, returning a default response, throwing a custom error, etc.
    throw new Error(axiosError.message);
  }
}
