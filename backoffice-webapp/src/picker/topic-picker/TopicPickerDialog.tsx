import { Dialog, DialogContent, DialogDescription, DialogHeader, DialogTitle } from '@/components/ui/dialog';
import { useGetTopics } from '@/hooks/api/useTopic';
import TopicPickerData from './components/TopicPickerData';
import { useMemo, useState } from 'react';
import type { GetTopicCriteria } from '@/lib/api/type/topic';

interface TopicPickerDialogProps {
  open?: boolean;
  onOpenChange?(open: boolean): void;
  currentId?: string;
  onChoose: (topicId: string) => void;
}

export default function TopicPickerDialog({ open, onOpenChange, currentId, onChoose }: TopicPickerDialogProps) {
  const [criteria, setCriteria] = useState<GetTopicCriteria>({
    idNotId: currentId ? [currentId] : undefined,
  });

  const {
    isLoading,
    data: dataList,
    error,
  } = useGetTopics(criteria, {
    enabled: open,
  });

  function onHandleChoose(topicId: string) {
    onOpenChange && onOpenChange(false);
    onChoose(topicId);
  }

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Answer history</DialogTitle>
          <DialogDescription asChild>
            <TopicPickerData
              isLoading={isLoading}
              dataList={dataList}
              error={error}
              onChoose={onHandleChoose}
            ></TopicPickerData>
          </DialogDescription>
        </DialogHeader>
      </DialogContent>
    </Dialog>
  );
}
