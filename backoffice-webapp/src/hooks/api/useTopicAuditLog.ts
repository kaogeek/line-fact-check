import { getTopicAuditLogs } from '@/lib/api/service/topic-audit-log';
import { useQuery } from '@tanstack/react-query';

export function useGetTopicAuditLogs(topicId: string, typeId?: string[]) {
  return useQuery({
    queryKey: ['topic-audit-logs', topicId, typeId],
    queryFn: () => getTopicAuditLogs(topicId, typeId),
  });
}
