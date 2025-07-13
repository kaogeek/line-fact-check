import avatar from '../../assets/feast.avif';
import { X } from 'lucide-react';
import { Button } from '../ui/button';
import UserAvatar from '../UserAvatar';
import SideBarMenuBtn from './SideBarMenuBtn';
import { TYMuted } from '../Typography';

interface SideBarProps {
  brand: string;
  menuList: SideBarMenu[];
  isOpen: boolean;
  setIsOpen: (open: boolean) => void;
}

export type SideBarMenu =
  | {
      label: string;
      link: string;
      onClick?: never;
    }
  | {
      label: string;
      link?: never;
      onClick: () => void;
    };

export default function SideBar({ brand, menuList, isOpen, setIsOpen }: SideBarProps) {
  const SidebarContent = (
    <aside className="flex flex-col h-full shadow-lg w-[238px] bg-secondary text-secondary-foreground">
      <header className="p-4 text-lg font-bold text-center text-primary">{brand}</header>
      <hr className="border-t border-secondary-light" />
      <section className="p-4">
        <UserAvatar avatarUrl={avatar} name="Username" />
      </section>
      <hr className="border-t border-secondary-light" />
      <nav className="flex-1 flex flex-col p-4 gap-2">
        {menuList.map((menu, idx) => (
          <SideBarMenuBtn key={idx} menu={menu} />
        ))}
      </nav>
      <hr className="border-t border-secondary-light" />
      <footer className="p-4 text-center">
        <TYMuted>Version 1.0.0</TYMuted>
      </footer>
    </aside>
  );

  return (
    <>
      {/* Sidebar for Desktop */}
      <div className="hidden md:block">{SidebarContent}</div>

      {/* Off-canvas for Mobile */}
      {isOpen && (
        <div className="lg:hidden">
          <Button
            variant="secondary"
            size="icon"
            className="fixed top-4 right-4 z-51"
            onClick={() => setIsOpen(false)}
            aria-label="Close sidebar"
          >
            <X />
          </Button>
          <div className="fixed inset-0 z-50 flex">{SidebarContent}</div>
        </div>
      )}
    </>
  );
}
