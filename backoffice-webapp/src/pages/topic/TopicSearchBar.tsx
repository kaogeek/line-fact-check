import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import type { GetTopicCriteria } from '@/lib/api/service/topic';
import { Search } from 'lucide-react';
import { useState } from 'react';

interface TopicTabProps {
  criteria: GetTopicCriteria;
  setCriteria: (criteria: GetTopicCriteria) => void;
}

export default function TopicSearchBar({ criteria, setCriteria }: TopicTabProps) {
  const [keyword, setKeyword] = useState<string>(criteria.keyword ?? '');

  const handleChange = (event: any) => {
    setKeyword(event.target.value);
  };

  const handleSearch = () => {
    setCriteria({
      ...criteria,
      keyword: keyword,
    });
  };

  return (
    <div className="flex gap-2">
      <Input className="flex-1" placeholder="Search keyword..." value={keyword} onChange={handleChange} />
      <Button onClick={handleSearch}>
        <Search />
      </Button>
    </div>
  );
}
