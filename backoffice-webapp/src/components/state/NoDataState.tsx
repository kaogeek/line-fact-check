import { type ReactElement } from 'react';
import EmptyState from './EmptyState';
import { useTranslation } from 'react-i18next';

interface ErrorStateProps {
  title?: string;
  msg?: string;
  icon?: ReactElement;
  showIcon?: boolean;
  action?: ReactElement;
  size?: 'sm' | 'md' | 'lg';
}

export default function NoDataState({ title, msg, icon, showIcon = true, action, size = 'md' }: ErrorStateProps) {
  const { t } = useTranslation();

  return (
    <EmptyState
      title={title || t('common.noData.defaultTitle')}
      msg={msg || t('common.noData.defaultMessage')}
      icon={!icon ? <img src="/assets/state/task-empty.svg"></img> : icon}
      showIcon={showIcon}
      action={action}
      size={size}
    ></EmptyState>
  );
}
