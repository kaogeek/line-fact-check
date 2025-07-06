import { getMessagesByTopicId } from '@/lib/api/service/message';
import type { Message } from '@/lib/api/type/message';
import { useMemo } from 'react';

export function useGetMessageByTopicId(id: string): Message[] {
  return useMemo(() => getMessagesByTopicId(id), [id]);
}
