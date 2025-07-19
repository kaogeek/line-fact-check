export enum TopicAuditLogType {
  UPDATE_ANSWER = 'UPDATE_ANSWER',
  APPROVED = 'APPROVED',
  REJECTED = 'REJECTED',
}

type TopicAuditLogTypeSpec = {
  actionDescription: string;
};

type TopicAuditLogTypeOption = Record<TopicAuditLogType, TopicAuditLogTypeSpec>;

export const topicAuditLogTypeOption: TopicAuditLogTypeOption = {
  UPDATE_ANSWER: {
    actionDescription: 'topicAuditLog.updateAnswer',
  },
  APPROVED: {
    actionDescription: 'topicAuditLog.approved',
  },
  REJECTED: {
    actionDescription: 'topicAuditLog.rejected',
  },
};

export interface TopicAuditLog {
  avatarUrl: string;
  username: string;
  actionDate: Date;
  status: TopicAuditLogType;
  detail: string;
  topicId: string;
}
