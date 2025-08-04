import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Search } from 'lucide-react';
import { useState } from 'react';
import { useTranslation } from 'react-i18next';

interface MessageGroupSearchBarProps {
  initCodeLike?: string;
  initMessageLike?: string;
  handleSearch: (criteria: { codeLike?: string; messageLike?: string }) => void;
}

export default function MessageGroupSearchBar({
  initCodeLike,
  initMessageLike,
  handleSearch,
}: MessageGroupSearchBarProps) {
  const { t } = useTranslation();
  const [codeLike, setCodeLike] = useState<string | undefined>(initCodeLike);
  const [messageLike, setMessageLike] = useState<string | undefined>(initMessageLike);

  function handleCodeLikeChange(event: React.ChangeEvent<HTMLInputElement>) {
    setCodeLike(event.target.value);
  }

  function handleMessageLikeChange(event: React.ChangeEvent<HTMLInputElement>) {
    setMessageLike(event.target.value);
  }

  function handleSearchClick() {
    handleSearch({
      ...(codeLike && { codeLike }),
      ...(messageLike && { messageLike }),
    });
  }

  return (
    <div className="flex flex-col md:flex-row gap-3 w-full items-stretch">
      <div className="flex flex-1 flex-col sm:flex-row gap-2">
        <div className="min-w-[120px]">
          <label className="text-sm font-medium text-gray-700 mb-1 block">{t('messageGroup.searchLabel.code')}</label>
          <Input
            className="w-full"
            placeholder={t('messageGroup.searchPlaceholder.code')}
            value={codeLike ?? ''}
            onChange={handleCodeLikeChange}
          />
        </div>
        <div className="min-w-[120px]">
          <label className="text-sm font-medium text-gray-700 mb-1 block">
            {t('messageGroup.searchLabel.message')}
          </label>
          <Input
            className="w-full"
            placeholder={t('messageGroup.searchPlaceholder.message')}
            value={messageLike ?? ''}
            onChange={handleMessageLikeChange}
          />
        </div>
      </div>
      <div className="flex items-end">
        <Button onClick={handleSearchClick}>
          <Search className="mr-2 h-4 w-4" />
          {t('common.searchButton')}
        </Button>
      </div>
    </div>
  );
}
