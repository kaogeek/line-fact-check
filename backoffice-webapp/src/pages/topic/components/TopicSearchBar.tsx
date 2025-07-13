import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Search } from 'lucide-react';
import { useState } from 'react';

interface TopicTabProps {
  initCodeLike?: string;
  initMessageLike?: string;
  handleSearch: (criteria: { codeLike?: string; messageLike?: string }) => void;
}

export default function TopicSearchBar({ initCodeLike, initMessageLike, handleSearch }: TopicTabProps) {
  const [codeLike, setCodeLike] = useState<string>(initCodeLike ?? '');
  const [messageLike, setMessageLike] = useState<string>(initMessageLike ?? '');

  function handleCodeLikeChange(event: any) {
    setCodeLike(event.target.value);
  }

  function handleMessageLikeChange(event: any) {
    setMessageLike(event.target.value);
  }

  function handleSearchClick() {
    const searchCriteria = {
      ...(codeLike && { codeLike }),
      ...(messageLike && { messageLike }),
    };

    handleSearch(searchCriteria);
  }

  return (
    <div className="flex flex-col md:flex-row gap-3 w-full items-stretch">
      <div className="flex flex-1 flex-col sm:flex-row gap-2">
        <div className="min-w-[120px]">
          <label className="text-sm font-medium text-gray-700 mb-1 block">Code</label>
          <Input className="w-full" placeholder="Enter code..." value={codeLike} onChange={handleCodeLikeChange} />
        </div>
        <div className="min-w-[120px]">
          <label className="text-sm font-medium text-gray-700 mb-1 block">Message</label>
          <Input
            className="w-full"
            placeholder="Enter message..."
            value={messageLike}
            onChange={handleMessageLikeChange}
          />
        </div>
      </div>
      <div className="flex items-end">
        <Button className="h-auto px-4 py-2" onClick={handleSearchClick}>
          <Search className="h-4 w-4" /> Search
        </Button>
      </div>
    </div>
  );
}
