import { useEffect, useState } from 'react';
import type { ChatMessage } from '../type';
import { TYH4 } from '@/components/Typography';
import ChatScreen from '@/components/chat/ChatScreen';
import { useGetAskAnswerById } from '@/hooks/api/message-answer';
import { Navigate, useParams } from 'react-router';
import ErrorState from '@/components/state/ErrorState';
import LoadingState from '@/components/state/LoadingState';

export default function AskDetailPage() {
  const { id: idParams } = useParams();
  const id = idParams!;
  const { isLoading, data: answer, error } = useGetAskAnswerById(id);
  const [messages, setMessages] = useState<ChatMessage[]>([]);

  useEffect(() => {
    if (answer) {
      setMessages([
        {
          sender: 'me',
          type: 'text',
          message: answer.message,
        },
        {
          sender: 'other',
          type: 'text',
          message: answer.hasAnswer ? answer.answer : 'รอคำตอบ',
        },
      ]);
    }
  }, [answer]);

  return (
    <div className="flex flex-col h-screen">
      <div className="p-4 bg-secondary text-secondary-foreground w-full">
        <TYH4>Code: {answer?.code}</TYH4>
      </div>
      {isLoading ? (
        <LoadingState className="flex-1" />
      ) : error ? (
        <ErrorState className="flex-1" />
      ) : !answer ? (
        <Navigate to="/404" replace />
      ) : (
        <ChatScreen className="flex-1 min-h-0" initialMessages={messages} hideMessageInput />
      )}
    </div>
  );
}
