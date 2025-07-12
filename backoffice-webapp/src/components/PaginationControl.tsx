import { Pagination, PaginationContent, PaginationItem, PaginationLink } from '@/components/ui/pagination';
import type { PaginationReq, PaginationRes } from '@/lib/api/type/base';

interface PaginationControlProps {
  paginationRes?: PaginationRes<any>;
  onPageChange: (paginationReq: PaginationReq) => void;
}

export default function PaginationControl({ paginationRes, onPageChange }: PaginationControlProps) {
  if (!paginationRes) {
    return <></>;
  }

  const { page: currentPage, pageSize, totalPages } = paginationRes;

  function handlePageChange(page: number) {
    onPageChange({
      page: page,
      pageSize: pageSize,
    });
  }

  // TODO implement <,> button and ...
  return (
    <Pagination>
      <PaginationContent>
        {Array.from({ length: totalPages }, (_, i) => {
          const page = i + 1;
          const isCurrentPage = page === currentPage;

          return (
            <PaginationItem key={page}>
              <PaginationLink isActive={page === currentPage} onClick={() => !isCurrentPage && handlePageChange(page)}>
                {page}
              </PaginationLink>
            </PaginationItem>
          );
        })}
      </PaginationContent>
    </Pagination>
  );
}
