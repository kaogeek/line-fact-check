import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Search } from 'lucide-react';
import { useState } from 'react';

interface TopicTabProps {
  initKeyword?: string;
  handleSearch: (keyword: string) => void;
}

export default function TopicSearchBar({ initKeyword, handleSearch }: TopicTabProps) {
  const [keyword, setKeyword] = useState<string>(initKeyword ?? '');

  const handleKeywordChange = (event: any) => {
    setKeyword(event.target.value);
  };

  return (
    <div className="flex gap-2">
      <Input className="flex-1" placeholder="Search keyword..." value={keyword} onChange={handleKeywordChange} />
      <Button onClick={() => handleSearch(keyword)}>
        <Search />
      </Button>
    </div>
  );
}
