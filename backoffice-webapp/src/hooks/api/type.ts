import type { UseQueryOptions } from '@tanstack/react-query';

export type BaseQueryOptions<T> = Omit<UseQueryOptions<T, Error, T>, 'queryKey' | 'queryFn'>;
