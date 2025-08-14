import { mockApi } from '@/lib/utils/mock-api-utils';
import type { MessageGroup, GetMessageGroupCriteria, CountMessageGroup } from '../type/message-group';
import { MessageGroupStatus } from '../type/message-group';
import type { PaginationReq } from '../type/base';

const MOCK_MESSAGE_GROUPS: MessageGroup[] = [
  {
    id: '1',
    status: MessageGroupStatus.MGROUP_PENDING,
    topic_id: 'topic1',
    name: 'Message Group 1',
    text: 'This is a sample message group 1',
    text_sha1: 'hash1',
    language: 'th',
    created_at: new Date().toISOString(),
    updated_at: new Date().toISOString(),
  },
  {
    id: '2',
    status: MessageGroupStatus.MGROUP_APPROVED,
    topic_id: 'topic2',
    name: 'Message Group 2',
    text: 'This is a sample message group 2',
    text_sha1: 'hash2',
    language: 'en',
    created_at: new Date(Date.now() - 1000 * 60 * 60).toISOString(),
    updated_at: new Date().toISOString(),
  },
];

export function getMessageGroups(
  criteria: GetMessageGroupCriteria,
  pagination: PaginationReq
): Promise<MessageGroup[]> {
  return mockApi(() => {
    let filtered = [...MOCK_MESSAGE_GROUPS];

    // Apply status filter
    if (criteria.statusIn && criteria.statusIn.length > 0) {
      filtered = filtered.filter((group) => criteria.statusIn?.includes(group.status));
    }

    // Apply search filter
    if (criteria.codeLike || criteria.messageLike) {
      const searchLower = criteria.codeLike?.toLowerCase() || criteria.messageLike?.toLowerCase() || '';
      filtered = filtered.filter(
        (group) => group.name.toLowerCase().includes(searchLower) || group.text.toLowerCase().includes(searchLower)
      );
    }

    // Apply pagination with defaults
    const page = pagination.page || 1;
    const pageSize = pagination.pageSize || 10;
    const start = (page - 1) * pageSize;
    const end = start + pageSize;
    return filtered.slice(start, end);
  }, 'getMessageGroups');
}

export function countMessageGroups(criteria: Omit<GetMessageGroupCriteria, 'statusIn'>): Promise<CountMessageGroup> {
  return mockApi(() => {
    let filtered = [...MOCK_MESSAGE_GROUPS];

    // Apply search filter
    if (criteria.codeLike || criteria.messageLike) {
      const searchLower = criteria.codeLike?.toLowerCase() || criteria.messageLike?.toLowerCase() || '';
      filtered = filtered.filter(
        (group) => group.name.toLowerCase().includes(searchLower) || group.text.toLowerCase().includes(searchLower)
      );
    }

    const statusCounts = filtered.reduce(
      (acc, group) => {
        acc[group.status] = (acc[group.status] || 0) + 1;
        return acc;
      },
      {} as Record<string, number>
    );

    return {
      total: filtered.length,
      MGROUP_PENDING: statusCounts[MessageGroupStatus.MGROUP_PENDING] || 0,
      MGROUP_APPROVED: statusCounts[MessageGroupStatus.MGROUP_APPROVED] || 0,
      MGROUP_ASSIGNED: statusCounts[MessageGroupStatus.MGROUP_ASSIGNED] || 0,
      MGROUP_REJECTED: statusCounts[MessageGroupStatus.MGROUP_REJECTED] || 0,
    };
  }, 'countMessageGroups');
}

export function updateMessageGroupStatus(id: string, status: MessageGroupStatus): Promise<void> {
  return mockApi(() => {
    const index = MOCK_MESSAGE_GROUPS.findIndex((g) => g.id === id);
    if (index !== -1) {
      MOCK_MESSAGE_GROUPS[index] = {
        ...MOCK_MESSAGE_GROUPS[index],
        status,
        updated_at: new Date().toISOString(),
      };
    }
  }, 'updateMessageGroupStatus');
}

export function getMessagesGroupByTopicId(topicId: string): Promise<MessageGroup[]> {
  return mockApi(() => {
    return MOCK_MESSAGE_GROUPS.filter((group) => group.topic_id === topicId);
  }, 'getMessagesGroupByTopicId');
}
