import { APP_NAME } from '@/constants/app';
import ChatScreen from '@/components/chat/ChatScreen';
import type { ChatMessage } from './type';
import { useState } from 'react';
import { useMutation } from '@tanstack/react-query';
import { askMessage } from '@/lib/api/service/message-answer';
import { TYH4 } from '@/components/Typography';

export default function AskPage() {
  const [messages, setMessages] = useState<ChatMessage[]>([
    {
      sender: 'other',
      type: 'text',
      message: `
      👋 สวัสดีครับ!
ยินดีต้อนรับสู่ แชทบอทตรวจสอบข้อเท็จจริง (Fact Check Bot) 🔍
ถ้าคุณสงสัยว่า ข่าวนี้จริงไหม? หรือ ข้อมูลนี้น่าเชื่อถือหรือเปล่า?
พิมพ์เข้ามาได้เลยครับ
      `,
    },
    {
      sender: 'other',
      type: 'sticker',
      src: '/assets/chat/stickers/hi.png',
    },
  ]);
  const { mutate: askMessageMutation } = useMutation({
    mutationFn: (message: string) => askMessage(message),
  });

  function handleSendMessage(message: string) {
    // Add loading message
    setMessages((prev) => [
      ...prev,
      {
        sender: 'me',
        type: 'text',
        message,
      },
    ]);

    setMessages((prev) => [
      ...prev,
      {
        sender: 'other',
        type: 'loading',
      },
    ]);

    // Simulate response (replace with actual API call)
    askMessageMutation(message, {
      onSettled: () => {
        setMessages((prev) => [...prev.filter((m) => m.type !== 'loading')]);
      },
      onSuccess: (data) => {
        if (!data.hasAnswer) {
          setMessages((prev) => [
            ...prev,
            {
              sender: 'other',
              type: 'sticker',
              src: '/assets/chat/stickers/waiting.png',
            },
          ]);

          const trackingLink = `${window.location.origin}/ask/${data.id}`;

          setMessages((prev) => [
            ...prev,
            {
              sender: 'other',
              type: 'text',
              message: `
              ตอนนี้ยังไม่มีคำตอบสำหรับข่าวนี้นะครับ 📰 ทีมงานกำลังตามข้อมูลเพิ่มเติมอยู่ แล้วจะรีบแจ้งให้ทราบทันทีเมื่อได้คำตอบครับ 
              สามารถติดตามสถานะได้ที่ลิงค์นี้ : <a href="${trackingLink}">${trackingLink}</a>`,
              isRichText: true,
            },
          ]);

          return;
        }

        setMessages((prev) => [
          ...prev,
          {
            sender: 'other',
            type: 'text',
            message: data.answer,
          },
        ]);
      },
    });
  }

  return (
    <div className="flex flex-col h-screen">
      <div className="p-4 bg-secondary text-secondary-foreground w-full">
        <TYH4>{APP_NAME}</TYH4>
      </div>
      <ChatScreen className="flex-1 min-h-0" initialMessages={messages} onSendMessage={handleSendMessage} />
    </div>
  );
}
