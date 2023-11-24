import { fetchUserFn } from '@features/user';
import { useQuery } from '@tanstack/react-query';
import { getAccessToken } from '@utils';

export const useUsers = () => {
  const { data, isLoading, error } = useQuery({
    queryKey: ['users'],
    queryFn: fetchUserFn,
    staleTime: 5 * 60 * 1000, // 5 minutes
    refetchInterval: 10 * 60 * 1000, // 10 minutes
  });
  return { data, isLoading, error };
};
