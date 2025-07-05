import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Search } from 'lucide-react';

export default function TopicSearchBar() {
  return (
    <div className="flex gap-2 p-4">
      <Input className="flex-1" placeholder="Search keyword..." />
      <Button>
        <Search />
      </Button>
    </div>
  );
}
