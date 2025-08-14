import { vi, describe, beforeAll, afterEach, afterAll, it, expect } from 'vitest';
import { mockApi, generateTraceId } from './mock-api-utils';
import { MOCKUP_API_LOADING_MS } from '@/constants/app';

describe('mock-api-utils', () => {
  describe('generateTraceId', () => {
    it('should generate a 6-character string', () => {
      const traceId = generateTraceId();
      expect(traceId).toHaveLength(6);
    });

    it('should only contain valid characters', () => {
      const traceId = generateTraceId();
      expect(traceId).toMatch(/^[A-Za-z0-9!@#$%^&*]{6}$/);
    });
  });

  describe('mockApi', () => {
    const MOCK_DELAY = 1000;
    const CALLER_NAME = 'testCaller';

    // Mock console methods
    const originalConsole = { ...console };
    const mockDebug = vi.fn();
    const mockError = vi.fn();

    beforeAll(() => {
      global.console = {
        ...originalConsole,
        debug: mockDebug,
        error: mockError,
      };
      vi.useFakeTimers();
    });

    afterEach(() => {
      vi.clearAllMocks();
      vi.clearAllTimers();
    });

    afterAll(() => {
      global.console = originalConsole;
      vi.useRealTimers();
    });

    it('should resolve with the correct value for synchronous callback', async () => {
      const expectedValue = { data: 'test' };
      const callback = vi.fn().mockReturnValue(expectedValue);

      const promise = mockApi(callback, CALLER_NAME, MOCK_DELAY);

      await vi.advanceTimersByTimeAsync(MOCK_DELAY);

      const result = await promise;

      expect(callback).toHaveBeenCalled();
      expect(mockDebug).toHaveBeenCalledWith(expect.stringContaining(`[${CALLER_NAME}] Request started`));
      expect(mockDebug).toHaveBeenCalledWith(expect.stringContaining('Response resolved:'), expect.anything());
      expect(result).toEqual(expectedValue);
    });

    it('should resolve with the correct value for asynchronous callback', async () => {
      const expectedValue = { data: 'async test' };
      const callback = vi.fn().mockResolvedValue(expectedValue);

      const promise = mockApi(callback, CALLER_NAME, MOCK_DELAY);

      // Fast-forward time
      await vi.advanceTimersByTimeAsync(MOCK_DELAY);

      const result = await promise;

      expect(callback).toHaveBeenCalled();
      expect(mockDebug).toHaveBeenCalledWith(expect.stringContaining(`[${CALLER_NAME}] Request started`));
      expect(mockDebug).toHaveBeenCalledWith(expect.stringContaining('Response resolved:'), expect.anything());
      expect(result).toEqual(expectedValue);
    });

    it('should reject when synchronous callback throws', async () => {
      const error = new Error('Test error');
      const callback = vi.fn().mockImplementation(() => {
        throw error;
      });
      let errorThrown;

      const promise = mockApi(callback, CALLER_NAME, MOCK_DELAY).catch((e) => {
        errorThrown = e;
      });

      // Fast-forward time
      await vi.advanceTimersByTimeAsync(MOCK_DELAY);
      await promise;

      expect(errorThrown).not.toBeUndefined();
      expect(errorThrown).toBe(error);
      expect(mockError).toHaveBeenCalledWith(expect.stringContaining('Callback threw error:'), expect.anything());
    });

    it('should reject when asynchronous callback rejects', async () => {
      const error = new Error('Async test error');
      const callback = vi.fn().mockRejectedValue(error);
      let errorThrown;

      const promise = mockApi(callback, CALLER_NAME, MOCK_DELAY).catch((e) => {
        errorThrown = e;
      });

      // Fast-forward time
      await vi.advanceTimersByTimeAsync(MOCK_DELAY);
      await promise;

      expect(errorThrown).not.toBeUndefined();
      expect(errorThrown).toBe(error);
      expect(mockError).toHaveBeenCalledWith(expect.stringContaining('Response error:'), expect.anything());
    });

    it('should use default delay when not provided', async () => {
      const callback = vi.fn().mockReturnValue('test');

      mockApi(callback, CALLER_NAME);

      // Verify callback not called before delay
      expect(callback).not.toHaveBeenCalled();

      // Fast-forward time
      await vi.advanceTimersByTimeAsync(MOCKUP_API_LOADING_MS); // Default delay
      expect(callback).toHaveBeenCalled();
    });
  });
});
