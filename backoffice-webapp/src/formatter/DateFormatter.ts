// src/utils/formatDate.ts
import { format } from 'date-fns';

export function formatDate(dateInput: Date | string | number): string {
  const date = new Date(dateInput);
  return format(date, 'dd/MM/yyyy HH:mm');
}
