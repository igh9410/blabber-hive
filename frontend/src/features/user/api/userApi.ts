import { axiosInstance } from '@lib/axios';
import axios from 'axios';
import { User } from '../types';
import { getAccessToken } from '@utils';

const userCheckURL = '/api/users/check';

export async function fetchUserFn(): Promise<User | null> {
  let accessToken = '';

  try {
    accessToken = await getAccessToken();
  } catch (error) {
    console.error('Error fetching user:', error);
    throw error;
  }

  try {
    const response = await axiosInstance.get(userCheckURL, {
      headers: {
        'Content-Type': 'application/json',
        Authorization: `Bearer ${accessToken}`,
      },
    });
    return response.data;
  } catch (error) {
    if (axios.isAxiosError(error)) {
      console.error(`Error (${error.response?.status}):`, error.response?.data);
      if (axios.isAxiosError(error)) {
        if (error.response?.status === 404) {
          // User not found, treat as first login
          return null;
        }
      }
    } else {
      console.error('Error fetching user:', error);
    }
    throw error;
  }
}
