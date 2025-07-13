import type { ReactNode } from 'react';
import clsx from 'clsx'; // utility to merge classNames nicely

interface TYProps {
  children: ReactNode;
  className?: string;
}

export function TYH1({ children, className }: TYProps) {
  return (
    <h1 className={clsx('scroll-m-20 text-4xl font-extrabold tracking-tight text-balance', className)}>{children}</h1>
  );
}

export function TYH2({ children, className }: TYProps) {
  return (
    <h2 className={clsx('scroll-m-20 text-3xl font-semibold tracking-tight first:mt-0', className)}>{children}</h2>
  );
}

export function TYH3({ children, className }: TYProps) {
  return <h3 className={clsx('scroll-m-20 text-2xl font-semibold tracking-tight', className)}>{children}</h3>;
}

export function TYH4({ children, className }: TYProps) {
  return <h4 className={clsx('scroll-m-20 text-xl font-semibold tracking-tight', className)}>{children}</h4>;
}

export function TYP({ children, className }: TYProps) {
  return <p className={clsx(className)}>{children}</p>;
}

export function TYLarge({ children, className }: TYProps) {
  return <div className={clsx('text-lg font-semibold', className)}>{children}</div>;
}

export function TYSmall({ children, className }: TYProps) {
  return <small className={clsx('text-sm leading-none font-medium', className)}>{children}</small>;
}

export function TYMuted({ children, className }: TYProps) {
  return <p className={clsx('text-muted-foreground text-sm', className)}>{children}</p>;
}

export function TYLink({ children, className }: TYProps) {
  return <p className={clsx('text-primary font-medium underline text-sm', className)}>{children}</p>;
}
