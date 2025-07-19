import { Badge } from '@/components/ui/badge';
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table';
import type { Topic } from '@/lib/api/type/topic';
import { formatDate } from '@/formatter/date-formatter';
import TopicStatusBadge from './TopicStatusBadge';
import { Link } from 'react-router';
import { TYLink } from '@/components/Typography';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu';
import { EllipsisVertical } from 'lucide-react';
import LoadingState from '@/components/state/LoadingState';
import ErrorState from '@/components/state/ErrorState';
import TableStateRow from '@/components/table/TableStateRow';
import NoDataState from '@/components/state/NoDataState';
import { Button } from '@/components/ui/button';
import { useTranslation } from 'react-i18next';

interface TopicDataProps {
  isLoading: boolean;
  dataList?: Topic[];
  error: Error | null;
  onReject?: (topic: Topic, idx: number) => void;
}

const colSpan = 6;

export default function TopicData({ isLoading, dataList, error, onReject }: TopicDataProps) {
  const { t } = useTranslation();

  return (
    <Table containerClassName="table-round h-full">
      <TableHeader>
        <TableRow>
          <TableHead className="w-[100px]">{t('topic.searchLabel.code')}</TableHead>
          <TableHead>{t('topic.searchLabel.message')}</TableHead>
          <TableHead className="w-[100px]">{t('topic.totalMessages')}</TableHead>
          <TableHead className="w-[100px]">{t('topic.createDate')}</TableHead>
          <TableHead className="w-[100px]">{t('topic.status.label')}</TableHead>
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
        ) : !dataList || dataList.length === 0 ? (
          <TableStateRow colSpan={colSpan}>
            <NoDataState />
          </TableStateRow>
        ) : (
          dataList.map((topic, idx) => (
            <TableRow key={topic.id}>
              <TableCell>
                <Link to={`/topic/${topic.id}`}>
                  <TYLink>{topic.code}</TYLink>
                </Link>
              </TableCell>
              <TableCell>
                {topic.description}{' '}
                {topic.countOfMessageGroup > 0 && <Badge variant="secondary">+{topic.countOfMessageGroup}</Badge>}
              </TableCell>
              <TableCell className="text-right">{topic.countOfTotalMessage}</TableCell>
              <TableCell>{formatDate(topic.createDate)}</TableCell>
              <TableCell>
                <TopicStatusBadge status={topic.status} />
              </TableCell>
              <TableCell>
                <DropdownMenu>
                  <DropdownMenuTrigger asChild>
                    <Button variant="ghost" className="h-8 w-8 p-0">
                      <EllipsisVertical className="h-4 w-4" />
                    </Button>
                  </DropdownMenuTrigger>
                  <DropdownMenuContent align="end">
                    <DropdownMenuItem asChild>
                      <Link to={`/topic/${topic.id}`}>{t('topic.viewDetail')}</Link>
                    </DropdownMenuItem>
                    {onReject && (
                      <DropdownMenuItem onClick={() => onReject(topic, idx)}>{t('topic.reject')}</DropdownMenuItem>
                    )}
                  </DropdownMenuContent>
                </DropdownMenu>
              </TableCell>
            </TableRow>
          ))
        )}
      </TableBody>
    </Table>
  );
}
