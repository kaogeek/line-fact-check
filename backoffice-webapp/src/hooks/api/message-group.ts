import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import type { BaseQueryOptions } from './type';
import type {
  MessageGroup,
  GetMessageGroupCriteria,
  MessageGroupStatus,
  CountMessageGroup,
} from '@/lib/api/type/message-group';
import {
  getMessageGroups,
  getMessagesGroupByTopicId,
  countMessageGroups,
  updateMessageGroupStatus,
} from '@/lib/api/service/message-group';
import type { PaginationReq } from '@/lib/api/type/base';

const domainKey = 'messageGroups';

export const messageGroupQueryKeys = {
  all: [domainKey] as const,
  lists: () => [...messageGroupQueryKeys.all, 'list'] as const,
  list: (criteria: GetMessageGroupCriteria, pagination: PaginationReq) =>
    [...messageGroupQueryKeys.lists(), { ...criteria, ...pagination }] as const,
  counts: () => [...messageGroupQueryKeys.all, 'count'] as const,
  count: (criteria: Omit<GetMessageGroupCriteria, 'statusIn'>) =>
    [...messageGroupQueryKeys.counts(), criteria] as const,
  byTopicId: (id: string) => [domainKey, 'topic', id] as const,
};

export function useGetMessageGroups(
  criteria: GetMessageGroupCriteria,
  pagination: PaginationReq,
  options?: BaseQueryOptions<MessageGroup[]>
) {
  return useQuery({
    ...options,
    queryKey: messageGroupQueryKeys.list(criteria, pagination),
    queryFn: () => getMessageGroups(criteria, pagination),
  });
}

export function useCountMessageGroups(
  criteria: GetMessageGroupCriteria,
  options?: BaseQueryOptions<CountMessageGroup>
) {
  return useQuery({
    ...options,
    queryKey: messageGroupQueryKeys.count(criteria),
    queryFn: () => countMessageGroups(criteria),
  });
}

export function useUpdateMessageGroupStatus() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ id, status }: { id: string; status: MessageGroupStatus }) => updateMessageGroupStatus(id, status),
    onSuccess: () => {
      // Invalidate all message group queries to refetch data
      queryClient.invalidateQueries({ queryKey: messageGroupQueryKeys.all });
    },
  });
}

export function useGetMessageGroupsByTopicId(id: string, options?: BaseQueryOptions<MessageGroup[]>) {
  return useQuery({
    ...options,
    queryKey: messageGroupQueryKeys.byTopicId(id),
    queryFn: () => getMessagesGroupByTopicId(id),
  });
}
