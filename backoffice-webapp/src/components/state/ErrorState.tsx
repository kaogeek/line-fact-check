import { t } from 'i18next';
import { CircleX } from 'lucide-react';
import EmptyState from './EmptyState';
import type { ReactElement } from 'react';

interface ErrorStateProps {
  title?: string;
  msg?: string;
  icon?: ReactElement;
  showIcon?: boolean;
  action?: ReactElement;
  size?: 'sm' | 'md' | 'lg';
}

export default function ErrorState({
  title,
  msg,
  icon,
  showIcon = true,
  action,
  size = 'md',
  className,
}: ErrorStateProps & { className?: string }) {
  return (
    <EmptyState
      title={title || t('common.error.defaultTitle')}
      msg={msg || t('common.error.defaultMessage')}
      icon={!icon ? <CircleX /> : icon}
      showIcon={showIcon}
      action={action}
      size={size}
      className={className}
    />
  );
}
