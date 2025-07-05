import { TopicStatus, type Topic } from '../type/topic';

function isInKeyword(data: Topic, keyword: string): boolean {
  return data.description.toLowerCase().includes(keyword.toLowerCase());
}

function isInStatus(data: Topic, statusIn: string[]): boolean {
  return statusIn.includes(data.status);
}

export function getTopics(criteria: GetTopicCriteria): Topic[] {
  const { keyword, statusIn } = criteria;
  const conditions: ((data: Topic) => boolean)[] = [];

  if (keyword) {
    conditions.push((data) => isInKeyword(data, keyword));
  }

  if (statusIn) {
    conditions.push((data) => isInStatus(data, statusIn));
  }

  return dataList.filter((data) => conditions.every((condition) => condition(data)));
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

export function countTopics(criteria: CountTopicCriteria): CountTopic {
  return {
    total: countByCriteriaAndStatus([TopicStatus.PENDING, TopicStatus.ANSWERED], criteria),
    pending: countByCriteriaAndStatus([TopicStatus.PENDING], criteria),
    answered: countByCriteriaAndStatus([TopicStatus.ANSWERED], criteria),
  };
}

export interface CountTopicCriteria {
  keyword?: string;
}

export interface CountTopic {
  total: number;
  pending: number;
  answered: number;
}

export interface GetTopicCriteria extends CountTopicCriteria {
  statusIn?: string[];
}

export const dataList: Topic[] = [
  {
    code: 'T001',
    status: TopicStatus.PENDING,
    description: 'This is the first topic.',
    createDate: new Date('2023-10-01T10:00:00Z'),
    countOfMessageGroup: 3,
    countOfTotalMessage: 12,
  },
  {
    code: 'T002',
    status: TopicStatus.ANSWERED,
    description: 'This is the second topic.',
    createDate: new Date('2023-10-02T14:30:00Z'),
    countOfMessageGroup: 5,
    countOfTotalMessage: 20,
  },
  {
    code: 'T003',
    status: TopicStatus.REJECTED,
    description: 'This is the third topic.',
    createDate: new Date('2023-10-03T09:15:00Z'),
    countOfMessageGroup: 1,
    countOfTotalMessage: 5,
  },
  {
    code: 'T004',
    status: TopicStatus.APPROVED,
    description: 'This is the fourth topic.',
    createDate: new Date('2023-10-04T11:45:00Z'),
    countOfMessageGroup: 4,
    countOfTotalMessage: 18,
  },
];
