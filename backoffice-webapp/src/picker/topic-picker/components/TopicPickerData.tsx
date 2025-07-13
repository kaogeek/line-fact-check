import { Badge } from '@/components/ui/badge';
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table';
import type { Topic } from '@/lib/api/type/topic';
import { formatDate } from '@/formatter/date-formatter';
import { Link } from 'react-router';
import { TYLink } from '@/components/Typography';
import LoadingState from '@/components/state/LoadingState';
import ErrorState from '@/components/state/ErrorState';
import TableStateRow from '@/components/table/TableStateRow';
import NoDataState from '@/components/state/NoDataState';
import TopicStatusBadge from '@/pages/topic/components/TopicStatusBadge';
import { Button } from '@/components/ui/button';
import { cn } from '@/lib/utils';

interface TopicPickerDataProps {
  className?: string;
  isLoading: boolean;
  dataList?: Topic[];
  error: Error | null;
  onChoose: (topicId: string) => void;
}

const colSpan = 6;

export default function TopicPickerData({ className, isLoading, dataList, error, onChoose }: TopicPickerDataProps) {
  return (
    <Table containerClassName={cn('table-round', className)}>
      <TableHeader className="sticky-header">
        <TableRow>
          <TableHead className="w-[100px]">Code</TableHead>
          <TableHead>Message</TableHead>
          <TableHead className="w-[100px]">Total message</TableHead>
          <TableHead className="w-[100px]">Create date</TableHead>
          <TableHead className="w-[100px]">Status</TableHead>
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
        ) : !dataList ? (
          <TableStateRow colSpan={colSpan}>
            <NoDataState />
          </TableStateRow>
        ) : (
          dataList.map((data, idx) => (
            <TableRow key={idx}>
              <TableCell>
                <Link to={`/topic/${data.id}`}>
                  <TYLink>{data.code}</TYLink>
                </Link>
              </TableCell>
              <TableCell>
                {data.description}{' '}
                {data.countOfMessageGroup > 0 && <Badge variant="secondary">+{data.countOfMessageGroup}</Badge>}
              </TableCell>
              <TableCell className="text-right">{data.countOfTotalMessage}</TableCell>
              <TableCell>{formatDate(data.createDate)}</TableCell>
              <TableCell>
                <TopicStatusBadge status={data.status}></TopicStatusBadge>
              </TableCell>
              <TableCell>
                <Button variant="default" onClick={() => onChoose(data.id)}>
                  Choose
                </Button>
              </TableCell>
            </TableRow>
          ))
        )}
      </TableBody>
    </Table>
  );
}
