import React from 'react';
import styles from './Header.module.scss';

export function Header() {
  return (
    <header>
      <div className={styles.chatLogo}>Blabber-Hive</div>
      <div className={styles.chatIcons}>🔔 📞 📷</div>
    </header>
  );
}
