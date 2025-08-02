import type { Pagination, PaginationReq, PaginationRes, StrictPaginationReq } from '../api/type/base';

export function paginate<T>(items: T[], paginationReq: StrictPaginationReq): PaginationRes<T> {
  const { page, pageSize } = paginationReq;
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

export function calPagination(totalItems: number, paginationReq: StrictPaginationReq): Pagination {
  const { pageSize } = paginationReq;
  const totalPages = Math.ceil(totalItems / pageSize);

  return {
    totalItems,
    totalPages,
  };
}

export function convertToOffsetLimit(pagination: PaginationReq): { offset: number; limit: number } {
  const { page = 1, pageSize = 10 } = pagination;

  return {
    offset: (page - 1) * pageSize,
    limit: pageSize,
  };
}
