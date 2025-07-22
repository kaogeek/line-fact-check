import type { RouteObject } from 'react-router';
import AskPage from './AskPage';

export const askRouter: RouteObject[] = [
  {
    path: '/ask',
    children: [
      {
        index: true,
        Component: AskPage,
      },
    ],
  },
];
