import { useQuery } from '@tanstack/react-query';
import type { TopicAnswer } from '@/lib/api/type/topic-answer';
import type { BaseQueryOptions } from './type';
import { getTopicAnswerByTopicId } from '@/lib/api/service/topic-answer';

const domainKey = 'topicAnswers';

export const topicAnswerQueryKeys = {
  all: [domainKey] as const,
  byTopicId: (topicId: string) => [domainKey, topicId] as const,
};

export function useGetTopicAnswerByTopicId(topicId: string, options?: BaseQueryOptions<TopicAnswer | null>) {
  return useQuery({
    ...options,
    queryKey: topicAnswerQueryKeys.byTopicId(topicId),
    queryFn: () => getTopicAnswerByTopicId(topicId),
  });
}
