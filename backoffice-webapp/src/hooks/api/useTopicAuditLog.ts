import { getTopicAuditLogs } from '@/lib/api/service/topic-audit-log';
import { useQuery, type UseQueryOptions } from '@tanstack/react-query';
import type { TopicAuditLog } from '@/lib/api/type/topic-audit-log';
import type { BaseQueryOptions } from './type';

export function useGetTopicAuditLogs(topicId: string, typeId?: string[], options?: BaseQueryOptions<TopicAuditLog[]>) {
  return useQuery({
    ...options,
    queryKey: ['topic-audit-logs', topicId, typeId],
    queryFn: () => getTopicAuditLogs(topicId, typeId),
  });
}
