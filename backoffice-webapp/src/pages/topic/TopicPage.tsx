import TopicSearchBar from './components/TopicSearchBar';
import TopicPagination from './components/TopicPagination';

import TopicData from './components/TopicData';
import { useCountTopics, useGetTopics } from '@/hooks/api/useTopic';
import { useEffect, useState } from 'react';
import TabIndex from '../../components/TabIndex';
import { TopicStatus, type GetTopicCriteria, type Topic } from '@/lib/api/type/topic';
import { TYH3 } from '@/components/Typography';

export default function TopicPage() {
  const [counts, setCounts] = useState<number[]>([0, 0, 0, 0, 0]);
  const [criteria, setCriteria] = useState<GetTopicCriteria>({
    statusIn: tabs[0].statusIn,
  });
  const [activeTab, setActiveTab] = useState<number>(0);
  const { data: topics, isLoading: topicsIsLoading, error: topicsError } = useGetTopics(criteria);
  const { data: countTopics } = useCountTopics(criteria);

  useEffect(() => {
    if (!countTopics) {
      return;
    }

    const { total, pending, answered } = countTopics;

    setCounts([total, pending, answered, 0, 0]);
  }, [countTopics]);

  const handleSearch = (keyword: string) => {
    setCriteria({
      ...criteria,
      keyword,
    });

    // TODO set pagination
  };

  const handleTabChange = (activeTabIdx: number) => {
    setActiveTab(activeTabIdx);
    const tab = tabs[activeTabIdx];

    setCriteria({
      ...criteria,
      statusIn: tab.statusIn,
    });
  };

  const handleReject = (topic: Topic, idx: number) => {
    console.log(topic);
    console.log(idx);
  };

  return (
    <div className="flex flex-col gap-4 p-4 h-full">
      <TYH3>Topic</TYH3>
      <TopicSearchBar initKeyword={criteria.keyword} handleSearch={handleSearch} />
      <TabIndex activeTab={activeTab} setActiveTab={handleTabChange} tabs={tabs} counts={counts} />
      <div className="flex-1 overflow-auto">
        <TopicData isLoading={topicsIsLoading} dataList={topics} error={topicsError}></TopicData>
      </div>
      <TopicPagination />
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
