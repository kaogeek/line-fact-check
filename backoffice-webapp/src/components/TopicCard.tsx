import { Card, CardContent } from './ui/card';
import { TYH3, TYMuted, TYSmall } from './Typography';
import type { TopicStatus } from '@/pages/topic/TopicPage';
import { Badge } from './ui/badge';
import { formatDate } from '@/formatter/DateFormatter';

interface TopicCardProps {
  code: string;
  status: TopicStatus;
  createDate: Date;
  description: string;
}

export default function TopicCard({ code, status, createDate, description }: TopicCardProps) {
  return (
    <Card>
      <CardContent className="flex flex-col">
        <div className="flex items-center gap-2 mb-1">
          <TYH3>{code}</TYH3>
          <Badge variant="secondary">{status}</Badge>
        </div>

        <TYMuted className="mb-3">Create at: {formatDate(createDate)}</TYMuted>
        <TYSmall>{description}</TYSmall>
      </CardContent>
    </Card>
  );
}
