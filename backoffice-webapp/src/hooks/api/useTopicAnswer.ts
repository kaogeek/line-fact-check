import { useQuery } from '@tanstack/react-query';
import type { TopicAnswer } from '@/lib/api/type/topic-answer';
import type { BaseQueryOptions } from './type';
import { getTopicAnswerByTopicId } from '@/lib/api/service/topic-answer';

export function useGetTopicAnswerByTopicId(topicId: string, options?: BaseQueryOptions<TopicAnswer | null>) {
  return useQuery({
    ...options,
    queryKey: ['topic-answer-by-topic-id', topicId],
    queryFn: () => getTopicAnswerByTopicId(topicId),
  });
}
