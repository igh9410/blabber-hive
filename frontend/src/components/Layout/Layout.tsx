import React from 'react';
import styles from './Layout.module.scss';
import { Header } from '@components/Header';
import { LeftSideBar } from '@components/SideBar';

type LayoutProps = {
  children: React.ReactNode;
};

export const Layout = ({ children }: LayoutProps) => {
  return (
    <div className={styles.layout}>
      <div className={styles.wrapper}>
        <LeftSideBar />
        <div className={styles.headerWrapper}>
          <Header />

          {children}
        </div>
      </div>
    </div>
  );
};
