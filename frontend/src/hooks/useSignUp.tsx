import { fetchUserFn, signUpFn } from '@features/user';
import { queryClient } from '@lib/react-query';
import { useMutation } from '@tanstack/react-query';
import { useNavigate } from 'react-router-dom';

export const useSignUp = () => {
  const navigate = useNavigate();
  const {
    mutate: signUp,
    isLoading,
    error,
  } = useMutation(signUpFn, {
    onSuccess: async () => {
      try {
        // Fetch the latest user data after successful sign-up
        const userData = await fetchUserFn();

        // Update the 'users' query data with the newly fetched user data
        queryClient.setQueryData(['users'], userData);

        // Navigate to the desired route
        navigate('/');
      } catch (error) {
        console.error('Error during post-sign-up user fetch:', error);
      }
    },
    onError: (error) => {
      // Error actions
      //  console.log('Sign Up Failed');
      console.error('Error: ', error);
    },
  });

  return {
    signUp,
    isLoading,
    error,
  };
};
