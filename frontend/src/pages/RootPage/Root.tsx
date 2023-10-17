import { SUPABASE_API_KEY, SUPABASE_URL } from '@config';
import { Auth } from '@supabase/auth-ui-react';
import { Session, createClient } from '@supabase/supabase-js';
import { ThemeSupa } from '@supabase/auth-ui-shared';
import { useState, useEffect } from 'react';
import styles from './Root.module.scss';
import { Header } from '@components/Header';
import { ChatArea, InputArea, MessageType } from '@features/chat';

const supabase = createClient(SUPABASE_URL, SUPABASE_API_KEY);

export function Root() {
  const [session, setSession] = useState<Session | null>(null);

  const [messages, setMessages] = useState<MessageType[]>([]);

  const handleNewMessage = (text: string) => {
    setMessages((prevMessages) => [
      ...prevMessages,
      { sender: 'sent', text: text },
    ]);
  };

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
    return <Auth supabaseClient={supabase} appearance={{ theme: ThemeSupa }} />;
  }
  return (
    <div className={styles.container}>
      <Header />
      <ChatArea messages={messages} />
      <InputArea onMessageSend={handleNewMessage} />
    </div>
  );
}
