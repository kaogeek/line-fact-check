export type ChatProps = {
  sender: 'me' | 'other';
};

export type ChatMessage = {
  sender: 'me' | 'other';
} & (
  | {
      type: 'text';
      message: string;
    }
  | {
      type: 'sticker';
      src: string;
    }
  | {
      type: 'loading';
    }
);
