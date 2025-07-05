import { Badge } from '@/components/ui/badge';
import {
  Table,
  TableBody,
  TableCaption,
  TableCell,
  TableFooter,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table';
import type { Topic } from '@/constants/topic';
import { formatDate } from '@/formatter/date-formatter';

interface TopicDataProps {
  dataList: Topic[];
}

export default function TopicData({ dataList }: TopicDataProps) {
  return (
    <Table>
      <TableHeader>
        <TableRow>
          <TableHead className="w-[100px]">Code</TableHead>
          <TableHead className="w-[100px]">Status</TableHead>
          <TableHead className="w-[100px]">Create date</TableHead>
          <TableHead>Message</TableHead>
          <TableHead className="w-[100px]">Total message</TableHead>
        </TableRow>
      </TableHeader>
      <TableBody>
        {dataList.map((data, idx) => (
          <TableRow key={idx}>
            <TableCell className="font-medium">{data.code}</TableCell>
            <TableCell>
              <Badge variant="secondary">{data.status}</Badge>
            </TableCell>
            <TableCell>{formatDate(data.createDate)}</TableCell>
            <TableCell>
              {data.description}{' '}
              {data.countOfMessageGroup > 0 && <Badge variant="secondary">{data.countOfMessageGroup}</Badge>}
            </TableCell>
            <TableCell className="text-right">{data.countOfTotalMessage}</TableCell>
          </TableRow>
        ))}
      </TableBody>
    </Table>
  );
}
