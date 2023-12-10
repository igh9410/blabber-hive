import React, { useState } from 'react';
import styles from './InputArea.module.scss';
import { useWebSocketConnection } from '@hooks/useWebSocketConnection';

type InputAreaProps = {
  onMessageSend: (text: string) => void;
  chatRoomId: string;
};

export function InputArea({
  onMessageSend,
  chatRoomId,
}: Readonly<InputAreaProps>) {
  const [message, setMessage] = useState(''); // State to hold the input value
  const { sendMessage } = useWebSocketConnection(chatRoomId);
  const handleInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setMessage(e.target.value); // Update the state when input changes
  };

  const handleSendClick = () => {
    if (message.trim()) {
      // Check if the message is not just whitespace
      sendMessage(message); // Send the message
      onMessageSend(message); // Update the UI
      setMessage(''); // Clear the input after sending
    }
  };

  const handleKeyDown = (e: React.KeyboardEvent<HTMLInputElement>) => {
    if (e.key === 'Enter') {
      e.preventDefault(); // Prevent default behavior like form submission or line break
      handleSendClick();
    }
  };

  return (
    <div className={styles.inputArea}>
      <input
        type="text"
        placeholder="Message..."
        value={message}
        onChange={handleInputChange}
        onKeyDown={handleKeyDown}
      />
      <button onClick={handleSendClick}>Send</button>
    </div>
  );
}
