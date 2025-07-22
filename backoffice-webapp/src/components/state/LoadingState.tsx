import { t } from 'i18next';
import Loader from '../Loader';
import { cn } from '@/lib/utils';

const LoadingState = ({ className }: { className?: string }) => {
  return (
    <>
      <div className={cn(`container min-h-[60vh] flex flex-col justify-center items-center`, className)}>
        <Loader className="mx-auto mb-4" />
        <h3 className="text-xl font-medium">{t('common.loading')}</h3>
      </div>
    </>
  );
};

export default LoadingState;
