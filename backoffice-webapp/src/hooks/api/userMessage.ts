import { getMessagesByTopicId } from '@/lib/api/service/message';
import { useQuery } from '@tanstack/react-query';

export function useGetMessageByTopicId(id: string) {
  return useQuery({
    queryKey: ['messages', id],
    queryFn: () => getMessagesByTopicId(id),
  });
}
