// tokenUtils.ts or a similar utility file

import { authTokenKey } from '@config';

export const getAccessToken = async (): Promise<string> => {
  return new Promise((resolve, reject) => {
    try {
      const supabaseData = localStorage.getItem(authTokenKey);

      if (!supabaseData) {
        throw new Error('Authentication token not found');
      }

      const parsedData = JSON.parse(supabaseData);
      const accessToken = parsedData.access_token;

      if (!accessToken) {
        throw new Error('Access token is missing in the stored data');
      }

      resolve(accessToken);
    } catch (error) {
      reject(error);
    }
  });
};
