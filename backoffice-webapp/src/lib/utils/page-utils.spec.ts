import { describe, it, expect } from 'vitest';
import { paginate, calPagination, convertToOffsetLimit } from './page-utils';
import type { StrictPaginationReq, PaginationReq } from '../api/type/base';

describe('page-utils', () => {
  describe('paginate', () => {
    const testItems = Array.from({ length: 50 }, (_, i) => ({ id: i + 1 }));

    it('should return first page of items with correct pagination info', () => {
      const paginationReq: StrictPaginationReq = { page: 1, pageSize: 10 };
      const result = paginate(testItems, paginationReq);

      expect(result.items).toHaveLength(10);
      expect(result.items[0].id).toBe(1);
      expect(result.items[9].id).toBe(10);
      expect(result.page).toBe(1);
      expect(result.pageSize).toBe(10);
      expect(result.totalItems).toBe(50);
      expect(result.totalPages).toBe(5);
    });

    it('should return second page of items', () => {
      const paginationReq: StrictPaginationReq = { page: 2, pageSize: 10 };
      const result = paginate(testItems, paginationReq);

      expect(result.items).toHaveLength(10);
      expect(result.items[0].id).toBe(11);
      expect(result.items[9].id).toBe(20);
      expect(result.page).toBe(2);
    });

    it('should handle last page with partial items', () => {
      const paginationReq: StrictPaginationReq = { page: 5, pageSize: 10 };
      const result = paginate(testItems, paginationReq);

      expect(result.items).toHaveLength(10);
      expect(result.items[0].id).toBe(41);
      expect(result.items[9].id).toBe(50);
      expect(result.page).toBe(5);
    });

    it('should handle empty array', () => {
      const paginationReq: StrictPaginationReq = { page: 1, pageSize: 10 };
      const result = paginate([], paginationReq);

      expect(result.items).toHaveLength(0);
      expect(result.totalItems).toBe(0);
      expect(result.totalPages).toBe(0);
    });
  });

  describe('calPagination', () => {
    it('should calculate pagination info correctly', () => {
      const paginationReq: StrictPaginationReq = { page: 2, pageSize: 15 };
      const totalItems = 50;
      const result = calPagination(totalItems, paginationReq);

      expect(result.totalItems).toBe(50);
      expect(result.totalPages).toBe(4); // 50 / 15 = 3.333... -> 4 pages
    });

    it('should handle totalItems of 0', () => {
      const paginationReq: StrictPaginationReq = { page: 1, pageSize: 10 };
      const result = calPagination(0, paginationReq);

      expect(result.totalItems).toBe(0);
      expect(result.totalPages).toBe(0);
    });
  });

  describe('convertToOffsetLimit', () => {
    it('should convert page and pageSize to offset and limit', () => {
      const pagination: PaginationReq = { page: 3, pageSize: 20 };
      const result = convertToOffsetLimit(pagination);

      expect(result.offset).toBe(40); // (3-1) * 20 = 40
      expect(result.limit).toBe(20);
    });

    it('should use default values when not provided', () => {
      const pagination: PaginationReq = {};
      const result = convertToOffsetLimit(pagination);

      expect(result.offset).toBe(0); // (1-1) * 10 = 0
      expect(result.limit).toBe(10);
    });

    it('should handle first page explicitly', () => {
      const pagination: PaginationReq = { page: 1, pageSize: 10 };
      const result = convertToOffsetLimit(pagination);

      expect(result.offset).toBe(0);
      expect(result.limit).toBe(10);
    });
  });
});
