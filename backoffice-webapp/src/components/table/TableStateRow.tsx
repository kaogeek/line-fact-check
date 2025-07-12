import React from 'react';

export interface TableStateRowProps {
  colSpan: number;
  children: React.ReactNode;
  className?: string;
}

export default function TableStateRow({ colSpan = 6, children, className }: TableStateRowProps) {
  return (
    <tr>
      <td colSpan={colSpan} className={`text-center py-8 ${className ?? ''}`}>
        {children}
      </td>
    </tr>
  );
}
