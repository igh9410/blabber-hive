import { axiosInstance } from '@lib/axios';
import { AxiosError } from 'axios';
import { ChatMessagesResponse } from '../types';
import { authTokenKey } from '@config';

// Function to fetch chat messages with proper types
export async function fetchChatMessagesFn(
  chatRoomId: string,
  pageParam?: string // Make pageParam optional as it may not be provided for the initial request
): Promise<ChatMessagesResponse> {
  // Parse the string into a Date object
  const dateTimeString = '2023-11-09 21:06:13.567';
  const date = new Date(dateTimeString);

  // Format the Date object as an ISO 8601 string (commonly used in APIs)
  const isoDateString = date.toISOString();
  const chatMessagesURL = `/api/chats/${chatRoomId}/messages?cursor=${isoDateString}`;

  console.log(isoDateString); // Outputs: 2023-11-09T21:06:13.567Z (note the 'Z' indicating UTC)

  const supabaseData = localStorage.getItem(authTokenKey);

  let accessToken = '';

  if (supabaseData) {
    const parsedData = JSON.parse(supabaseData);
    accessToken = parsedData.access_token;
  } else {
    console.error('No Supabase data found in Local Storage');
  }
  try {
    console.log('Token = ', accessToken);
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
