import { Badge } from '@/components/ui/badge';
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table';
import type { Topic } from '@/lib/api/type/topic';
import { formatDate } from '@/formatter/date-formatter';
import TopicStatusBadge from './TopicStatusBadge';

interface TopicDataProps {
  dataList: Topic[];
}

export default function TopicData({ dataList }: TopicDataProps) {
  return (
    <div className="rounded-md border h-full">
      <Table className="[&>*]:whitespace-nowrap sticky top-0 bg-background after:content-[''] after:inset-x-0 after:h-px after:bg-border after:absolute after:bottom-0">
        <TableHeader>
          <TableRow>
            <TableHead className="w-[100px]">Code</TableHead>
            <TableHead>Message</TableHead>
            <TableHead className="w-[100px]">Total message</TableHead>
            <TableHead className="w-[100px]">Create date</TableHead>
            <TableHead className="w-[100px]">Status</TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          {dataList.map((data, idx) => (
            <TableRow key={idx} className="odd:bg-muted/50 [&>*]:whitespace-nowrap">
              <TableCell className="font-medium">{data.code}</TableCell>
              <TableCell>
                {data.description}{' '}
                {data.countOfMessageGroup > 0 && <Badge variant="secondary">+{data.countOfMessageGroup}</Badge>}
              </TableCell>
              <TableCell className="text-right">{data.countOfTotalMessage}</TableCell>
              <TableCell>{formatDate(data.createDate)}</TableCell>
              <TableCell>
                <TopicStatusBadge status={data.status}></TopicStatusBadge>
              </TableCell>
            </TableRow>
          ))}
        </TableBody>
      </Table>
    </div>
  );
}
