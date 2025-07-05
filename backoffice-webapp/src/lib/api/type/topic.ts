export enum TopicStatus {
  PENDING = 'PENDING',
  ANSWERED = 'ANSWERED',
  REJECTED = 'REJECTED',
  APPROVED = 'APPROVED',
}

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
