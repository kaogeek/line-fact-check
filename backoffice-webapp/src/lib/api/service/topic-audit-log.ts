import { MOCKUP_API_LOADING_MS } from '@/constants/app';
import { TopicAuditLogType, type TopicAuditLog } from '../type/topic-audit-log';

function isHasTopicId(data: TopicAuditLog, topicId: string) {
  return data.topicId === topicId;
}

function isInType(data: TopicAuditLog, typeIn: string[]): boolean {
  return typeIn.includes(data.status);
}

export function getTopicAuditLogs(topicId: string, typeIn?: string[]): Promise<TopicAuditLog[]> {
  return new Promise((resolve) => {
    setTimeout(() => {
      const conditions: ((data: TopicAuditLog) => boolean)[] = [];

      conditions.push((data) => isHasTopicId(data, topicId));

      if (typeIn) {
        conditions.push((data) => isInType(data, typeIn));
      }

      const filteredLogs = dataList.filter((data) => conditions.every((condition) => condition(data)));
      resolve(filteredLogs);
    }, MOCKUP_API_LOADING_MS);
  });
}

export const dataList: TopicAuditLog[] = [
  {
    avatarUrl: '/assets/avatars/mockup/1.jpg',
    username: 'user1',
    actionDate: new Date('2023-10-01T10:00:00Z'),
    status: TopicAuditLogType.UPDATE_ANSWER,
    detail: 'Created new topic',
    topicId: '1',
  },
  {
    avatarUrl: '/assets/avatars/mockup/1.jpg',
    username: 'user1',
    actionDate: new Date('2023-10-02T10:00:00Z'),
    status: TopicAuditLogType.UPDATE_ANSWER,
    detail: 'Test test',
    topicId: '1',
  },
  {
    avatarUrl: '/assets/avatars/mockup/1.jpg',
    username: 'user1',
    actionDate: new Date('2023-10-01T12:00:00Z'),
    status: TopicAuditLogType.APPROVED,
    detail: 'Updated topic title',
    topicId: '1',
  },
  {
    avatarUrl: '/assets/avatars/mockup/1.jpg',
    username: 'user1',
    actionDate: new Date('2023-10-01T14:00:00Z'),
    status: TopicAuditLogType.UPDATE_ANSWER,
    detail: 'Added initial comments',
    topicId: '1',
  },
  {
    avatarUrl: '/assets/avatars/mockup/2.jpg',
    username: 'user2',
    actionDate: new Date('2023-10-02T11:30:00Z'),
    status: TopicAuditLogType.REJECTED,
    detail: 'Updated topic description',
    topicId: '2',
  },
  {
    avatarUrl: '/assets/avatars/mockup/2.jpg',
    username: 'user2',
    actionDate: new Date('2023-10-02T13:45:00Z'),
    status: TopicAuditLogType.APPROVED,
    detail: 'Added new section',
    topicId: '2',
  },
  {
    avatarUrl: '/assets/avatars/mockup/3.jpg',
    username: 'user3',
    actionDate: new Date('2023-10-03T12:45:00Z'),
    status: TopicAuditLogType.UPDATE_ANSWER,
    detail: 'Deleted irrelevant content',
    topicId: '3',
  },
  {
    avatarUrl: '/assets/avatars/mockup/3.jpg',
    username: 'user3',
    actionDate: new Date('2023-10-03T15:00:00Z'),
    status: TopicAuditLogType.REJECTED,
    detail: 'Updated references',
    topicId: '3',
  },
  {
    avatarUrl: '/assets/avatars/mockup/4.jpg',
    username: 'user4',
    actionDate: new Date('2023-10-04T14:15:00Z'),
    status: TopicAuditLogType.UPDATE_ANSWER,
    detail: 'Added new comments',
    topicId: '4',
  },
  {
    avatarUrl: '/assets/avatars/mockup/4.jpg',
    username: 'user4',
    actionDate: new Date('2023-10-04T16:30:00Z'),
    status: TopicAuditLogType.APPROVED,
    detail: 'Replied to comment',
    topicId: '4',
  },
];
