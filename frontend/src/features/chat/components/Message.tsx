import React from 'react';
import styles from './Message.module.scss';

type MessageProps = {
  sender: 'received' | 'sent';
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
      {sender === 'received' && (
        <img src={img} alt="Sender's Profile" className="profileImage" />
      )}
      <p>{text}</p>
    </div>
  );
};
