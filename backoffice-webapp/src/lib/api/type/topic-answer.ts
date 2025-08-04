export enum TopicAnswerType {
  REAL = 'REAL',
  FAKE = 'FAKE',
}
export interface TopicAnswer {
  id: string;
  user_id: string;
  topic_id: string;
  text: string;
  type: TopicAnswerType;
  created_at: string;
}
