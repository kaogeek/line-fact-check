import { Button } from '@/components/ui/button';
import { Menu } from 'lucide-react';
import { useTranslation } from 'react-i18next';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu';
import { getSavedLanguage, saveLanguage } from '@/lib/language-storage';

interface NavigatorBarProps {
  brand: string;
  setIsOpen: (open: boolean) => void;
  className?: string;
}

export default function NavigatorBar({ brand, setIsOpen, className }: NavigatorBarProps) {
  const { t, i18n } = useTranslation();

  const savedLanguage = getSavedLanguage() || i18n.language;

  const changeLanguage = (lng: string) => {
    i18n.changeLanguage(lng);
    saveLanguage(lng);
  };

  return (
    <header>
      <nav
        className={`flex items-center gap-2 p-4 bg-secondary text-secondary-foreground ${className ?? ''} min-h-[60px]`}
        aria-label="Main navigation"
      >
        <Button
          variant="ghost"
          size="icon"
          className="md:hidden"
          onClick={() => setIsOpen(true)}
          aria-label="Open sidebar menu"
        >
          <Menu />
        </Button>
        <strong className="text-lg font-bold md:hidden">{brand}</strong>
        <div className="flex-1"></div>
        <DropdownMenu>
          <DropdownMenuTrigger asChild>
            <Button variant="ghost" size="sm">
              {savedLanguage === 'th' ? 'ไทย' : 'English'}
            </Button>
          </DropdownMenuTrigger>
          <DropdownMenuContent align="end">
            <DropdownMenuItem onClick={() => changeLanguage('en')}>{t('localeSwitcher.english')}</DropdownMenuItem>
            <DropdownMenuItem onClick={() => changeLanguage('th')}>{t('localeSwitcher.thai')}</DropdownMenuItem>
          </DropdownMenuContent>
        </DropdownMenu>
      </nav>
    </header>
  );
}
