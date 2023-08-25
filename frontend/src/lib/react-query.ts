import { DefaultOptions, QueryClient } from '@tanstack/react-query';

const queryConfig: DefaultOptions = {
  queries: {
    refetchOnWindowFocus: false,
    refetchOnReconnect: false,
    retry: 1,
    staleTime: Infinity,
  },
};

export const queryClient = new QueryClient({
  defaultOptions: queryConfig,
});
