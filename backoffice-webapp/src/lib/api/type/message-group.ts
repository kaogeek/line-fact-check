export enum MessageGroupStatus {
  MGROUP_PENDING = 'MGROUP_PENDING',
  MGROUP_APPROVED = 'MGROUP_APPROVED',
  MGROUP_ASSIGNED = 'MGROUP_ASSIGNED',
  MGROUP_REJECTED = 'MGROUP_REJECTED',
}

type MessageGroupStatusOptionSpec = {
  // TODO use badge variant instead or something
  variant: 'warning' | 'blue' | 'danger' | 'success';
  label: string;
};

type MessageGroupStatusOption = Record<MessageGroupStatus, MessageGroupStatusOptionSpec>;

export const messageGroupStatusOption: MessageGroupStatusOption = {
  MGROUP_PENDING: {
    variant: 'warning',
    label: 'messageGroup.status.pending',
  },
  MGROUP_APPROVED: {
    variant: 'blue',
    label: 'messageGroup.status.approved',
  },
  MGROUP_ASSIGNED: {
    variant: 'success',
    label: 'messageGroup.status.assigned',
  },
  MGROUP_REJECTED: {
    variant: 'danger',
    label: 'messageGroup.status.rejected',
  },
};

export interface MessageGroup {
  id: string;
  status: MessageGroupStatus;
  topic_id: string;
  name: string;
  text: string;
  text_sha1: string;
  language: string;
  created_at: string;
  updated_at: string;
}
export interface GetMessageGroupCriteria {
  statusIn?: MessageGroupStatus[];
  codeLike?: string;
  messageLike?: string;
}

export interface CountMessageGroup {
  total: number;
  MGROUP_PENDING: number;
  MGROUP_APPROVED: number;
  MGROUP_ASSIGNED: number;
  MGROUP_REJECTED: number;
}
