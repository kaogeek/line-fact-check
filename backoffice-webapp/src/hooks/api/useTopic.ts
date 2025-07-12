import { countTopics, getTopicById, getTopics } from '@/lib/api/service/topic';
import type { CountTopic, CountTopicCriteria, GetTopicCriteria, Topic } from '@/lib/api/type/topic';
import { useQuery } from '@tanstack/react-query';
import type { BaseQueryOptions } from './type';
import type { PaginationReq, PaginationRes } from '@/lib/api/type/base';

export function useGetTopics(
  criteria: GetTopicCriteria,
  pagination: PaginationReq,
  options?: BaseQueryOptions<PaginationRes<Topic>>
) {
  return useQuery({
    ...options,
    queryKey: ['topics', criteria, pagination],
    queryFn: () => getTopics(criteria, pagination),
  });
}

export function useGetTopicById(id: string, options?: BaseQueryOptions<Topic | undefined>) {
  return useQuery({
    ...options,
    queryKey: ['topic', id],
    queryFn: () => getTopicById(id),
  });
}

export function useCountTopics(criteria: CountTopicCriteria, options?: BaseQueryOptions<CountTopic>) {
  return useQuery({
    ...options,
    queryKey: ['countTopics', criteria],
    queryFn: () => countTopics(criteria),
  });
}
