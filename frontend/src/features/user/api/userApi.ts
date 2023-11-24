import { axiosInstance } from '@lib/axios';
import axios from 'axios';
import { User } from '../types';
import { getAccessToken } from '@utils';
import Cookies from 'js-cookie';

const userCheckURL = '/api/users/check';

export async function fetchUserFn(): Promise<User> {
  let accessToken = '';

  try {
    accessToken = getAccessToken();
  } catch (error) {
    if (error instanceof Error) {
      console.error(error.message);
    } else {
      console.error('Failed to retrieve access token');
    }
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
    } else if (error instanceof Error) {
      console.error('Non-Axios error:', error.message);
    } else {
      console.error('An unknown error occurred');
    }
    throw error;
  }
}
