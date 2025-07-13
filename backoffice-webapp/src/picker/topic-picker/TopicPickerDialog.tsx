import { Dialog, DialogContent, DialogDescription, DialogHeader, DialogTitle } from '@/components/ui/dialog';
import { useGetTopics } from '@/hooks/api/useTopic';
import TopicPickerData from './components/TopicPickerData';
import { useState } from 'react';
import type { GetTopicCriteria } from '@/lib/api/type/topic';
import type { PaginationReq } from '@/lib/api/type/base';
import PaginationControl from '@/components/PaginationControl';
import TopicSearchBar from '@/pages/topic/components/TopicSearchBar';

interface TopicPickerDialogProps {
  open?: boolean;
  onOpenChange?(open: boolean): void;
  currentId?: string;
  onChoose: (topicId: string) => void;
}

export default function TopicPickerDialog({ open, onOpenChange, currentId, onChoose }: TopicPickerDialogProps) {
  const [criteria, setCriteria] = useState<GetTopicCriteria>({
    idNotIn: currentId ? [currentId] : undefined,
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

  function handleSearch(criteria: { codeLike?: string; messageLike?: string }) {
    setCriteria((prev) => {
      return { ...prev, codeLike: criteria.codeLike, messageLike: criteria.messageLike };
    });
  }

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="max-w-[95vw] lg:max-w-[95vw] max-h-[80vh] flex flex-col">
        <DialogHeader>
          <DialogTitle>Answer history</DialogTitle>
        </DialogHeader>
        <DialogDescription asChild></DialogDescription>
        <div className="flex flex-col gap-4 flex-1 min-h-0">
          <TopicSearchBar
            initCodeLike={criteria.codeLike}
            initMessageLike={criteria.messageLike}
            handleSearch={handleSearch}
          />

          <div className="flex-1 min-h-0 flex flex-col">
            <TopicPickerData
              className="flex-1 overflow-auto"
              isLoading={isLoading}
              dataList={data?.items}
              error={error}
              onChoose={handleChoose}
            />
          </div>

          <PaginationControl paginationRes={data} onPageChange={handlePageChange} />
        </div>
      </DialogContent>
    </Dialog>
  );
}
