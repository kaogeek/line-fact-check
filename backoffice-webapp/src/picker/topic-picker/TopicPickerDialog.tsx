import { Dialog, DialogContent, DialogHeader, DialogTitle } from '@/components/ui/dialog';
import { useCountTopics, useGetTopics } from '@/hooks/api/topic';
import TopicPickerData from './components/TopicPickerData';
import { useEffect, useState } from 'react';
import type { GetTopicCriteria } from '@/lib/api/type/topic';
import type { Pagination } from '@/lib/api/type/base';
import PaginationControl from '@/components/PaginationControl';
import TopicSearchBar from '@/pages/topic/components/TopicSearchBar';
import { usePaginationReqState } from '@/hooks/usePagination';
import { calPagination } from '@/lib/utils/page-utils';

interface TopicPickerDialogProps {
  title?: string;
  open?: boolean;
  onOpenChange?(open: boolean): void;
  currentId?: string;
  onChoose: (topicId: string) => void;
}

export default function TopicPickerDialog({
  title = 'Select topic',
  open,
  onOpenChange,
  currentId,
  onChoose,
}: TopicPickerDialogProps) {
  const [criteria, setCriteria] = useState<GetTopicCriteria>({
    idNotIn: currentId ? [currentId] : undefined,
  });

  const [, setCount] = useState<number>(0);

  const [pagination, setPagination] = useState<Pagination>({
    totalItems: 0,
    totalPages: 1,
  });

  const { paginationReq, setPage } = usePaginationReqState({
    /* TODO: make this to default value */
    page: 1,
    pageSize: 10,
  });

  const { isLoading, data, error } = useGetTopics(criteria, paginationReq, {
    enabled: open,
  });

  const { data: countTopics } = useCountTopics(criteria);

  useEffect(() => {
    if (!countTopics) {
      return;
    }

    const { total } = countTopics;

    setCount(total);
    setPagination(calPagination(total, paginationReq));
  }, [countTopics, paginationReq]);

  function handleChoose(topicId: string) {
    if (onOpenChange) {
      onOpenChange(false);
    }

    onChoose(topicId);
  }

  function handlePageChange(page: number) {
    setPage(page);
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
          <DialogTitle>{title}</DialogTitle>
        </DialogHeader>
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
              dataList={data}
              error={error}
              onChoose={handleChoose}
            />
          </div>

          <PaginationControl paginationReq={paginationReq} pagination={pagination} onPageChange={handlePageChange} />
        </div>
      </DialogContent>
    </Dialog>
  );
}
