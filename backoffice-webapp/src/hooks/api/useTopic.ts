import { countTopics, getTopicById, getTopics } from '@/lib/api/service/topic';
import type { CountTopic, CountTopicCriteria, GetTopicCriteria, Topic } from '@/lib/api/type/topic';
import { useMemo } from 'react';

export function useGetTopics(criteria: GetTopicCriteria): Topic[] {
  return useMemo(() => getTopics(criteria), [criteria]);
}

export function useGetTopicById(id: string): Topic | undefined {
  return useMemo(() => getTopicById(id), [id]);
}

export function useCountTopics(criteria: CountTopicCriteria): CountTopic {
  return useMemo(() => countTopics(criteria), [criteria]);
}
