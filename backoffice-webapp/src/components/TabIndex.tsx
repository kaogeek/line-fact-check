import { Tabs, TabsList, TabsTrigger } from '@/components/ui/tabs';
import { Badge } from '@/components/ui/badge';
import type { TopicPageTab } from '../pages/topic/TopicPage';

interface TabIndexProps {
  tabs: TopicPageTab[];
  counts: number[];
  activeTab: number;
  setActiveTab: (activeTab: number) => void;
}

export default function TabIndex({ activeTab, setActiveTab, tabs, counts }: TabIndexProps) {
  const handleTabChange = (currentTab: string) => {
    const tabIdx = Number(currentTab);

    setActiveTab(tabIdx);
  };

  return (
    <Tabs value={activeTab.toString()} onValueChange={handleTabChange}>
      <TabsList>
        {tabs.map((stat, idx) => (
          <TabsTrigger key={idx} value={idx.toString()}>
            <div className="flex gap-2">
              <span>{stat.label}</span>
              {counts[idx] > 0 && (
                <Badge variant="strongWarning" className="rounded-full">
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
