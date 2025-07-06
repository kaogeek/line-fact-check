import { StrictMode } from 'react';
import { createRoot } from 'react-dom/client';
import './index.css';
import { createBrowserRouter, RouterProvider } from 'react-router';
import RootLayout from './layouts/RootLayout.tsx';
import DashboardPage from './pages/dashboard/DashboardPage.tsx';
import HomePage from './pages/home/HomePage.tsx';
import { topicRouter } from './pages/topic/topicRouter.ts';

const router = createBrowserRouter([
  {
    path: '/',
    Component: RootLayout,
    children: [
      {
        index: true,
        Component: HomePage,
      },
      ...topicRouter,
      {
        path: '/dashboard',
        Component: DashboardPage,
      },
    ],
  },
]);

createRoot(document.getElementById('root')!).render(
  <StrictMode>
    <RouterProvider router={router} />
  </StrictMode>
);
