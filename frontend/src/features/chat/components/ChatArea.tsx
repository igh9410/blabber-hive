import React from 'react';
import styles from './ChatArea.module.scss';

export function ChatArea() {
  return (
    <div className={styles.chatArea}>
      <div className={styles.messageReceived}>
        <img
          src="/IMG_0004.PNG"
          alt="Sender's Profile"
          className={styles.profileImage}
        />
        <p>Hello! How are you?</p>
      </div>
      <div className={styles.messageSent}>
        <p>I'm good, thanks for asking!</p>
      </div>
    </div>
  );
}
