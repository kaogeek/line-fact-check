import { getMessagesByTopicId } from '@/lib/api/service/message';
import type { Message } from '@/lib/api/type/message';
import { useQuery } from '@tanstack/react-query';
import type { BaseQueryOptions } from './type';

export function useGetMessageByTopicId(id: string, options?: BaseQueryOptions<Message[]>) {
  return useQuery({
    ...options,
    queryKey: ['messages', id],
    queryFn: () => getMessagesByTopicId(id),
  });
}
