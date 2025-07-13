export enum TopicAnswerType {
  REAL = 'REAL',
  FAKE = 'FAKE',
}

export interface TopicAnswer {
  answer: string;
  type: TopicAnswerType;
  topicId: string;
}
