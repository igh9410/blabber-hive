import { RouterProvider, createBrowserRouter } from 'react-router-dom';
import { Root } from '@pages';
import { Lobby } from '@pages/LobbyPage';

export function Routes() {
  const router = createBrowserRouter([
    {
      path: '/',
      element: <Root />,
      children: [
        {
          path: '/',
          element: <Lobby />,
        },
      ],
    },
  ]);

  return <RouterProvider router={router} />;
}
