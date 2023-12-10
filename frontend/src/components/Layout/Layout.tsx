import React from 'react';
import styles from './Layout.module.scss';
import { Header } from '@components/Header';
import { Dashboard } from '@components/Dashboard/Dashboard';

type LayoutProps = {
  children: React.ReactNode;
};

export const Layout = ({ children }: LayoutProps) => {
  return (
    <div className={styles.layout}>
      <div className={styles.headerWrapper}>
        <Header />
        <Dashboard>{children}</Dashboard>
      </div>
    </div>
  );
};
