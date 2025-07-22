import { cn } from '@/lib/utils';
import { cva } from 'class-variance-authority';
import type { ChatProps } from './type';

export interface ChatBubbleProps extends ChatProps {
  message: string;
  isTyping?: boolean;
  isRichText?: boolean;
  className?: string;
}

const wrapperVariants = cva('w-fit max-w-[70%] min-w-0', {
  variants: {
    sender: {
      me: 'self-end',
      other: 'self-start',
    },
  },
  defaultVariants: {
    sender: 'other',
  },
});

const chatBubbleVariants = cva('p-4', {
  variants: {
    sender: {
      me: 'rounded-bl-2xl rounded-br-2xl rounded-tl-2xl bg-secondary text-secondary-foreground',
      other: 'rounded-br-2xl rounded-bl-2xl rounded-tr-2xl bg-primary text-primary-foreground',
    },
  },
  defaultVariants: {
    sender: 'other',
  },
});

export default function ChatBubble({ sender, message, isTyping = false, className, isRichText }: ChatBubbleProps) {
  return (
    <div className={cn(wrapperVariants({ sender, className }))}>
      <div className={cn(chatBubbleVariants({ sender }))}>
        {isTyping ? (
          <div className="flex gap-1">
            <div
              className="w-2 h-2 rounded-full bg-current opacity-60 animate-bounce"
              style={{ animationDelay: '0ms' }}
            />
            <div
              className="w-2 h-2 rounded-full bg-current opacity-60 animate-bounce"
              style={{ animationDelay: '150ms' }}
            />
            <div
              className="w-2 h-2 rounded-full bg-current opacity-60 animate-bounce"
              style={{ animationDelay: '300ms' }}
            />
          </div>
        ) : isRichText ? (
          <div dangerouslySetInnerHTML={{ __html: message }} />
        ) : (
          message
        )}
      </div>
    </div>
  );
}
