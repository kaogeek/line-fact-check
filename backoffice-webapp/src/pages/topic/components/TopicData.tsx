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

interface TopicDataProps {
  dataList: Topic[];
  onReject?: (topic: Topic, idx: number) => void;
}

export default function TopicData({ dataList, onReject }: TopicDataProps) {
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
            <TableHead className="w-[20px]"></TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          {dataList.map((data, idx) => (
            <TableRow key={idx} className="odd:bg-muted/50 [&>*]:whitespace-nowrap">
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
                <DropdownMenu>
                  <DropdownMenuTrigger>
                    <EllipsisVertical />
                  </DropdownMenuTrigger>
                  <DropdownMenuContent>
                    <DropdownMenuItem onSelect={() => onReject && onReject(data, idx)}>Reject</DropdownMenuItem>
                  </DropdownMenuContent>
                </DropdownMenu>
              </TableCell>
            </TableRow>
          ))}
        </TableBody>
      </Table>
    </div>
  );
}
