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
    label: 'topic.status.pending',
  },
  ANSWERED: {
    variant: 'blue',
    label: 'topic.status.answered',
  },
  REJECTED: {
    variant: 'danger',
    label: 'topic.status.rejected',
  },
  APPROVED: {
    variant: 'success',
    label: 'topic.status.approved',
  },
};

export interface Stat {
  label: string;
  value: number;
}

export interface Topic {
  id: string;
  code: string;
  status: TopicStatus;
  description: string;
  createDate: Date;
  countOfMessageGroup: number;
  countOfTotalMessage: number;
}

export interface CountTopic {
  total: number;
  pending: number;
  answered: number;
}

export interface CountTopicCriteria {
  codeLike?: string;
  messageLike?: string;
}

export interface GetTopicCriteria extends CountTopicCriteria {
  idNotIn?: string[];
  statusIn?: string[];
}
