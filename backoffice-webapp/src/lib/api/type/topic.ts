export enum TopicStatus {
  PENDING = 'PENDING',
  ANSWERED = 'ANSWERED',
  REJECTED = 'REJECTED',
  APPROVED = 'APPROVED',
}

type TopicStatusOptionSpec = {
  // TODO use badge variant instead or something
  variant: 'warning' | 'blue' | 'danger' | 'success';
  label: string;
};

type TopicStatusOption = Record<TopicStatus, TopicStatusOptionSpec>;

export const topicStatusOption: TopicStatusOption = {
  PENDING: {
    variant: 'warning',
    label: 'Pending',
  },
  ANSWERED: {
    variant: 'blue',
    label: 'Answered',
  },
  REJECTED: {
    variant: 'danger',
    label: 'Rejected',
  },
  APPROVED: {
    variant: 'success',
    label: 'Approved',
  },
};

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
