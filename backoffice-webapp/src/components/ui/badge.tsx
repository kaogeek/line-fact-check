import * as React from 'react';
import { Slot } from '@radix-ui/react-slot';
import { cva, type VariantProps } from 'class-variance-authority';

import { cn } from '@/lib/utils';

const badgeVariants = cva(
  'inline-flex items-center justify-center rounded-md border px-2 py-0.5 text-xs font-medium w-fit whitespace-nowrap shrink-0 [&>svg]:size-3 gap-1 [&>svg]:pointer-events-none focus-visible:border-ring focus-visible:ring-ring/50 focus-visible:ring-[3px] aria-invalid:ring-destructive/20 dark:aria-invalid:ring-destructive/40 aria-invalid:border-destructive transition-[color,box-shadow] overflow-hidden',
  {
    variants: {
      variant: {
        default: 'border-transparent bg-primary text-primary-foreground [a&]:hover:bg-primary/90',
        secondary: 'border-transparent bg-secondary text-secondary-foreground [a&]:hover:bg-secondary/90',
        destructive:
          'border-transparent bg-destructive text-white [a&]:hover:bg-destructive/90 focus-visible:ring-destructive/20 dark:focus-visible:ring-destructive/40 dark:bg-destructive/60',
        outline: 'text-foreground [a&]:hover:bg-accent [a&]:hover:text-accent-foreground',
        // Custom
        gray: 'border-transparent bg-gray-soft text-gray-soft-foreground',
        blue: 'border-transparent bg-blue-soft text-blue-soft-foreground',
        success: 'border-transparent bg-success-soft text-success-soft-foreground',
        warning: 'border-transparent bg-warning-soft text-warning-soft-foreground',
        orange: 'border-transparent bg-orange-soft text-orange-soft-foreground',
        danger: 'border-transparent bg-danger-soft text-danger-soft-foreground',
        purple: 'border-transparent bg-purple-soft text-purple-soft-foreground',
        strongGray: 'border-transparent bg-gray text-gray-foreground',
        strongBlue: 'border-transparent bg-blue text-blue-foreground',
        strongSuccess: 'border-transparent bg-success text-success-foreground',
        strongWarning: 'border-transparent bg-warning text-warning-foreground',
        strongOrange: 'border-transparent bg-orange text-orange-foreground',
        strongDanger: 'border-transparent bg-danger text-danger-foreground',
        strongPurple: 'border-transparent bg-purple text-purple-foreground',
      },
    },
    defaultVariants: {
      variant: 'default',
    },
  }
);

function Badge({
  className,
  variant,
  asChild = false,
  ...props
}: React.ComponentProps<'span'> & VariantProps<typeof badgeVariants> & { asChild?: boolean }) {
  const Comp = asChild ? Slot : 'span';

  return <Comp data-slot="badge" className={cn(badgeVariants({ variant }), className)} {...props} />;
}

export { Badge, badgeVariants };
