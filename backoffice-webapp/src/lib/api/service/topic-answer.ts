import { TopicAnswerType, type TopicAnswer } from '../type/topic-answer';
import apiClient from '../client';

export async function getTopicAnswerByTopicId(topicId: string): Promise<TopicAnswer | null> {
  const response = await apiClient.get<TopicAnswer | null>(`/topics/${topicId}/answer`);
  return patchData(response.data);
}

function patchData(topic: TopicAnswer | null): TopicAnswer | null {
  // TODO resolve this with real code
  if (!topic) {
    return null;
  }

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
