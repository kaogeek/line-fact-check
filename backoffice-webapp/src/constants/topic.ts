export const stats: Stat[] = [
  {
    label: 'Total ticket',
    value: 7,
  },
  {
    label: 'Pending',
    value: 3,
  },
  {
    label: 'Answered',
    value: 4,
  },
  {
    label: 'Rejected',
    value: 0,
  },
  {
    label: 'Approved',
    value: 0,
  },
];

export enum TopicStatus {
  PENDING = 'PENDING',
  ANSWERED = 'ANSWERED',
  REJECTED = 'REJECTED',
  APPROVED = 'APPROVED',
}

export const topics: Topic[] = [
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

export interface Stat {
  label: string;
  value: number;
}

export interface Topic {
  code: string;
  status: TopicStatus;
  description: string;
  createDate: Date;
  countOfMessageGroup: number;
  countOfTotalMessage: number;
}
