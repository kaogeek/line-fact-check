import { MOCKUP_API_LOADING_MS } from '@/constants/app';
import { TopicAnswerType, type TopicAnswer } from '../type/topic-answer';

function isHasTopicId(data: TopicAnswer, topicId: string) {
  return data.topicId === topicId;
}

export function getTopicAnswerByTopicId(topicId: string): Promise<TopicAnswer | null> {
  return new Promise((resolve) => {
    setTimeout(() => {
      resolve(dataList.find((data) => isHasTopicId(data, topicId)) || null);
    }, MOCKUP_API_LOADING_MS);
  });
}

export async function updateAnswer(topicId: string, answerId: string, content: string) {
  console.log(`Updating answer ${answerId} in topic ${topicId} with content: ${content}`);
  return new Promise((resolve) => setTimeout(resolve, MOCKUP_API_LOADING_MS));
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
