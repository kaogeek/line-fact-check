import type { ChatProps } from './type';
import { cva } from 'class-variance-authority';

export interface ChatStickerProps extends ChatProps {
  src: string;
}

const wrapperVariants = cva('w-[200px]', {
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

export default function ChatSticker({ sender, src }: ChatStickerProps) {
  return (
    <div className={wrapperVariants({ sender })}>
      <img src={src} alt="sticker" className="w-full" />
    </div>
  );
}
