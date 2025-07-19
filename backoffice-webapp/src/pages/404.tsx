import { Home } from 'lucide-react';
import NoDataState from '../components/state/NoDataState';
import { Button } from '../components/ui/button';
import { useTranslation } from 'react-i18next';

export default function NotFoundPage() {
  const { t } = useTranslation();

  return (
    <div className="flex items-center justify-center h-screen">
      <NoDataState
        title={t('notFound.title')}
        msg={t('notFound.message')}
        action={
          <Button variant="default">
            <Home />
            <a href="/">{t('notFound.goHome')}</a>
          </Button>
        }
      />
    </div>
  );
}
