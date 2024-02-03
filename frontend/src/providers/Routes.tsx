import {
  Navigate,
  Outlet,
  RouterProvider,
  createBrowserRouter,
  redirect,
} from 'react-router-dom';
import { Root } from '@pages';
import { ChatRoom } from '@pages/ChatRoomPage';
import { SignUp } from '@pages/SignUpPage';
import { RenewLobby } from '@pages/RenewLobbyPage';
import { useUsers } from '@hooks';

const ProtectedRoutes = () => {
  const { data: userData, isLoading } = useUsers();
  if (isLoading) {
    return <div>Loading...</div>;
  }

  return userData ? <Outlet /> : <Navigate to="/signup" />;
};

export function Routes() {
  const router = createBrowserRouter([
    {
      path: '/',
      element: <Root />,
      children: [
        {
          element: <ProtectedRoutes />, // Wrap the children with ProtectedRoute
          children: [
            {
              index: true,
              element: <RenewLobby />,
            },
            {
              path: 'chats/:id',
              element: <ChatRoom />,
            },
          ],
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
