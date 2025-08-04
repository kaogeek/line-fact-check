import apiClient from '../client';
import type { MessageGroup } from '../type/message-group';

export async function getMessagesGroupByTopicId(topicId: string): Promise<MessageGroup[]> {
  const response = await apiClient.get<MessageGroup[]>(`/topics/${topicId}/message-group`);
  return response.data;
}
