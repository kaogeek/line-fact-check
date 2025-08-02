import {
  Pagination,
  PaginationContent,
  PaginationItem,
  PaginationLink,
  PaginationPrevious,
  PaginationNext,
  PaginationEllipsis,
} from '@/components/ui/pagination';
import type { StrictPaginationReq } from '@/lib/api/type/base';

interface PaginationControlProps {
  paginationReq: StrictPaginationReq;
  // TODO use interface instead
  pagination: {
    totalPages: number;
  };
  onPageChange: (page: number) => void;
}

export default function PaginationControl({ paginationReq, pagination, onPageChange }: PaginationControlProps) {
  const { page: currentPage } = paginationReq;
  const { totalPages } = pagination;

  const maxVisiblePages = 5;
  const shouldCollapse = totalPages > maxVisiblePages;

  function handlePageChange(page: number) {
    onPageChange(page);
  }

  const renderPageItems = () => {
    if (!shouldCollapse) {
      return Array.from({ length: totalPages }, (_, i) => renderPageItem(i + 1));
    }

    const pages = [];
    // Always show first page
    pages.push(renderPageItem(1));

    // Show ellipsis if current page is beyond first 3 pages
    if (currentPage > 3) {
      pages.push(
        <PaginationItem key="ellipsis-start">
          <PaginationEllipsis />
        </PaginationItem>
      );
    }

    // Show current page and neighbors
    const start = Math.max(2, currentPage - 1);
    const end = Math.min(totalPages - 1, currentPage + 1);

    for (let i = start; i <= end; i++) {
      pages.push(renderPageItem(i));
    }

    // Show ellipsis if current page is not in last 3 pages
    if (currentPage < totalPages - 2) {
      pages.push(
        <PaginationItem key="ellipsis-end">
          <PaginationEllipsis />
        </PaginationItem>
      );
    }

    // Always show last page
    if (totalPages > 1) {
      pages.push(renderPageItem(totalPages));
    }

    return pages;
  };

  const renderPageItem = (page: number) => (
    <PaginationItem key={page}>
      <PaginationLink isActive={page === currentPage} onClick={() => !(page === currentPage) && handlePageChange(page)}>
        {page}
      </PaginationLink>
    </PaginationItem>
  );

  return (
    <Pagination>
      <PaginationContent>
        <PaginationItem>
          <PaginationPrevious
            onClick={() => (currentPage > 1 ? handlePageChange(currentPage - 1) : null)}
            className={currentPage === 1 ? 'opacity-50 cursor-not-allowed' : ''}
          />
        </PaginationItem>
        {renderPageItems()}
        <PaginationItem>
          <PaginationNext
            onClick={() => (currentPage < totalPages ? handlePageChange(currentPage + 1) : null)}
            className={currentPage === totalPages ? 'opacity-50 cursor-not-allowed' : ''}
          />
        </PaginationItem>
      </PaginationContent>
    </Pagination>
  );
}
