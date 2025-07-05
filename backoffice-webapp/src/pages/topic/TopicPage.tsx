import SectionTopic from '@/components/SectionTopic';
import StatCard from '@/components/StatCard';
import TopicCard from '@/components/TopicCard';
import TopicSearchBar from './TopicSearchBar';
import TopicPagination from './TopicPagination';

export default function TopicPage() {
  return (
    <div className="flex flex-col">
      <SectionTopic label="Topic" />
      <div className="p-4 grid grid-cols-1 md:grid-cols-3 lg:grid-cols-5 gap-4">
        {/* TODO resolve with convert number to string */}
        {stats.map((stat, idx) => (
          <StatCard key={idx} label={stat.label} value={stat.value.toString()}></StatCard>
        ))}
      </div>
      <TopicSearchBar />
      <SectionTopic label={`Topic (${topics.length})`} />
      <div className="flex flex-col p-4 gap-4">
        {topics.map((topic, idx) => (
          <TopicCard
            key={idx}
            code={topic.code}
            status={topic.status}
            createDate={topic.createDate}
            description={topic.description}
          ></TopicCard>
        ))}
      </div>
      <TopicPagination />
    </div>
  );
}

export const stats: Stat[] = [
  {
    label: 'Total ticket',
    value: 9,
  },
  {
    label: 'Pending',
    value: 3,
  },
  {
    label: 'Answered',
    value: 2,
  },
  {
    label: 'Rejected',
    value: 2,
  },
  {
    label: 'Approved',
    value: 2,
  },
];

export enum TopicStatus {
  PENDING = 'PENDING',
  ANSWERED = 'ANSWERED',
  REJECTED = 'REJECTED',
  APPROVED = 'APPROVED',
}

export const topics: Topic[] = [
  {
    code: 'T001',
    status: TopicStatus.PENDING,
    description: 'This is the first topic.',
    createDate: new Date('2023-10-01T10:00:00Z'),
    countOfMessageGroup: 3,
    countOfTotalMessage: 12,
  },
  {
    code: 'T002',
    status: TopicStatus.ANSWERED,
    description: 'This is the second topic.',
    createDate: new Date('2023-10-02T14:30:00Z'),
    countOfMessageGroup: 5,
    countOfTotalMessage: 20,
  },
  {
    code: 'T003',
    status: TopicStatus.REJECTED,
    description: 'This is the third topic.',
    createDate: new Date('2023-10-03T09:15:00Z'),
    countOfMessageGroup: 1,
    countOfTotalMessage: 5,
  },
  {
    code: 'T004',
    status: TopicStatus.APPROVED,
    description: 'This is the fourth topic.',
    createDate: new Date('2023-10-04T11:45:00Z'),
    countOfMessageGroup: 4,
    countOfTotalMessage: 18,
  },
];

export interface Stat {
  label: string;
  value: number;
}

export interface Topic {
  code: string;
  status: TopicStatus;
  description: string;
  createDate: Date;
  countOfMessageGroup: number;
  countOfTotalMessage: number;
}
