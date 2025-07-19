import { useTranslation } from 'react-i18next';
import { TYH3 } from '@/components/Typography';
import { Button } from '@/components/ui/button';
import { MoveRight, Plus } from 'lucide-react';
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table';
import { formatDate } from '@/formatter/date-formatter';
import { useGetMessageByTopicId } from '@/hooks/api/message';
import TableStateRow from '@/components/table/TableStateRow';
import LoadingState from '@/components/state/LoadingState';
import ErrorState from '@/components/state/ErrorState';
import NoDataState from '@/components/state/NoDataState';

interface TopicMessageDetailProps {
  topicId: string | null;
  onClickMove: (messageId: string) => void;
  onClickCreate: () => void;
}

const colSpan = 4;

export default function TopicMessageDetail({ topicId, onClickMove, onClickCreate }: TopicMessageDetailProps) {
  const { t } = useTranslation();

  if (!topicId) {
    return <></>;
  }

  const { isLoading, data: dataList, error } = useGetMessageByTopicId(topicId);

  return (
    <div className="flex flex-col gap-2">
      <div className="flex gap-2">
        <TYH3 className="flex-1">{t('topicMessageDetail.title')}</TYH3>
        <Button variant="default" size="icon" onClick={onClickCreate}>
          <Plus />
        </Button>
      </div>
      <Table containerClassName="table-round h-full">
        <TableHeader>
          <TableRow>
            <TableHead className="w-[100px]">{t('topicMessageDetail.tableHeaders.code')}</TableHead>
            <TableHead>{t('topicMessageDetail.tableHeaders.message')}</TableHead>
            <TableHead className="w-[100px]">{t('topicMessageDetail.tableHeaders.totalMessage')}</TableHead>
            <TableHead className="w-[100px]">{t('topicMessageDetail.tableHeaders.createDate')}</TableHead>
            <TableHead className="w-[20px]"></TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          {isLoading ? (
            <TableStateRow colSpan={colSpan}>
              <LoadingState />
            </TableStateRow>
          ) : error ? (
            <TableStateRow colSpan={colSpan}>
              <ErrorState />
            </TableStateRow>
          ) : !dataList || !dataList.length ? (
            <TableStateRow colSpan={colSpan}>
              <NoDataState />
            </TableStateRow>
          ) : (
            dataList.map((data, idx) => (
              <TableRow key={idx}>
                <TableCell>{data.code}</TableCell>
                <TableCell>{data.message}</TableCell>
                <TableCell className="text-right">{data.countOfMessageGroup}</TableCell>
                <TableCell>{formatDate(data.createDate)}</TableCell>
                <TableCell>
                  <Button
                    variant="outline"
                    size="icon"
                    onClick={() => onClickMove(data.id)}
                    aria-label={t('topicMessageDetail.moveButton')}
                  >
                    <MoveRight />
                  </Button>
                </TableCell>
              </TableRow>
            ))
          )}
        </TableBody>
      </Table>
    </div>
  );
}
