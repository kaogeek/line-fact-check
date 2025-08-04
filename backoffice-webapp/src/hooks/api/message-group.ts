import { useQuery } from '@tanstack/react-query';
import type { BaseQueryOptions } from './type';
import type { MessageGroup } from '@/lib/api/type/message-group';
import { getMessagesGroupByTopicId } from '@/lib/api/service/message-group';

const domainKey = 'messagesGroup';

export const messageQueryKeys = {
  all: [domainKey] as const,
  byTopicId: (id: string) => [domainKey, id] as const,
};

export function useGetMessageGroupsByTopicId(id: string, options?: BaseQueryOptions<MessageGroup[]>) {
  return useQuery({
    ...options,
    queryKey: messageQueryKeys.byTopicId(id),
    queryFn: () => getMessagesGroupByTopicId(id),
  });
}
