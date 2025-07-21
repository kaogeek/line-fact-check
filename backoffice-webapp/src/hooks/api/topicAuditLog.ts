import { getTopicAuditLogs } from '@/lib/api/service/topic-audit-log';
import { useQuery } from '@tanstack/react-query';
import type { TopicAuditLog } from '@/lib/api/type/topic-audit-log';
import type { BaseQueryOptions } from './type';

const domainKey = 'topicAuditLogs';

export const topicAuditLogQueryKeys = {
  all: [domainKey] as const,
  byTopicId: (topicId: string, typeId?: string[]) => [domainKey, topicId, typeId] as const,
};

export function useGetTopicAuditLogs(topicId: string, typeId?: string[], options?: BaseQueryOptions<TopicAuditLog[]>) {
  return useQuery({
    ...options,
    queryKey: topicAuditLogQueryKeys.byTopicId(topicId, typeId),
    queryFn: () => getTopicAuditLogs(topicId, typeId),
  });
}
