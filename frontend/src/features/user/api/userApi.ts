import { axiosInstance } from '@lib/axios';
import axios from 'axios';
import { User } from '../types';
import { getAccessToken } from '@utils';

const userCheckURL = '/api/users/check';
const signUpURL = '/api/users/register';

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

export async function signUpFn({ username }: { username: string }) {
  try {
    // Retrieve the access token
    const accessToken = await getAccessToken();
    console.log('Access token:', accessToken);

    const userData: { username: string } = {
      username,
    };

    const headers = {
      'Content-Type': 'application/json',
      Authorization: `Bearer ${accessToken}`, // Include the access token in the Authorization header
    };

    const response = await axiosInstance.post(
      signUpURL,
      JSON.stringify(userData),
      { headers }
    );

    console.log('User created');

    return response.data;
  } catch (error) {
    console.error('Error:', error);
    // Handle any errors that occur during the token retrieval or the POST request
  }
}
