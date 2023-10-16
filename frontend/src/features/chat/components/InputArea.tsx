import React from 'react';
import styles from './InputArea.module.scss';

export function InputArea() {
  return (
    <div className={styles.inputArea}>
      <input type="text" placeholder="Type a message..." />
      <button>Send</button>
    </div>
  );
}
