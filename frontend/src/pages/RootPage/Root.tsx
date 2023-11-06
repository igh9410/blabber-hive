import { SUPABASE_API_KEY, SUPABASE_URL } from '@config';
import { Auth } from '@supabase/auth-ui-react';
import { Session, createClient } from '@supabase/supabase-js';
import { ThemeSupa } from '@supabase/auth-ui-shared';
import { useState, useEffect } from 'react';
import { Layout } from '@components/Layout/Layout';
import { Outlet } from 'react-router-dom';
import styles from './Root.module.scss';

const supabase = createClient(SUPABASE_URL, SUPABASE_API_KEY);

export function Root() {
  const [session, setSession] = useState<Session | null>(null);

  useEffect(() => {
    supabase.auth.getSession().then(({ data: { session } }) => {
      setSession(session);
    });

    const {
      data: { subscription },
    } = supabase.auth.onAuthStateChange((_event, session) => {
      setSession(session);
    });

    return () => subscription.unsubscribe();
  }, []);

  if (!session) {
    return (
      <div className={styles.authContainer}>
        <Auth
          supabaseClient={supabase}
          providers={['google', 'apple', 'kakao', 'github']}
          appearance={{ theme: ThemeSupa }}
        />
      </div>
    );
  }
  return (
    <Layout>
      <Outlet />
    </Layout>
  );
}
