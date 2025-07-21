import { mockApi } from '@/lib/utils/mock-api-utils';
import type { Message } from '../type/message';

function isHasTopicId(data: Message, topicId: string) {
  return data.topicId === topicId;
}

export function getMessagesByTopicId(topicId: string): Promise<Message[]> {
  return mockApi(() => {
    return dataList.filter((data) => isHasTopicId(data, topicId));
  }, 'getMessagesByTopicId');
}

export function createMessage(topicId: string, message: string): Promise<Message> {
  return mockApi(() => {
    const newMessage = {
      id: `MSG${dataList.length + 1}`.padStart(3, '0'),
      code: `MSG${dataList.length + 1}`.padStart(3, '0'),
      message,
      createDate: new Date(),
      countOfMessageGroup: 1,
      topicId,
    };
    dataList.push(newMessage);
    return newMessage;
  }, 'createMessage');
}

export const dataList: Message[] = [
  {
    id: '1',
    code: 'MSG001',
    message: 'This is message 1',
    createDate: new Date(),
    countOfMessageGroup: 1,
    topicId: '1',
  },
  {
    id: '2',
    code: 'MSG002',
    message: 'This is message 2',
    createDate: new Date(),
    countOfMessageGroup: 1,
    topicId: '1',
  },
  {
    id: '3',
    code: 'MSG003',
    message: 'This is message 3',
    createDate: new Date(),
    countOfMessageGroup: 1,
    topicId: '1',
  },
  {
    id: '4',
    code: 'MSG004',
    message: 'This is message 4',
    createDate: new Date(),
    countOfMessageGroup: 1,
    topicId: '2',
  },
  {
    id: '5',
    code: 'MSG005',
    message: 'This is message 5',
    createDate: new Date(),
    countOfMessageGroup: 1,
    topicId: '2',
  },
  {
    id: '6',
    code: 'MSG006',
    message: 'This is message 6',
    createDate: new Date(),
    countOfMessageGroup: 1,
    topicId: '3',
  },
  {
    id: '7',
    code: 'MSG007',
    message: 'This is message 7',
    createDate: new Date(),
    countOfMessageGroup: 1,
    topicId: '3',
  },
  {
    id: '8',
    code: 'MSG008',
    message: 'This is message 8',
    createDate: new Date(),
    countOfMessageGroup: 1,
    topicId: '4',
  },
  {
    id: '9',
    code: 'MSG009',
    message: 'This is message 9',
    createDate: new Date(),
    countOfMessageGroup: 1,
    topicId: '4',
  },
];
