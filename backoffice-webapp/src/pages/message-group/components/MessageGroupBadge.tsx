import { Badge } from '@/components/ui/badge';
import { MessageGroupStatus, messageGroupStatusOption } from '@/lib/api/type/message-group';
import { useTranslation } from 'react-i18next';

interface MessageGroupStatusBadgeProps {
  status: MessageGroupStatus;
}

export default function MessageGroupStatusBadge({ status }: MessageGroupStatusBadgeProps) {
  const { t } = useTranslation();
  const { variant, label } = messageGroupStatusOption[status];
  return <Badge variant={variant}>{t(label)}</Badge>;
}
