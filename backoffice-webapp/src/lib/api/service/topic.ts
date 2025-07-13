import { MOCKUP_API_LOADING_MS } from '@/constants/app';
import {
  TopicStatus,
  type CountTopic,
  type CountTopicCriteria,
  type GetTopicCriteria,
  type Topic,
} from '../type/topic';
import type { PaginationReq, PaginationRes } from '../type/base';
import { paginate } from '@/lib/utils/page-utils';

function isCodeLike(data: Topic, keyword: string): boolean {
  return data.code.toLowerCase().includes(keyword.toLowerCase());
}

function isMessageLike(data: Topic, keyword: string): boolean {
  return data.description.toLowerCase().includes(keyword.toLowerCase());
}

function isIdIn(data: Topic, idIn: string[]): boolean {
  return idIn.includes(data.id);
}

function isInStatus(data: Topic, statusIn: string[]): boolean {
  return statusIn.includes(data.status);
}

export function getTopics(criteria: GetTopicCriteria, pagination: PaginationReq): Promise<PaginationRes<Topic>> {
  console.log(criteria);
  return new Promise((resolve) => {
    setTimeout(() => {
      const { codeLike, messageLike, statusIn, idNotIn } = criteria;
      const conditions: ((data: Topic) => boolean)[] = [];

      if (codeLike) {
        conditions.push((data) => isCodeLike(data, codeLike));
      }

      if (messageLike) {
        conditions.push((data) => isMessageLike(data, messageLike));
      }

      if (statusIn) {
        conditions.push((data) => isInStatus(data, statusIn));
      }

      if (idNotIn) {
        conditions.push((data) => !isIdIn(data, idNotIn));
      }

      const filteredTopics = dataList.filter((data) => conditions.every((condition) => condition(data)));
      resolve(paginate(filteredTopics, pagination));
    }, MOCKUP_API_LOADING_MS);
  });
}

export function getTopicById(id: string): Promise<Topic | null> {
  return new Promise((resolve) => {
    setTimeout(() => {
      const topic = dataList.find((data) => data.id === id) || null;
      resolve(topic);
    }, MOCKUP_API_LOADING_MS);
  });
}

function countByCriteriaAndStatus(statusIn: TopicStatus[], criteria: CountTopicCriteria): number {
  const { codeLike, messageLike } = criteria;
  const conditions: ((data: Topic) => boolean)[] = [];

  if (codeLike) {
    conditions.push((data) => isCodeLike(data, codeLike));
  }

  if (messageLike) {
    conditions.push((data) => isMessageLike(data, messageLike));
  }

  conditions.push((data) => isInStatus(data, statusIn));

  return dataList.filter((data) => conditions.every((condition) => condition(data))).length;
}

export function countTopics(criteria: CountTopicCriteria): Promise<CountTopic> {
  return new Promise((resolve) => {
    setTimeout(() => {
      resolve({
        total: countByCriteriaAndStatus([TopicStatus.PENDING, TopicStatus.ANSWERED], criteria),
        pending: countByCriteriaAndStatus([TopicStatus.PENDING], criteria),
        answered: countByCriteriaAndStatus([TopicStatus.ANSWERED], criteria),
      });
    }, MOCKUP_API_LOADING_MS);
  });
}

export async function approveTopic(topicId: string) {
  console.log(`Approving topic ${topicId}`);
  const topic = dataList.find(t => t.id === topicId);
  if (topic) {
    topic.status = TopicStatus.APPROVED;
  }
  return new Promise((resolve) => setTimeout(resolve, MOCKUP_API_LOADING_MS));
}

export async function rejectTopic(topicId: string) {
  console.log(`Rejecting topic ${topicId}`);
  const topic = dataList.find(t => t.id === topicId);
  if (topic) {
    topic.status = TopicStatus.REJECTED;
  }
  return new Promise((resolve) => setTimeout(resolve, MOCKUP_API_LOADING_MS));
}

export const dataList: Topic[] = [
  {
    id: '1',
    code: 'T001',
    status: TopicStatus.PENDING,
    description: 'This is the first topic.',
    createDate: new Date('2023-10-01T10:00:00Z'),
    countOfMessageGroup: 3,
    countOfTotalMessage: 12,
  },
  {
    id: '2',
    code: 'T002',
    status: TopicStatus.ANSWERED,
    description: 'This is the second topic.',
    createDate: new Date('2023-10-02T14:30:00Z'),
    countOfMessageGroup: 5,
    countOfTotalMessage: 20,
  },
  {
    id: '3',
    code: 'T003',
    status: TopicStatus.REJECTED,
    description: 'This is the third topic.',
    createDate: new Date('2023-10-03T09:15:00Z'),
    countOfMessageGroup: 1,
    countOfTotalMessage: 5,
  },
  {
    id: '4',
    code: 'T004',
    status: TopicStatus.APPROVED,
    description: 'This is the fourth topic.',
    createDate: new Date('2023-10-04T11:45:00Z'),
    countOfMessageGroup: 4,
    countOfTotalMessage: 18,
  },
  {
    id: '5',
    code: 'T005',
    status: TopicStatus.PENDING,
    description: 'Topic about product feedback',
    createDate: new Date('2023-10-05T08:20:00Z'),
    countOfMessageGroup: 2,
    countOfTotalMessage: 8,
  },
  {
    id: '6',
    code: 'T006',
    status: TopicStatus.ANSWERED,
    description: 'Customer support inquiry',
    createDate: new Date('2023-10-06T13:10:00Z'),
    countOfMessageGroup: 7,
    countOfTotalMessage: 25,
  },
  {
    id: '7',
    code: 'T007',
    status: TopicStatus.REJECTED,
    description: 'Feature request discussion',
    createDate: new Date('2023-10-07T16:45:00Z'),
    countOfMessageGroup: 1,
    countOfTotalMessage: 4,
  },
  {
    id: '8',
    code: 'T008',
    status: TopicStatus.APPROVED,
    description: 'Bug report investigation',
    createDate: new Date('2023-10-08T09:30:00Z'),
    countOfMessageGroup: 3,
    countOfTotalMessage: 15,
  },
  {
    id: '9',
    code: 'T009',
    status: TopicStatus.PENDING,
    description: 'General question about service',
    createDate: new Date('2023-10-09T11:15:00Z'),
    countOfMessageGroup: 4,
    countOfTotalMessage: 12,
  },
  {
    id: '10',
    code: 'T010',
    status: TopicStatus.ANSWERED,
    description: 'Account management issue',
    createDate: new Date('2023-10-10T14:50:00Z'),
    countOfMessageGroup: 6,
    countOfTotalMessage: 22,
  },
  {
    id: '11',
    code: 'T011',
    status: TopicStatus.REJECTED,
    description: 'Billing question',
    createDate: new Date('2023-10-11T10:25:00Z'),
    countOfMessageGroup: 2,
    countOfTotalMessage: 7,
  },
  {
    id: '12',
    code: 'T012',
    status: TopicStatus.APPROVED,
    description: 'Technical support request',
    createDate: new Date('2023-10-12T15:40:00Z'),
    countOfMessageGroup: 5,
    countOfTotalMessage: 19,
  },
];
