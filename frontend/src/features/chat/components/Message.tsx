import React from 'react';
import styles from './Message.module.scss';

type MessageProps = {
  sender: string;
  text: string;
  img?: string;
};

export const Message = ({ sender, text, img }: MessageProps) => {
  return (
    <div
      className={
        sender === 'received' ? styles.messageReceived : styles.messageSent
      }
    >
      <p>{text}</p>
    </div>
  );
};
