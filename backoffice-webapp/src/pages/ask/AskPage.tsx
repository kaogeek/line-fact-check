import { APP_NAME } from '@/constants/app';
import ChatScreen from '@/components/chat/ChatScreen';
import type { ChatMessage } from './type';
import { useState } from 'react';

export default function AskPage() {
  const [messages, setMessages] = useState<ChatMessage[]>([
    {
      sender: 'me',
      type: 'text',
      message: 'sampleTextMessage',
    },
    {
      sender: 'other',
      type: 'sticker',
      src: 'assets/chat/stickers/query.png',
    },
    {
      sender: 'other',
      type: 'text',
      message: 'eiei',
    },
  ]);

  const handleSendMessage = (message: string) => {
    // Add user message
    setMessages((prev) => [
      ...prev,
      {
        sender: 'me' as const,
        type: 'text' as const,
        message,
      },
    ]);

    // Add loading message
    setMessages((prev) => [
      ...prev,
      {
        sender: 'other' as const,
        type: 'loading' as const,
      },
    ]);

    // Simulate response (replace with actual API call)
    setTimeout(() => {
      setMessages((prev) => [
        ...prev.filter((m) => m.type !== 'loading'),
        {
          sender: 'other' as const,
          type: 'text' as const,
          message: `Response to: ${message}`,
        },
      ]);
    }, 500);
  };

  return (
    <div className="flex flex-col h-screen justify-center items-center">
      <div className="p-4 bg-secondary text-secondary-foreground w-full">{APP_NAME}</div>
      <ChatScreen initialMessages={messages} onSendMessage={handleSendMessage} />
    </div>
  );
}
