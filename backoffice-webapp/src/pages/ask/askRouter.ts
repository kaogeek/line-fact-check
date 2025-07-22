import type { RouteObject } from 'react-router';
import AskPage from './AskPage';
import AskDetailPage from './detail/AskDetailPage';

export const askRouter: RouteObject[] = [
  {
    path: '/ask',
    children: [
      {
        index: true,
        Component: AskPage,
      },
      {
        path: ':id',
        Component: AskDetailPage,
      },
    ],
  },
];
