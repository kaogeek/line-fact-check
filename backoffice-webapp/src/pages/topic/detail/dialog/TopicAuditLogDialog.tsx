import { Dialog, DialogContent, DialogDescription, DialogHeader, DialogTitle } from '@/components/ui/dialog';
import { useGetTopicAuditLogs } from '@/hooks/api/useTopicAuditLog';
import TopicAuditLogCard from '../components/TopicAuditLogCard';

interface TopicAuditLogDialog {
  open?: boolean;
  onOpenChange?(open: boolean): void;
  topicId: string;
}

export default function TopicAuditLogDialog({ open, onOpenChange, topicId }: TopicAuditLogDialog) {
  const topicAuditLogs = useGetTopicAuditLogs(topicId);

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Topic history</DialogTitle>
          <DialogDescription>
            <div className="flex flex-col gap-2">
              {topicAuditLogs.map((log, idx) => (
                <TopicAuditLogCard
                  key={idx}
                  avatarUrl={log.avatarUrl}
                  username={log.username}
                  actionDate={log.actionDate}
                  status={log.status}
                  detail={log.detail}
                />
              ))}
            </div>
          </DialogDescription>
        </DialogHeader>
      </DialogContent>
    </Dialog>
  );
}
