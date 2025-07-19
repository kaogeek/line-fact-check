import { TYMuted, TYP } from '@/components/Typography';
import { formatDate } from '@/formatter/date-formatter';
import { useTranslation } from 'react-i18next';

interface AuditLogCardProps {
  avatarUrl: string;
  username: string;
  actionDate: Date;
  actionDescription: string;
  actionDetail: string;
}

export default function AuditLogCard({
  avatarUrl,
  username,
  actionDate,
  actionDescription,
  actionDetail,
}: AuditLogCardProps) {
  const { t } = useTranslation();
  return (
    <div className="flex gap-2">
      <div className="flex flex-col items-center gap-2">
        <img src={avatarUrl} alt={`avatar`} className="w-[32px] h-[32px] object-cover rounded-full shadow-md" />
        <div className="flex-1 w-1 bg-muted rounded-md"></div>
      </div>
      <div className="flex-1 text-black">
        <div className="flex gap-2">
          <TYP className="font-bold">{username}</TYP> <p>{t(actionDescription)}</p>
        </div>

        <TYMuted>{formatDate(actionDate)}</TYMuted>
        <TYP>{actionDetail}</TYP>
      </div>
    </div>
  );
}
