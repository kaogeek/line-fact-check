import { countTopics, getTopicById, getTopics } from '@/lib/api/service/topic';
import type { CountTopic, CountTopicCriteria, GetTopicCriteria, Topic } from '@/lib/api/type/topic';
import { useQuery } from '@tanstack/react-query';
import type { BaseQueryOptions } from './type';
import type { PaginationReq } from '@/lib/api/type/base';
import { countKey } from './key';

const domainKey = 'topics';

export const topicQueryKeys = {
  all: [domainKey] as const,
  list: (criteria: GetTopicCriteria, pagination: PaginationReq) => [domainKey, criteria, pagination] as const,
  detail: (id: string) => [domainKey, id] as const,
  count: (criteria: CountTopicCriteria) => [domainKey, countKey, criteria] as const,
};

export function useGetTopics(
  criteria: GetTopicCriteria,
  pagination: PaginationReq,
  options?: BaseQueryOptions<Topic[]>
) {
  return useQuery({
    ...options,
    queryKey: topicQueryKeys.list(criteria, pagination),
    queryFn: () => getTopics(criteria, pagination),
  });
}

export function useGetTopicById(id: string, options?: BaseQueryOptions<Topic | null>) {
  return useQuery({
    ...options,
    queryKey: topicQueryKeys.detail(id),
    queryFn: () => getTopicById(id),
  });
}

export function useCountTopics(criteria: CountTopicCriteria, options?: BaseQueryOptions<CountTopic>) {
  return useQuery({
    ...options,
    queryKey: topicQueryKeys.count(criteria),
    queryFn: () => countTopics(criteria),
  });
}
