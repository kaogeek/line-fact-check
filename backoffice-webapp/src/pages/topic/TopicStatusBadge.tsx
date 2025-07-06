import { Badge } from '@/components/ui/badge';
import { topicStatusOption, type TopicStatus } from '@/lib/api/type/topic';

interface TopicStatusBadgeProps {
  status: TopicStatus;
}

export default function TopicStatusBadge({ status }: TopicStatusBadgeProps) {
  const { variant, label } = topicStatusOption[status];
  return <Badge variant={variant}>{label}</Badge>;
}
