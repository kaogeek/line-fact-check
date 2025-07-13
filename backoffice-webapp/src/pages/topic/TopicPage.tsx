import TopicSearchBar from './components/TopicSearchBar';

import TopicData from './components/TopicData';
import { useCountTopics, useGetTopics } from '@/hooks/api/useTopic';
import { useEffect, useState } from 'react';
import TabIndex from '../../components/TabIndex';
import { TopicStatus, type GetTopicCriteria, type Topic } from '@/lib/api/type/topic';
import { TYH3 } from '@/components/Typography';
import type { PaginationReq } from '@/lib/api/type/base';
import PaginationControl from '@/components/PaginationControl';
import { rejectTopic } from '@/lib/api/service/topic';
import { useMutation } from '@tanstack/react-query';
import { toast } from 'sonner';
import { useLoader } from '@/hooks/useLoader';
import { ConfirmAlertDialog } from '@/components/ConfirmAlertDialog';

export default function TopicPage() {
  const [counts, setCounts] = useState<number[]>([0, 0, 0, 0, 0]);
  const [criteria, setCriteria] = useState<GetTopicCriteria>({
    statusIn: tabs[0].statusIn,
    codeLike: '',
    messageLike: '',
  });
  const [paginationReq, setPaginationReq] = useState<PaginationReq>({
    page: 1,
  });
  const [activeTab, setActiveTab] = useState<number>(0);
  const [showRejectDialog, setShowRejectDialog] = useState(false);
  const [topicToReject, setTopicToReject] = useState<Topic | null>(null);
  const { data: data, isLoading, error } = useGetTopics(criteria, paginationReq);
  const { data: countTopics } = useCountTopics(criteria);

  const { startLoading, stopLoading } = useLoader();

  const { mutate: rejectTopicMutation } = useMutation({
    mutationFn: (topicId: string) => rejectTopic(topicId),
    onSettled: () => {
      stopLoading();
    },
    onSuccess: () => {
      toast.success('Topic rejected successfully');
      // Optionally refresh the data here
    },
    onError: (err) => {
      toast.error('Failed to reject topic');
      console.error(err);
    },
  });

  useEffect(() => {
    if (!countTopics) {
      return;
    }

    const { total, pending, answered } = countTopics;

    setCounts([total, pending, answered, 0, 0]);
  }, [countTopics]);

  function handleSearch(criteria: { codeLike?: string; messageLike?: string }) {
    setCriteria((prev) => {
      return { ...prev, codeLike: criteria.codeLike, messageLike: criteria.messageLike };
    });
  }

  function handleTabChange(activeTabIdx: number) {
    setActiveTab(activeTabIdx);
    const tab = tabs[activeTabIdx];

    setCriteria({
      ...criteria,
      statusIn: tab.statusIn,
    });
  }

  function handlePageChange(paginationReq: PaginationReq) {
    setPaginationReq(paginationReq);
  }

  function handleRejectClick(topic: Topic) {
    setTopicToReject(topic);
    setShowRejectDialog(true);
  }

  const handleConfirmReject = () => {
    if (topicToReject) {
      startLoading();
      rejectTopicMutation(topicToReject.id);
    }
    setShowRejectDialog(false);
  };

  return (
    <div className="flex flex-col gap-4 p-4 h-full">
      <TYH3>Topic</TYH3>
      <TopicSearchBar
        initCodeLike={criteria.codeLike}
        initMessageLike={criteria.messageLike}
        handleSearch={handleSearch}
      />
      <TabIndex activeTab={activeTab} setActiveTab={handleTabChange} tabs={tabs} counts={counts} />
      <div className="flex-1 overflow-auto">
        <TopicData isLoading={isLoading} dataList={data?.items} error={error} onReject={handleRejectClick}></TopicData>
      </div>
      <PaginationControl paginationRes={data} onPageChange={handlePageChange} />
      <ConfirmAlertDialog
        open={showRejectDialog}
        onOpenChange={setShowRejectDialog}
        title="Reject Topic"
        description="Are you sure you want to reject this topic?"
        confirmText="Reject"
        onConfirm={handleConfirmReject}
      />
    </div>
  );
}

const tabs = [
  {
    label: 'Total',
    statusIn: [TopicStatus.PENDING, TopicStatus.ANSWERED],
  },
  {
    label: 'Pending',
    statusIn: [TopicStatus.PENDING],
  },
  {
    label: 'Answered',
    statusIn: [TopicStatus.ANSWERED],
  },
  {
    label: 'Rejected',
    statusIn: [TopicStatus.REJECTED],
  },
  {
    label: 'Approved',
    statusIn: [TopicStatus.APPROVED],
  },
];

export interface TopicPageTab {
  label: string;
  statusIn: TopicStatus[];
}
