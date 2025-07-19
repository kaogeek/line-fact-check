import { TYH3, TYMuted } from '@/components/Typography';
import { Navigate, useParams } from 'react-router';
import TopicStatusBadge from '../components/TopicStatusBadge';
import { useGetTopicById } from '@/hooks/api/topic';
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
import AddMessageDialog from './dialog/AddMessageDialog';
import { useLoader } from '@/hooks/loader';
import { createMessage } from '@/lib/api/service/message';
import { useMutation } from '@tanstack/react-query';
import { toast } from 'sonner';
import { updateAnswer } from '@/lib/api/service/topic-answer';
import { approveTopic, rejectTopic } from '@/lib/api/service/topic';
import { ConfirmAlertDialog } from '../../../components/ConfirmAlertDialog';

export default function TopicDetailPage() {
  const [selectedMessageId, setSelectedMessageId] = useState<string | null>(null);
  const [openAddMessageDialog, setOpenAddMessageDialog] = useState<boolean>(false);
  const [openTopicPickerDialog, setOpenTopicPickerDialog] = useState<boolean>(false);
  const [openTopicHistoryDialog, setOpenTopicHistoryDialog] = useState<boolean>(false);
  const [openAnswerHistoryDialog, setOpenAnswerHistoryDialog] = useState<boolean>(false);
  const [showApproveDialog, setShowApproveDialog] = useState(false);
  const [showRejectDialog, setShowRejectDialog] = useState(false);
  const { id } = useParams();
  const { startLoading, stopLoading } = useLoader();

  if (!id) {
    return <Navigate to="/404" replace />;
  }

  const { isLoading, data: topic, error } = useGetTopicById(id);

  const { mutate: createMessageMutation } = useMutation({
    mutationFn: (message: string) => createMessage(id!, message),
    onSettled: () => {
      stopLoading();
    },
    onSuccess: () => {
      toast.success('Message has been created.');
    },
    onError: (err) => {
      toast.error('Failed to create message.');
      console.error(err);
    },
  });

  const { mutate: updateAnswerMutation } = useMutation({
    mutationFn: ({ answerId, content }: { answerId: string; content: string }) => updateAnswer(id!, answerId, content),
    onSettled: () => {
      stopLoading();
    },
    onSuccess: () => {
      toast.success('Answer updated successfully');
    },
    onError: (err) => {
      toast.error('Failed to update answer');
      console.error(err);
    },
  });

  const { mutate: approveTopicMutation } = useMutation({
    mutationFn: () => approveTopic(id!),
    onSettled: () => {
      stopLoading();
    },
    onSuccess: () => {
      toast.success('Topic approved successfully');
    },
    onError: (err) => {
      toast.error('Failed to approve topic');
      console.error(err);
    },
  });

  const { mutate: rejectTopicMutation } = useMutation({
    mutationFn: () => rejectTopic(id!),
    onSettled: () => {
      stopLoading();
    },
    onSuccess: () => {
      toast.success('Topic rejected successfully');
    },
    onError: (err) => {
      toast.error('Failed to reject topic');
      console.error(err);
    },
  });

  function handleClickMoveMessage(messageId: string) {
    setSelectedMessageId(messageId);
    setOpenTopicPickerDialog(true);
  }

  function handleOpenAddMessageDialog() {
    setOpenAddMessageDialog(true);
  }

  async function handleCreateMessage(message: string) {
    startLoading();
    createMessageMutation(message);
  }

  async function handleUpdateAnswer(answerId: string, content: string) {
    startLoading();
    updateAnswerMutation({ answerId, content });
  }

  function handleChooseDestination(topicId: string) {
    console.log(`Move ${selectedMessageId} to ${topicId}`);
  }

  function handleClickAnswerHistory() {
    setOpenAnswerHistoryDialog(true);
  }

  function handleClickTopicHistory() {
    setOpenTopicHistoryDialog(true);
  }

  const handleApproveClick = () => {
    setShowApproveDialog(true);
  };

  const handleRejectClick = () => {
    setShowRejectDialog(true);
  };

  const handleConfirmApprove = () => {
    setShowApproveDialog(false);
    startLoading();
    approveTopicMutation();
  };

  const handleConfirmReject = () => {
    setShowRejectDialog(false);
    startLoading();
    rejectTopicMutation();
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
          <AddMessageDialog
            open={openAddMessageDialog}
            onOpenChange={setOpenAddMessageDialog}
            onSubmit={handleCreateMessage}
          />
          <TopicPickerDialog
            open={openTopicPickerDialog}
            onOpenChange={setOpenTopicPickerDialog}
            currentId={id}
            onChoose={handleChooseDestination}
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
          <ConfirmAlertDialog
            open={showApproveDialog}
            onOpenChange={setShowApproveDialog}
            title="Approve Topic"
            description="Are you sure you want to approve this topic?"
            confirmText="Approve"
            onConfirm={handleConfirmApprove}
          />
          <ConfirmAlertDialog
            open={showRejectDialog}
            onOpenChange={setShowRejectDialog}
            title="Reject Topic"
            description="Are you sure you want to reject this topic?"
            confirmText="Reject"
            onConfirm={handleConfirmReject}
          />
          <div className="flex flex-col gap-4 p-4">
            <div className="flex flex-col">
              <div className="flex gap-2">
                <TYH3 className="flex-1">Topic: {topic.code}</TYH3>
                <TopicStatusBadge status={topic.status} />
                <Button variant="outline" onClick={handleClickTopicHistory}>
                  History
                </Button>
                <DropdownMenu>
                  <DropdownMenuTrigger asChild>
                    <Button variant="outline">
                      <EllipsisVertical />
                    </Button>
                  </DropdownMenuTrigger>
                  <DropdownMenuContent>
                    <DropdownMenuItem onClick={handleApproveClick}>Approve</DropdownMenuItem>
                    <DropdownMenuItem onClick={handleRejectClick}>Reject</DropdownMenuItem>
                  </DropdownMenuContent>
                </DropdownMenu>
              </div>
              <TYMuted>Create at: {formatDate(topic.createDate)}</TYMuted>
            </div>
            <TopicMessageDetail
              topicId={topic.id}
              onClickMove={handleClickMoveMessage}
              onClickCreate={handleOpenAddMessageDialog}
            />
            <TopicMessageAnswer
              topicId={topic.id}
              onClickHistory={handleClickAnswerHistory}
              onUpdateAnswer={handleUpdateAnswer}
            />
          </div>
        </>
      )}
    </>
  );
}
