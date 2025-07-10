import { getTopicAuditLogs } from '@/lib/api/service/topic-audit-log';
import type { TopicAuditLog } from '@/lib/api/type/topic-audit-log';
import { useMemo } from 'react';

export function useGetTopicAuditLogs(topicId: string, typeId?: string[]): TopicAuditLog[] {
  return useMemo(() => getTopicAuditLogs(topicId, typeId), [topicId, typeId]);
}
