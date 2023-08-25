import { RouterProvider, createBrowserRouter } from 'react-router-dom';
import { Root } from '@pages';

export function Routes() {
  const router = createBrowserRouter([
    {
      path: '/',
      element: <Root />,
    },
  ]);

  return <RouterProvider router={router} />;
}
