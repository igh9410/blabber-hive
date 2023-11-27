import { signUpFn } from '@features/user';
import { useMutation } from '@tanstack/react-query';
import { useNavigate } from 'react-router-dom';

export const useSignUp = () => {
  const navigate = useNavigate();
  const {
    mutate: signUp,
    isLoading,
    error,
  } = useMutation(signUpFn, {
    onSuccess: () => {
      // Success actions
      navigate('/');
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
