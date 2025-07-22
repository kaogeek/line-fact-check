import { useState, useRef, useEffect } from 'react';
import ChatBubble from '@/components/chat/ChatBubble';
import ChatSticker from '@/components/chat/ChatSticker';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { SendHorizontal } from 'lucide-react';
import type { ChatMessage } from './type';
import { cn } from '@/lib/utils';

export function RenderChatMessage({ message }: { message: ChatMessage }) {
  switch (message.type) {
    case 'text':
      return <ChatBubble sender={message.sender} message={message.message} isRichText={message.isRichText} />;
    case 'sticker':
      return <ChatSticker sender={message.sender} src={message.src} />;
    case 'loading':
      return <ChatBubble sender={message.sender} isTyping message="loading" />;
  }
}

type ChatScreenProps = {
  initialMessages?: ChatMessage[];
  onSendMessage?: (message: string) => void;
  className?: string;
  hideMessageInput?: boolean;
};

export default function ChatScreen({
  initialMessages = [],
  onSendMessage,
  className,
  hideMessageInput = false,
}: ChatScreenProps) {
  const [messageInput, setMessageInput] = useState<string>('');
  const messagesEndRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' });
  }, [initialMessages]);

  return (
    <div className={cn('flex flex-col w-full', className)}>
      <div className="flex flex-col flex-1 gap-4 p-4 w-full overflow-y-auto">
        {initialMessages.map((message, index) => (
          <RenderChatMessage key={index} message={message} />
        ))}
        <div ref={messagesEndRef} />
      </div>
      {!hideMessageInput && (
        <div className="w-full flex gap-2 p-4">
          <Input
            className="flex-1"
            value={messageInput}
            onChange={(e) => setMessageInput(e.target.value)}
            onKeyDown={(e) => {
              if (e.key === 'Enter' && messageInput.trim() && onSendMessage) {
                onSendMessage(messageInput);
                setMessageInput('');
              }
            }}
          />
          <Button
            variant="secondary"
            size="icon"
            onClick={() => {
              if (messageInput.trim() && onSendMessage) {
                onSendMessage(messageInput);
                setMessageInput('');
              }
            }}
          >
            <SendHorizontal />
          </Button>
        </div>
      )}
    </div>
  );
}
