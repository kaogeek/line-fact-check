import { useTranslation } from 'react-i18next';
import TopicSearchBar from './components/TopicSearchBar';

import TopicData from './components/TopicData';
import { topicQueryKeys, useCountTopics, useGetTopics } from '@/hooks/api/topic';
import { useEffect, useState } from 'react';
import TabIndex from '../../components/TabIndex';
import { TopicStatus, type GetTopicCriteria, type Topic } from '@/lib/api/type/topic';
import { TYH3 } from '@/components/Typography';
import { rejectTopic } from '@/lib/api/service/topic';
import { useMutation, useQueryClient } from '@tanstack/react-query';
import { toast } from 'sonner';
import { useLoader } from '@/hooks/loader';
import { ConfirmAlertDialog } from '@/components/ConfirmAlertDialog';
import type { Pagination } from '@/lib/api/type/base';
import { calPagination } from '@/lib/utils/page-utils';
import PaginationControl from '@/components/PaginationControl';
import { usePaginationReqState } from '@/hooks/usePagination';

export default function TopicPage() {
  const { t } = useTranslation();

  const tabs = [
    {
      label: t('topic.total'),
      statusIn: [TopicStatus.TOPIC_PENDING, TopicStatus.TOPIC_RESOLVED],
    },
    {
      label: t('topic.status.pending'),
      statusIn: [TopicStatus.TOPIC_PENDING],
    },
    {
      label: t('topic.status.answered'),
      statusIn: [TopicStatus.TOPIC_RESOLVED],
    },
    {
      label: t('topic.status.rejected'),
      statusIn: [TopicStatus.REJECTED],
    },
    {
      label: t('topic.status.approved'),
      statusIn: [TopicStatus.APPROVED],
    },
  ];

  const [criteria, setCriteria] = useState<GetTopicCriteria>({
    statusIn: tabs[0].statusIn,
    codeLike: '',
    messageLike: '',
  });

  const [counts, setCounts] = useState<number[]>([0, 0, 0, 0, 0]);

  const { paginationReq, setPage, resetPage } = usePaginationReqState({
    page: 1,
    pageSize: 10,
  });

  const [pagination, setPagination] = useState<Pagination>({
    totalItems: 0,
    totalPages: 1,
  });

  const [activeTab, setActiveTab] = useState<number>(0);
  const [showRejectDialog, setShowRejectDialog] = useState(false);
  const [topicToReject, setTopicToReject] = useState<Topic | null>(null);
  const { data: data, isLoading, error } = useGetTopics(criteria, paginationReq);
  const { data: countTopics } = useCountTopics(criteria);
  const queryClient = useQueryClient();

  const { startLoading, stopLoading } = useLoader();

  const { mutate: rejectTopicMutation } = useMutation({
    mutationFn: (topicId: string) => rejectTopic(topicId),
    onSettled: () => {
      stopLoading();
    },
    onSuccess: () => {
      toast.success(t('topic.deleteSuccess'));
      queryClient.removeQueries({ queryKey: topicQueryKeys.all });
    },
    onError: (err) => {
      toast.error(t('topic.deleteError'));
      console.error(err);
    },
  });

  useEffect(() => {
    if (!countTopics) {
      return;
    }

    const { total, TOPIC_PENDING, TOPIC_RESOLVED } = countTopics;

    setCounts([total, TOPIC_PENDING, TOPIC_RESOLVED, 0, 0]);
    setPagination(calPagination(total, paginationReq));
  }, [countTopics, paginationReq]);

  function handleSearch(criteria: { codeLike?: string; messageLike?: string }) {
    resetPage();
    setCriteria((prev) => ({
      ...prev,
      codeLike: criteria.codeLike,
      messageLike: criteria.messageLike,
    }));
    queryClient.removeQueries({ queryKey: topicQueryKeys.all });
  }

  function handleTabChange(activeTabIdx: number) {
    setActiveTab(activeTabIdx);
    const tab = tabs[activeTabIdx];

    resetPage();
    setCriteria((prev) => ({
      ...prev,
      statusIn: tab.statusIn,
    }));
  }

  function handlePageChange(page: number) {
    setPage(page);
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
      <TYH3>{t('topic.title')}</TYH3>
      <TopicSearchBar
        initCodeLike={criteria.codeLike}
        initMessageLike={criteria.messageLike}
        handleSearch={handleSearch}
      />
      <TabIndex activeTab={activeTab} setActiveTab={handleTabChange} tabs={tabs} counts={counts} />
      <div className="flex-1 overflow-auto">
        <TopicData isLoading={isLoading} dataList={data} error={error} onReject={handleRejectClick}></TopicData>
      </div>
      <PaginationControl paginationReq={paginationReq} pagination={pagination} onPageChange={handlePageChange} />
      <ConfirmAlertDialog
        open={showRejectDialog}
        onOpenChange={setShowRejectDialog}
        title={t('topic.rejectTitle')}
        description={t('topic.rejectDescription')}
        confirmText={t('topic.rejectConfirm')}
        onConfirm={handleConfirmReject}
      />
    </div>
  );
}

export interface TopicPageTab {
  label: string;
  statusIn: TopicStatus[];
}
