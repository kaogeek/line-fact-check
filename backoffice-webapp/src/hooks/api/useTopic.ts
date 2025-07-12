import { countTopics, getTopicById, getTopics } from '@/lib/api/service/topic';
import type { CountTopicCriteria, GetTopicCriteria } from '@/lib/api/type/topic';
import { useQuery } from '@tanstack/react-query';

export function useGetTopics(criteria: GetTopicCriteria) {
  return useQuery({
    queryKey: ['topics', criteria],
    queryFn: () => getTopics(criteria),
  });
}

export function useGetTopicById(id: string) {
  return useQuery({
    queryKey: ['topic', id],
    queryFn: () => getTopicById(id),
  });
}

export function useCountTopics(criteria: CountTopicCriteria) {
  return useQuery({
    queryKey: ['countTopics', criteria],
    queryFn: () => countTopics(criteria),
  });
}
