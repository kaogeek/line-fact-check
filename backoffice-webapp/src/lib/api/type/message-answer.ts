type AskAnswerBase = {
  id: string;
  code: string;
  topicId: string;
  message: string;
  createDate: Date;
};

export type AskAnswer =
  | ({
      hasAnswer: true;
      answer: string;
    } & AskAnswerBase)
  | ({
      hasAnswer: false;
    } & AskAnswerBase);
