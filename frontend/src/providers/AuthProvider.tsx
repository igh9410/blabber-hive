import React, {
  createContext,
  useContext,
  useState,
  useEffect,
  ReactNode,
} from 'react';
import { createClient, SupabaseClient, Session } from '@supabase/supabase-js';
import { SUPABASE_API_KEY, SUPABASE_URL } from '@config';

export const supabase: SupabaseClient = createClient(
  SUPABASE_URL,
  SUPABASE_API_KEY
);

interface AuthContextType {
  session: Session | null;
  isFirstLogin: boolean;
  setIsFirstLogin: React.Dispatch<React.SetStateAction<boolean>>;
}

const AuthContext = createContext<AuthContextType>(null!); // Using `null!` as a workaround for initial value

export const useAuth = () => useContext(AuthContext);

interface AuthProviderProps {
  children: ReactNode;
}

export const AuthProvider = ({ children }: AuthProviderProps) => {
  const [session, setSession] = useState<Session | null>(null);
  const [isFirstLogin, setIsFirstLogin] = useState<boolean>(false);

  useEffect(() => {
    supabase.auth.getSession().then(({ data: { session } }) => {
      setSession(session);
    });

    const {
      data: { subscription },
    } = supabase.auth.onAuthStateChange(async (event, session) => {
      setSession(session);
    });

    const login = async (email: string, password: string): Promise<void> => {
      // Implement login logic
    };

    const logout = async (): Promise<void> => {
      // Implement logout logic
    };

    return () => subscription.unsubscribe();
  }, []);

  return (
    <AuthContext.Provider value={{ session, isFirstLogin, setIsFirstLogin }}>
      {children}
    </AuthContext.Provider>
  );
};
