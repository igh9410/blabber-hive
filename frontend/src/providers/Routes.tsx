import { RouterProvider, createBrowserRouter } from 'react-router-dom';
import { Root } from '@pages';
import { Lobby } from '@pages/LobbyPage';
import { SignUp } from '@pages/SignUpPage';

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
        {
          path: '/signup',
          element: <SignUp />,
        },
      ],
    },
  ]);

  return <RouterProvider router={router} />;
}
