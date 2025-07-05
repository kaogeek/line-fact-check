import { StrictMode } from 'react';
import { createRoot } from 'react-dom/client';
import './index.css';
import { createBrowserRouter, RouterProvider } from 'react-router';
import RootLayout from './layouts/RootLayout.tsx';
import DashboardPage from './pages/dashboard/DashboardPage.tsx';
import HomePage from './pages/home/HomePage.tsx';
import TopicPage from './pages/topic/TopicPage.tsx';

const router = createBrowserRouter([
  {
    path: '/',
    element: <RootLayout />,
    children: [
      {
        index: true,
        element: <HomePage />,
      },
      {
        path: '/topic',
        element: <TopicPage />,
      },
      {
        path: '/dashboard',
        element: <DashboardPage />,
      },
    ],
  },
]);

createRoot(document.getElementById('root')!).render(
  <StrictMode>
    <RouterProvider router={router} />
  </StrictMode>
);
