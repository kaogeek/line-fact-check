import { Dialog, DialogContent, DialogDescription, DialogHeader, DialogTitle } from '@/components/ui/dialog';
import { useGetTopicAuditLogs } from '@/hooks/api/topicAuditLog';
import TopicAuditLogCard from '../components/TopicAuditLogCard';
import { TopicAuditLogType } from '@/lib/api/type/topic-audit-log';
import LoadingState from '@/components/state/LoadingState';
import ErrorState from '@/components/state/ErrorState';
import NoDataState from '@/components/state/NoDataState';

interface AnswerAuditLogDialogProps {
  open?: boolean;
  onOpenChange?(open: boolean): void;
  topicId: string;
}

export default function AnswerAuditLogDialog({ open, onOpenChange, topicId }: AnswerAuditLogDialogProps) {
  const {
    isLoading,
    data: topicAuditLogs,
    error,
  } = useGetTopicAuditLogs(topicId, [TopicAuditLogType.UPDATE_ANSWER], {
    enabled: open,
  });

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Answer history</DialogTitle>
          <DialogDescription asChild>
            <div className="flex flex-col gap-2">
              {isLoading ? (
                <LoadingState />
              ) : error ? (
                <ErrorState />
              ) : !topicAuditLogs ? (
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
