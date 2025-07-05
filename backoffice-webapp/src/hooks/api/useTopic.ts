import {
  countTopics,
  getTopics,
  type CountTopic,
  type CountTopicCriteria,
  type GetTopicCriteria,
} from '@/lib/api/service/topic';
import type { Topic } from '@/lib/api/type/topic';
import { useMemo } from 'react';

export function useGetTopics(criteria: GetTopicCriteria): Topic[] {
  return useMemo(() => getTopics(criteria), [criteria]);
}

export function useCountTopics(criteria: CountTopicCriteria): CountTopic {
  return useMemo(() => countTopics(criteria), [criteria]);
}
