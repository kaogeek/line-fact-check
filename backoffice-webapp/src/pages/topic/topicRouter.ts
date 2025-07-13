import type { RouteObject } from 'react-router';
import TopicPage from './TopicPage';
import TopicDetailPage from './detail/TopicDetailPage';

export const topicRouter: RouteObject[] = [
  {
    path: '/topic',
    children: [
      {
        index: true,
        Component: TopicPage,
      },
      {
        path: ':id',
        Component: TopicDetailPage,
      },
    ],
  },
];
