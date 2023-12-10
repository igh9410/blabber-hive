import { RouterProvider, createBrowserRouter } from 'react-router-dom';
import { Root } from '@pages';
import { ChatRoom } from '@pages/ChatRoomPage';
import { SignUp } from '@pages/SignUpPage';
import { RenewLobby } from '@pages/RenewLobbyPage';

export function Routes() {
  const router = createBrowserRouter([
    {
      path: '/',
      element: <Root />,
      children: [
        {
          path: '/',
          element: <RenewLobby />,
        },
        {
          path: '/chats/:id',
          element: <ChatRoom />,
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
