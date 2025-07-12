import { countTopics, getTopicById, getTopics } from '@/lib/api/service/topic';
import type { CountTopic, CountTopicCriteria, GetTopicCriteria, Topic } from '@/lib/api/type/topic';
import { useQuery, type UseQueryOptions } from '@tanstack/react-query';
import type { BaseQueryOptions } from './type';

export function useGetTopics(criteria: GetTopicCriteria, options?: BaseQueryOptions<Topic[]>) {
  return useQuery({
    ...options,
    queryKey: ['topics', criteria],
    queryFn: () => getTopics(criteria),
  });
}

export function useGetTopicById(id: string, options?: BaseQueryOptions<Topic | undefined>) {
  return useQuery({
    ...options,
    queryKey: ['topic', id],
    queryFn: () => getTopicById(id),
  });
}

export function useCountTopics(criteria: CountTopicCriteria, options?: UseQueryOptions<CountTopic, Error, CountTopic>) {
  return useQuery({
    ...options,
    queryKey: ['countTopics', criteria],
    queryFn: () => countTopics(criteria),
  });
}
