import { TopicAnswerType, type TopicAnswer } from '../type/topic-answer';
import apiClient from '../client';

export async function getTopicAnswerByTopicId(topicId: string): Promise<TopicAnswer | null> {
  const response = await apiClient.get<TopicAnswer | null>(`/topics/${topicId}/answer`);
  const data = response.data;

  if (!data) {
    return null;
  }

  return patchData(data);
}

function patchData(topic: TopicAnswer): TopicAnswer {
  return {
    ...topic,
    type: TopicAnswerType.REAL,
  };
}

export async function updateAnswer(topicId: string, content: string) {
  await apiClient.post<TopicAnswer | null>(`/admin/topics/resolve/${topicId}`, {
    text: content,
  });
}
