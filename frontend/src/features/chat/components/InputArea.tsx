import React, { useState } from 'react';
import styles from './InputArea.module.scss';
import { useWebSocketConnection } from '@hooks/useWebSocketConnection';

export function InputArea() {
  const [message, setMessage] = useState(''); // State to hold the input value
  const { sendMessage } = useWebSocketConnection();
  const handleInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setMessage(e.target.value); // Update the state when input changes
  };

  const handleSendClick = () => {
    if (message.trim()) {
      // Check if the message is not just whitespace
      sendMessage(message); // Send the message
      setMessage(''); // Clear the input after sending
    }
  };

  return (
    <div className={styles.inputArea}>
      <input
        type="text"
        placeholder="Type a message..."
        value={message}
        onChange={handleInputChange}
      />
      <button onClick={handleSendClick}>Send</button>
    </div>
  );
}
