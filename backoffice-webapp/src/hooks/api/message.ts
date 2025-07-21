import { getMessagesByTopicId } from '@/lib/api/service/message';
import type { Message } from '@/lib/api/type/message';
import { useQuery } from '@tanstack/react-query';
import type { BaseQueryOptions } from './type';

const domainKey = 'messages';

export const messageQueryKeys = {
  all: [domainKey] as const,
  byTopicId: (id: string) => [domainKey, id] as const,
};

export function useGetMessageByTopicId(id: string, options?: BaseQueryOptions<Message[]>) {
  return useQuery({
    ...options,
    queryKey: messageQueryKeys.byTopicId(id),
    queryFn: () => getMessagesByTopicId(id),
  });
}
