export enum TopicStatus {
  TOPIC_PENDING = 'TOPIC_PENDING',
  TOPIC_RESOLVED = 'TOPIC_RESOLVED',
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
  TOPIC_PENDING: {
    variant: 'warning',
    label: 'topic.status.pending',
  },
  TOPIC_RESOLVED: {
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
  status: TopicStatus;
  result: string;
  description: string;
  created_at: Date;
  replied_at: string;
  updated_at: string;

  // not have in backend
  code: string;
  countOfMessageGroup: number;
  countOfTotalMessage: number;
}

export interface CountTopic {
  total: number;
  TOPIC_PENDING: number;
  TOPIC_RESOLVED: number;
}

export interface CountTopicCriteria {
  idNotIn?: string[];
  codeLike?: string;
  messageLike?: string;
}

export interface GetTopicCriteria extends CountTopicCriteria {
  statusIn?: string[];
}
