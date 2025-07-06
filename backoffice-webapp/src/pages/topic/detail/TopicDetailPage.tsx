import { TYH3, TYMuted } from '@/components/Typography';
import { Navigate, useParams } from 'react-router';
import TopicStatusBadge from '../components/TopicStatusBadge';
import { useGetTopicById } from '@/hooks/api/useTopic';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu';
import { EllipsisVertical } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { formatDate } from '@/formatter/date-formatter';
import TopicMessageDetail from './components/TopicMessageDetail';
import { useGetMessageByTopicId } from '@/hooks/api/userMessage';
import TopicMessageAnswer from './components/TopicMessageAnswer';

export default function TopicDetailPage() {
  const { id } = useParams();

  if (!id) {
    return <Navigate to="/404" replace />;
  }

  const topic = useGetTopicById(id);

  if (!topic) {
    return <Navigate to="/404" replace />;
  }

  const messages = useGetMessageByTopicId(topic.id);

  return (
    <div className="flex flex-col gap-4 p-4 h-full">
      <div className="flex flex-col">
        <div className="flex gap-2">
          <TYH3 className="flex-1">Topic: {topic.code}</TYH3>
          <TopicStatusBadge status={topic.status} />
          <DropdownMenu>
            <DropdownMenuTrigger asChild>
              <Button variant="outline">
                <EllipsisVertical />
              </Button>
            </DropdownMenuTrigger>
            <DropdownMenuContent>
              <DropdownMenuItem>Approve</DropdownMenuItem>
              <DropdownMenuItem>Reject</DropdownMenuItem>
            </DropdownMenuContent>
          </DropdownMenu>
        </div>
        <TYMuted>Create at: {formatDate(topic.createDate)}</TYMuted>
      </div>
      <TopicMessageDetail dataList={messages} />
      <TopicMessageAnswer />
    </div>
  );
}
