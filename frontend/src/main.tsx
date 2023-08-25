import React from 'react';
import ReactDOM from 'react-dom/client';
import './index.css';
import { AppProvider } from '@providers/AppProvider.tsx';
import { Routes } from '@providers/Routes';

ReactDOM.createRoot(document.getElementById('root')!).render(
  <React.StrictMode>
    <AppProvider>
      <Routes />
    </AppProvider>
  </React.StrictMode>
);
