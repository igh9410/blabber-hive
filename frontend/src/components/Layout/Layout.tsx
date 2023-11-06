import React from 'react';
import styles from './Layout.module.scss';
import { Header } from '@components/Header';

type LayoutProps = {
  children: React.ReactNode;
};

export const Layout = ({ children }: LayoutProps) => {
  return (
    <div className={styles.layout}>
      <Header />
      {children}
    </div>
  );
};
