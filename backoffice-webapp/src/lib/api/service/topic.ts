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

function isInKeyword(data: Topic, keyword: string): boolean {
  return data.description.toLowerCase().includes(keyword.toLowerCase());
}

function isIdIn(data: Topic, idIn: string[]): boolean {
  return idIn.includes(data.id);
}

function isInStatus(data: Topic, statusIn: string[]): boolean {
  return statusIn.includes(data.status);
}

export function getTopics(criteria: GetTopicCriteria, pagination: PaginationReq): Promise<PaginationRes<Topic>> {
  return new Promise((resolve) => {
    setTimeout(() => {
      const { keyword, statusIn, idNotId } = criteria;
      const conditions: ((data: Topic) => boolean)[] = [];

      if (keyword) {
        conditions.push((data) => isInKeyword(data, keyword));
      }

      if (statusIn) {
        conditions.push((data) => isInStatus(data, statusIn));
      }

      if (idNotId) {
        conditions.push((data) => !isIdIn(data, idNotId));
      }

      const filteredTopics = dataList.filter((data) => conditions.every((condition) => condition(data)));
      resolve(paginate(filteredTopics, pagination));
    }, MOCKUP_API_LOADING_MS);
  });
}

export function getTopicById(id: string): Promise<Topic | undefined> {
  return new Promise((resolve) => {
    setTimeout(() => {
      const topic = dataList.find((data) => data.id === id);
      resolve(topic);
    }, MOCKUP_API_LOADING_MS);
  });
}

function countByCriteriaAndStatus(statusIn: TopicStatus[], criteria: CountTopicCriteria): number {
  const { keyword } = criteria;
  const conditions: ((data: Topic) => boolean)[] = [];

  if (keyword) {
    conditions.push((data) => isInKeyword(data, keyword));
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
];
