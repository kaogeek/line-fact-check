import { Dialog, DialogContent, DialogDescription, DialogHeader, DialogTitle } from '@/components/ui/dialog';
import { useGetTopics } from '@/hooks/api/useTopic';
import TopicPickerData from './components/TopicPickerData';
import { useState } from 'react';
import type { GetTopicCriteria } from '@/lib/api/type/topic';
import type { PaginationReq } from '@/lib/api/type/base';
import PaginationControl from '@/components/PaginationControl';

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
  const [paginationReq, setPaginationReq] = useState<PaginationReq>({
    page: 1,
  });

  const { isLoading, data, error } = useGetTopics(criteria, paginationReq, {
    enabled: open,
  });

  function handleChoose(topicId: string) {
    onOpenChange && onOpenChange(false);
    onChoose(topicId);
  }

  function handlePageChange(paginationReq: PaginationReq) {
    setPaginationReq(paginationReq);
  }

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Answer history</DialogTitle>
          <DialogDescription asChild>
            <div className="flex flex-col gap-4">
              <TopicPickerData
                isLoading={isLoading}
                dataList={data?.items}
                error={error}
                onChoose={handleChoose}
              ></TopicPickerData>
              <PaginationControl paginationRes={data} onPageChange={handlePageChange} />
            </div>
          </DialogDescription>
        </DialogHeader>
      </DialogContent>
    </Dialog>
  );
}
