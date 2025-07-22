import { mockApi } from '@/lib/utils/mock-api-utils';
import type { AskAnswer } from '../type/message-answer';

export function askMessage(message: string): Promise<AskAnswer> {
  return mockApi(() => {
    const answer = dataList.find((data) => data.message === message);

    if (!answer) {
      return dataList[1];
    }

    return answer;
  }, 'askMessage');
}

export function getAskAnswerById(messageId: string): Promise<AskAnswer | null> {
  return mockApi(() => {
    return dataList.find((data) => data.id === messageId) || null;
  }, 'getAskMessageAnswerById');
}

const dataList: AskAnswer[] = [
  {
    id: '1',
    code: `ANS001`.padStart(3, '0'),
    message: 'hello',
    createDate: new Date(),
    hasAnswer: true,
    topicId: '1',
    answer: 'This is answer',
  },
  {
    id: '2',
    code: `ANS002`.padStart(3, '0'),
    message: 'User ask',
    createDate: new Date(),
    hasAnswer: false,
    topicId: '2',
  },
];
