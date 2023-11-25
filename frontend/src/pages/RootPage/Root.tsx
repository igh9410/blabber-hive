import { Auth } from '@supabase/auth-ui-react';
import { ThemeSupa } from '@supabase/auth-ui-shared';
import { Layout } from '@components/Layout/Layout';
import { Outlet } from 'react-router-dom';
import styles from './Root.module.scss';
import { supabase, useAuth } from '@providers';

export function Root() {
  const { session } = useAuth();

  if (!session) {
    return (
      <div className={styles.wrapperContainer}>
        <div className={styles.authContainer}>
          <Auth
            supabaseClient={supabase}
            providers={['google', 'apple', 'github']}
            appearance={{ theme: ThemeSupa }}
          />
        </div>
      </div>
    );
  }
  return (
    <Layout>
      <Outlet />
    </Layout>
  );
}
