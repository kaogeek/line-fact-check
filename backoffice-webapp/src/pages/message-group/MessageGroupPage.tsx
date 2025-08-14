import { useTranslation } from 'react-i18next';
import { useEffect, useState } from 'react';
import { TYH3 } from '@/components/Typography';
import { usePaginationReqState } from '@/hooks/usePagination';
import { MessageGroupStatus } from '@/lib/api/type/message-group';
import { messageGroupQueryKeys, useCountMessageGroups, useGetMessageGroups } from '@/hooks/api/message-group';
import type { GetMessageGroupCriteria } from '@/lib/api/type/message-group';
import TabIndex from '@/components/TabIndex';
import MessageGroupSearchBar from './components/MessageGroupSearchBar';
import MessageGroupData from './components/MessageGroupData';
import type { Pagination } from '@/lib/api/type/base';
import { calPagination } from '@/lib/utils/page-utils';
import { useQueryClient } from '@tanstack/react-query';
import PaginationControl from '@/components/PaginationControl';

export default function MessageGroupPage() {
  const { t } = useTranslation();

  const tabs: MessageGroupPageTab[] = [
    {
      label: t('common.total'),
      statusIn: [MessageGroupStatus.MGROUP_PENDING, MessageGroupStatus.MGROUP_APPROVED],
    },
    {
      label: t('messageGroup.status.pending'),
      statusIn: [MessageGroupStatus.MGROUP_PENDING],
    },
    {
      label: t('messageGroup.status.assigned'),
      statusIn: [MessageGroupStatus.MGROUP_ASSIGNED],
    },
    {
      label: t('messageGroup.status.rejected'),
      statusIn: [MessageGroupStatus.MGROUP_REJECTED],
    },
    {
      label: t('messageGroup.status.approved'),
      statusIn: [MessageGroupStatus.MGROUP_APPROVED],
    },
  ];

  const [criteria, setCriteria] = useState<GetMessageGroupCriteria>({
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

  const { data: data, isLoading, error } = useGetMessageGroups(criteria, paginationReq);
  const { data: countMessageGroups } = useCountMessageGroups(criteria);
  const queryClient = useQueryClient();

  useEffect(() => {
    if (!countMessageGroups) {
      return;
    }

    const { total, MGROUP_PENDING, MGROUP_APPROVED, MGROUP_ASSIGNED, MGROUP_REJECTED } = countMessageGroups;

    setCounts([total, MGROUP_PENDING, MGROUP_APPROVED, MGROUP_ASSIGNED, MGROUP_REJECTED]);
    setPagination(calPagination(total, paginationReq));
  }, [countMessageGroups, paginationReq]);

  function handleSearch(criteria: { codeLike?: string; messageLike?: string }) {
    resetPage();
    setCriteria((prev) => ({
      ...prev,
      codeLike: criteria.codeLike,
      messageLike: criteria.messageLike,
    }));
    queryClient.removeQueries({ queryKey: messageGroupQueryKeys.all });
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

  return (
    <div className="flex flex-col gap-4 p-4 h-full">
      <TYH3>{t('messageGroup.title')}</TYH3>
      <MessageGroupSearchBar initCodeLike={criteria.codeLike} handleSearch={handleSearch} />
      <TabIndex activeTab={activeTab} setActiveTab={handleTabChange} tabs={tabs} counts={counts} />
      <div className="flex-1 overflow-auto">
        <MessageGroupData isLoading={isLoading} dataList={data} error={error} />
      </div>
      <PaginationControl paginationReq={paginationReq} pagination={pagination} onPageChange={handlePageChange} />
    </div>
  );
}

export interface MessageGroupPageTab {
  label: string;
  statusIn: MessageGroupStatus[];
}
