import { convertToOffsetLimit } from '@/lib/utils/page-utils';
import type { PaginationReq } from '../type/base';

export function appendPaginationParams(params: URLSearchParams, pagination: PaginationReq): void {
  const { offset, limit } = convertToOffsetLimit(pagination);
  params.append('offset', offset.toString());
  params.append('limit', limit.toString());
}
