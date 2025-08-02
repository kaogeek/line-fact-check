import { describe, it, expect } from 'vitest';
import { appendPaginationParams } from './base';

describe('appendPaginationParams', () => {
  it('should append default pagination params when no values provided', () => {
    const params = new URLSearchParams();
    appendPaginationParams(params, {});

    expect(params.get('offset')).toBe('0'); // (1 - 1) * 10 = 0
    expect(params.get('limit')).toBe('10');
  });

  it('should append correct pagination params when page and pageSize are provided', () => {
    const params = new URLSearchParams();
    appendPaginationParams(params, { page: 2, pageSize: 20 });

    expect(params.get('offset')).toBe('20'); // (2 - 1) * 20 = 20
    expect(params.get('limit')).toBe('20');
  });

  it('should handle only page being provided', () => {
    const params = new URLSearchParams();
    appendPaginationParams(params, { page: 3 });

    expect(params.get('offset')).toBe('20'); // (3 - 1) * 10 = 20
    expect(params.get('limit')).toBe('10');
  });

  it('should handle only pageSize being provided', () => {
    const params = new URLSearchParams();
    appendPaginationParams(params, { pageSize: 5 });

    expect(params.get('offset')).toBe('0'); // (1 - 1) * 5 = 0
    expect(params.get('limit')).toBe('5');
  });

  it('should append to existing params without overriding them', () => {
    const params = new URLSearchParams('foo=bar');
    appendPaginationParams(params, { page: 2, pageSize: 15 });

    expect(params.get('foo')).toBe('bar');
    expect(params.get('offset')).toBe('15'); // (2 - 1) * 15 = 15
    expect(params.get('limit')).toBe('15');
  });
});
