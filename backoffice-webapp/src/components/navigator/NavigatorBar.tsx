import { Button } from '@/components/ui/button';
import { Menu } from 'lucide-react';

interface NavigatorBarProps {
  brand: string;
  setIsOpen: (open: boolean) => void;
  className?: string;
}

export default function NavigatorBar({ brand, setIsOpen, className }: NavigatorBarProps) {
  return (
    <header>
      <nav
        className={`flex items-center gap-2 p-4 bg-primary text-primary-foreground ${className ?? ''}`}
        aria-label="Main navigation"
      >
        <Button
          variant="secondary"
          size="icon"
          className="md:hidden"
          onClick={() => setIsOpen(true)}
          aria-label="Open sidebar menu"
        >
          <Menu />
        </Button>
        <strong className="text-lg font-bold">{brand}</strong>
      </nav>
    </header>
  );
}
