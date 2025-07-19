import { Dialog, DialogContent, DialogDescription, DialogHeader, DialogTitle } from '@/components/ui/dialog';
import { useGetTopicAuditLogs } from '@/hooks/api/topicAuditLog';
import TopicAuditLogCard from '../components/TopicAuditLogCard';
import NoDataState from '@/components/state/NoDataState';
import ErrorState from '@/components/state/ErrorState';
import LoadingState from '@/components/state/LoadingState';

interface TopicAuditLogDialog {
  open?: boolean;
  onOpenChange?(open: boolean): void;
  topicId: string;
}

export default function TopicAuditLogDialog({ open, onOpenChange, topicId }: TopicAuditLogDialog) {
  const {
    isLoading,
    data: topicAuditLogs,
    error,
  } = useGetTopicAuditLogs(topicId, undefined, {
    enabled: open,
  });

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Topic history</DialogTitle>
          <DialogDescription asChild>
            <div className="flex flex-col gap-2">
              {isLoading ? (
                <LoadingState />
              ) : error ? (
                <ErrorState />
              ) : !topicAuditLogs || topicAuditLogs.length === 0 ? (
                <NoDataState />
              ) : (
                topicAuditLogs.map((log, idx) => (
                  <TopicAuditLogCard
                    key={idx}
                    avatarUrl={log.avatarUrl}
                    username={log.username}
                    actionDate={log.actionDate}
                    status={log.status}
                    detail={log.detail}
                  />
                ))
              )}
            </div>
          </DialogDescription>
        </DialogHeader>
      </DialogContent>
    </Dialog>
  );
}
