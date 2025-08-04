import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table';
import { formatDate } from '@/formatter/date-formatter';
import { MessageGroupStatus, type MessageGroup } from '@/lib/api/type/message-group';
import { useTranslation } from 'react-i18next';
import LoadingState from '@/components/state/LoadingState';
import ErrorState from '@/components/state/ErrorState';
import TableStateRow from '@/components/table/TableStateRow';
import NoDataState from '@/components/state/NoDataState';
import MessageGroupStatusBadge from './MessageGroupBadge';
import { TYLink } from '@/components/Typography';
import { Link } from 'react-router';

interface MessageGroupDataProps {
  isLoading: boolean;
  dataList?: MessageGroup[];
  error: Error | null;
  onStatusUpdate?: (id: string, status: MessageGroupStatus) => void;
}

const colSpan = 5;

export default function MessageGroupData({ isLoading, dataList = [], error }: MessageGroupDataProps) {
  const { t } = useTranslation();
  return (
    <Table containerClassName="table-round h-full">
      <TableHeader>
        <TableRow>
          <TableHead className="w-[100px]">{t('messageGroup.table.id')}</TableHead>
          <TableHead>{t('messageGroup.table.text')}</TableHead>
          <TableHead className="w-[100px]">{t('messageGroup.table.status')}</TableHead>
          <TableHead className="w-[100px]">{t('messageGroup.table.createdAt')}</TableHead>
          <TableHead className="w-[20px]">{t('common.actions')}</TableHead>
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
        ) : !dataList || dataList.length === 0 ? (
          <TableStateRow colSpan={colSpan}>
            <NoDataState />
          </TableStateRow>
        ) : (
          dataList.map((message) => (
            <TableRow key={message.id}>
              <TableCell>
                <Link to={`/message-group/${message.id}`}>
                  <TYLink>{message.id}</TYLink>
                </Link>
              </TableCell>
              <TableCell>{message.text}</TableCell>
              <TableCell>
                <MessageGroupStatusBadge status={message.status} />
              </TableCell>
              <TableCell>{formatDate(message.created_at)}</TableCell>
              <TableCell></TableCell>
            </TableRow>
          ))
        )}
      </TableBody>
    </Table>
  );
}
