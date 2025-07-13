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
import TopicMessageAnswer from './components/TopicMessageAnswer';
import { useState } from 'react';
import AnswerAuditLogDialog from './dialog/AnswerAuditLogDialog';
import TopicAuditLogDialog from './dialog/TopicAuditLogDialog';
import LoadingState from '@/components/state/LoadingState';
import ErrorState from '@/components/state/ErrorState';
import TopicPickerDialog from '@/picker/topic-picker/TopicPickerDialog';

export default function TopicDetailPage() {
  const [selectedMessageId, setSelectedMessageId] = useState<string | null>(null);
  const [openTopicPickerDialog, setOpenTopicPickerDialog] = useState<boolean>(false);
  const [openTopicHistoryDialog, setOpenTopicHistoryDialog] = useState<boolean>(false);
  const [openAnswerHistoryDialog, setOpenAnswerHistoryDialog] = useState<boolean>(false);
  const { id } = useParams();

  if (!id) {
    return <Navigate to="/404" replace />;
  }

  const { isLoading, data: topic, error } = useGetTopicById(id);

  const onHandleClickMoveMessage = (messageId: string) => {
    setSelectedMessageId(messageId);
    setOpenTopicPickerDialog(true);
  };

  const onHandleChooseDestination = (topicId: string) => {
    console.log(`Move ${selectedMessageId} to ${topicId}`);
  };

  const onHandleClickAnswerHistory = () => {
    setOpenAnswerHistoryDialog(true);
  };

  const onHandleClickTopicHistory = () => {
    setOpenTopicHistoryDialog(true);
  };

  return (
    <>
      {isLoading ? (
        <LoadingState />
      ) : error ? (
        <ErrorState />
      ) : !topic ? (
        <Navigate to="/404" replace />
      ) : (
        <>
          <TopicPickerDialog
            open={openTopicPickerDialog}
            onOpenChange={setOpenTopicPickerDialog}
            currentId={id}
            onChoose={onHandleChooseDestination}
          ></TopicPickerDialog>
          <AnswerAuditLogDialog
            open={openAnswerHistoryDialog}
            onOpenChange={setOpenAnswerHistoryDialog}
            topicId={id}
          ></AnswerAuditLogDialog>
          <TopicAuditLogDialog
            open={openTopicHistoryDialog}
            onOpenChange={setOpenTopicHistoryDialog}
            topicId={id}
          ></TopicAuditLogDialog>
          <div className="flex flex-col gap-4 p-4 h-full overflow-auto">
            <div className="flex flex-col">
              <div className="flex gap-2">
                <TYH3 className="flex-1">Topic: {topic.code}</TYH3>
                <TopicStatusBadge status={topic.status} />
                <Button variant="outline" onClick={onHandleClickTopicHistory}>
                  History
                </Button>
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
            <TopicMessageDetail topicId={topic.id} onClickMove={onHandleClickMoveMessage} />
            <TopicMessageAnswer topicId={topic.id} onClickHistory={onHandleClickAnswerHistory} />
          </div>
        </>
      )}
    </>
  );
}
