import SectionTopic from '@/components/SectionTopic';
import TopicCard from '@/components/TopicCard';
import TopicSearchBar from './TopicSearchBar';
import TopicPagination from './TopicPagination';
import { stats, topics } from '@/constants/topic';
import { Tabs, TabsList, TabsTrigger } from '@/components/ui/tabs';
import { Badge } from '@/components/ui/badge';
import TopicData from './TopicData';

export default function TopicPage() {
  return (
    <div className="flex flex-col h-full">
      <SectionTopic label="Topic" />
      <TopicSearchBar />
      <div className="p-4">
        <Tabs defaultValue="0" className="w-[400px]">
          <TabsList>
            {stats.map((stat, idx) => (
              <TabsTrigger key={idx} value={idx.toString()}>
                <div className="flex gap-2">
                  <span>{stat.label}</span>
                  {stat.value > 0 && (
                    <Badge variant="secondary" className="rounded-full">
                      {stat.value}
                    </Badge>
                  )}
                </div>
              </TabsTrigger>
            ))}
          </TabsList>
        </Tabs>
      </div>
      {/* <SectionTopic label={`Topic (${topics.length})`} /> */}
      <div className="flex-1 flex flex-col p-4 gap-4 overflow-auto">
        <TopicData dataList={topics}></TopicData>
      </div>
      <TopicPagination />
    </div>
  );
}
