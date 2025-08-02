import type { StrictPaginationReq } from '@/lib/api/type/base';
import { useState } from 'react';

export function usePaginationReqState(initial: StrictPaginationReq) {
  const [paginationReq, setPaginationReq] = useState(initial);

  function resetPage() {
    setPaginationReq((prev) => ({ ...prev, page: 1 }));
  }

  function setPage(page: number) {
    setPaginationReq((prev) => ({ ...prev, page }));
  }

  function setPageSize(pageSize: number) {
    setPaginationReq((prev) => ({ ...prev, pageSize }));
  }

  return { paginationReq, setPage, setPageSize, resetPage };
}
