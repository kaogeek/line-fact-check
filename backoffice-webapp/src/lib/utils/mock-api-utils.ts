import { MOCKUP_API_LOADING_MS } from '@/constants/app';

export function generateTraceId(): string {
  const chars = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789!@#$%^&*';
  let result = '';
  for (let i = 0; i < 6; i++) {
    result += chars.charAt(Math.floor(Math.random() * chars.length));
  }
  return result;
}

export function mockApi<T>(callback: () => T | Promise<T>, delay = MOCKUP_API_LOADING_MS): Promise<T> {
  const traceId = generateTraceId();
  console.debug(`[Mock API] [${traceId}] Request started...`);
  return new Promise((resolve, reject) => {
    setTimeout(() => {
      try {
        const result = callback();
        if (result instanceof Promise) {
          result
            .then((data) => {
              console.debug(`[Mock API] [${traceId}] Response resolved:`, data);
              resolve(data);
            })
            .catch((err) => {
              console.error(`[Mock API] [${traceId}] Response error:`, err);
              reject(err);
            });
        } else {
          console.debug(`[Mock API] [${traceId}] Response resolved:`, result);
          resolve(result);
        }
      } catch (error) {
        console.error(`[Mock API] [${traceId}] Callback threw error:`, error);
        reject(error);
      }
    }, delay);
  });
}
