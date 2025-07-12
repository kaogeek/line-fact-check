export interface PaginationReq {
  page?: number;
  pageSize?: number;
}

export interface PaginationRes<T> {
  items: T[];
  page: number;
  pageSize: number;
  totalItems: number;
  totalPages: number;
}
