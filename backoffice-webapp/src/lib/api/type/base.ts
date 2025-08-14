export interface PaginationReq {
  page?: number;
  pageSize?: number;
}

export type StrictPaginationReq = Required<PaginationReq>;

export interface Pagination {
  totalItems: number;
  totalPages: number;
}

export interface PaginationRes<T> extends Pagination {
  page: number;
  pageSize: number;
  items: T[];
}
