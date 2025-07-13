import type { PaginationReq, PaginationRes } from '../api/type/base';

export function paginate<T>(items: T[], paginationReq: PaginationReq): PaginationRes<T> {
  const { page = 1, pageSize = 10 } = paginationReq;
  const totalItems = items.length;
  const totalPages = Math.ceil(totalItems / pageSize);

  const startIndex = (page - 1) * pageSize;
  const endIndex = startIndex + pageSize;

  const pagedItems = items.slice(startIndex, endIndex);

  return {
    items: pagedItems,
    page,
    pageSize,
    totalItems,
    totalPages,
  };
}
