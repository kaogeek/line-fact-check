export type ChatBase = {
  sender: 'me' | 'other';
};

export interface TextChatMessage extends ChatBase {
  type: 'text';
  message: string;
}

export interface StickerChatMessage extends ChatBase {
  type: 'sticker';
  src: string;
}

export interface LoadingChatMessage {
  sender: 'other';
  type: 'loading';
}

export type ChatMessage = TextChatMessage | StickerChatMessage | LoadingChatMessage;
