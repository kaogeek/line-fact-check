import { Badge } from '@/components/ui/badge';
import { topicStatusOption, type TopicStatus } from '@/lib/api/type/topic';
import { useTranslation } from 'react-i18next';

interface TopicStatusBadgeProps {
  status: TopicStatus;
}

export default function TopicStatusBadge({ status }: TopicStatusBadgeProps) {
  const { t } = useTranslation();
  const { variant, label } = topicStatusOption[status];
  return <Badge variant={variant}>{t(label)}</Badge>;
}
