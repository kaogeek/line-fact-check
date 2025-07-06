import { TYH3 } from '@/components/Typography';
import { Button } from '@/components/ui/button';
import { Plus } from 'lucide-react';
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table';
import type { Message } from '@/lib/api/type/message';
import { formatDate } from '@/formatter/date-formatter';

interface TopicMessageDetailProps {
  dataList: Message[];
}

export default function TopicMessageDetail({ dataList }: TopicMessageDetailProps) {
  return (
    <div className="flex flex-col gap-2">
      <div className="flex gap-2">
        <TYH3 className="flex-1">Message</TYH3>
        <Button variant="default">
          <Plus />
        </Button>
      </div>
      <div className="rounded-md border h-full">
        <Table className="[&>*]:whitespace-nowrap sticky top-0 bg-background after:content-[''] after:inset-x-0 after:h-px after:bg-border after:absolute after:bottom-0">
          <TableHeader>
            <TableRow>
              <TableHead className="w-[100px]">Code</TableHead>
              <TableHead>Message</TableHead>
              <TableHead className="w-[100px]">Total message</TableHead>
              <TableHead className="w-[100px]">Create date</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {dataList.map((data, idx) => (
              <TableRow key={idx} className="odd:bg-muted/50 [&>*]:whitespace-nowrap">
                <TableCell>{data.code}</TableCell>
                <TableCell>{data.message}</TableCell>
                <TableCell className="text-right">{data.countOfMessageGroup}</TableCell>
                <TableCell>{formatDate(data.createDate)}</TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
      </div>
    </div>
  );
}
