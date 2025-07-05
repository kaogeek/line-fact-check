import { Tabs, TabsList, TabsTrigger } from '@/components/ui/tabs';
import { Badge } from '@/components/ui/badge';
import type { GetTopicCriteria } from '@/lib/api/service/topic';
import type { TopicPageTab } from './TopicPage';

interface TopicTabProps {
  activeTab: string;
  setActiveTab: (activeTab: string) => void;
  criteria: GetTopicCriteria;
  setCriteria: (criteria: GetTopicCriteria) => void;
  tabs: TopicPageTab[];
  counts: number[];
}

export default function TopicTab({ activeTab, setActiveTab, criteria, setCriteria, tabs, counts }: TopicTabProps) {
  const handleTabChange = (currentTab: string) => {
    setActiveTab(currentTab);
    const tabIdx = Number(currentTab);
    const tab = tabs[tabIdx];

    setCriteria({
      ...criteria,
      statusIn: tab.statusIn,
    });
  };

  return (
    <Tabs value={activeTab} onValueChange={handleTabChange}>
      <TabsList>
        {tabs.map((stat, idx) => (
          <TabsTrigger key={idx} value={idx.toString()}>
            <div className="flex gap-2">
              <span>{stat.label}</span>
              {counts[idx] > 0 && (
                <Badge variant="secondary" className="rounded-full">
                  {counts[idx]}
                </Badge>
              )}
            </div>
          </TabsTrigger>
        ))}
      </TabsList>
    </Tabs>
  );
}
