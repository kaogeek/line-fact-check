import type { RouteObject } from 'react-router';
import MessageGroupPage from './MessageGroupPage';

export const messageGroupRouter: RouteObject[] = [
  {
    path: '/message-group',
    children: [
      {
        index: true,
        Component: MessageGroupPage,
      },
    ],
  },
];
