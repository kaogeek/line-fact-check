import { mockApi } from '@/lib/utils/mock-api-utils';
import { TopicAnswerType, type TopicAnswer } from '../type/topic-answer';

function isHasTopicId(data: TopicAnswer, topicId: string) {
  return data.topicId === topicId;
}

export function getTopicAnswerByTopicId(topicId: string): Promise<TopicAnswer | null> {
  return mockApi(() => {
    return dataList.find((data) => isHasTopicId(data, topicId)) || null;
  }, 'getTopicAnswerByTopicId');
}

export async function updateAnswer(topicId: string, answerId: string, content: string) {
  return mockApi(() => {
    console.log(`Updating answer ${answerId} in topic ${topicId} with content: ${content}`);
  }, 'updateAnswer');
}

export const dataList: TopicAnswer[] = [
  {
    answer: 'This claim has been verified by multiple independent sources.',
    type: TopicAnswerType.REAL,
    topicId: '1',
  },
  {
    answer: 'Official government reports confirm this statement is accurate.',
    type: TopicAnswerType.REAL,
    topicId: '3',
  },
  {
    answer: 'Fact-checkers have debunked this claim as false.',
    type: TopicAnswerType.FAKE,
    topicId: '4',
  },
];
