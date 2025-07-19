import './i18n';
import { createBrowserRouter, Navigate, RouterProvider } from 'react-router';
import RootLayout from './layouts/RootLayout';
import { topicRouter } from './pages/topic/topicRouter';
import DashboardPage from './pages/dashboard/DashboardPage';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import NotFoundPage from './pages/404';
import { LoaderProvider } from './hooks/loader';
import { Toaster } from 'sonner';

const router = createBrowserRouter([
  {
    path: '/',
    Component: RootLayout,
    children: [
      {
        index: true,
        element: <Navigate to="/topic" replace />,
      },
      ...topicRouter,
      {
        path: '/dashboard',
        Component: DashboardPage,
      },
    ],
  },
  {
    path: '/404',
    Component: NotFoundPage,
  },
  {
    path: '*',
    element: <Navigate to="/404" replace />,
  },
]);

const queryClient = new QueryClient();

export default function App() {
  return (
    <LoaderProvider>
      <QueryClientProvider client={queryClient}>
        <RouterProvider router={router} />
        <Toaster />
      </QueryClientProvider>
    </LoaderProvider>
  );
}
