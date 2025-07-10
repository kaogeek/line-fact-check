import AuditLogCard from '@/components/AuditLogCard';
import { topicAuditLogTypeOption, type TopicAuditLogType } from '@/lib/api/type/topic-audit-log';

interface AuditLogCardProps {
  avatarUrl: string;
  username: string;
  actionDate: Date;
  status: TopicAuditLogType;
  detail: string;
}

export default function TopicAuditLogCard({ avatarUrl, username, actionDate, status, detail }: AuditLogCardProps) {
  const { actionDescription } = topicAuditLogTypeOption[status];
  return (
    <AuditLogCard
      avatarUrl={avatarUrl}
      username={username}
      actionDate={actionDate}
      actionDescription={actionDescription}
      actionDetail={detail}
    />
  );
}
