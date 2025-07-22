import { useQuery } from '@tanstack/react-query';
import type { BaseQueryOptions } from './type';
import type { AskAnswer } from '@/lib/api/type/message-answer';
import { getAskAnswerById } from '@/lib/api/service/message-answer';

const domainKey = 'messageAnswers';

export const messageQueryKeys = {
  all: [domainKey] as const,
  useGetAskAnswerById: (id: string) => [domainKey, id] as const,
};

export function useGetAskAnswerById(id: string, options?: BaseQueryOptions<AskAnswer | null>) {
  return useQuery({
    ...options,
    queryKey: messageQueryKeys.useGetAskAnswerById(id),
    queryFn: () => getAskAnswerById(id),
  });
}
