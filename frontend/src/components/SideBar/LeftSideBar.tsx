import React, { useState } from 'react';
import styles from './LeftSideBar.module.scss';

export function LeftSideBar() {
  return (
    <div className={styles.leftSideBarWrapper}>
      <div className={styles.leftSideBar}>
        {/* Add your sidebar content here */}
        <div className={styles.sidebarItem}>Item 1</div>
        <div className={styles.sidebarItem}>Item 2</div>
        <div className={styles.sidebarItem}>Item 3</div>
        {/* ... more items */}
      </div>
    </div>
  );
}
