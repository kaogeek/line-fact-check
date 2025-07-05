import avatar from '../assets/feast.avif';
import { NavLink } from 'react-router';
import { X } from 'lucide-react';
import { Button } from './ui/button';
import UserAvatar from './UserAvatar';

interface SideBarProps {
  brand: string;
  menuList: SideBarMenu[];
  isOpen: boolean;
  setIsOpen: (open: boolean) => void;
}

type SideBarMenu =
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
    <aside className="flex flex-col h-full shadow-lg w-[238px] bg-white">
      <header className="p-4 text-lg font-bold text-center">{brand}</header>

      <section className="p-4">
        <UserAvatar avatarUrl={avatar} name="Username" />
      </section>

      <nav className="flex-1 flex flex-col">
        {menuList.map((menu, idx) =>
          menu.link ? (
            <NavLink
              key={idx}
              to={menu.link}
              className={({ isActive }) =>
                ['p-4 block', 'transition duration-300', 'hover:bg-gray-200', isActive ? 'bg-gray-300' : ''].join(' ')
              }
              onClick={() => setIsOpen(false)}
            >
              {menu.label}
            </NavLink>
          ) : (
            <button
              key={idx}
              onClick={menu.onClick}
              className="p-4 text-left block transition duration-300 hover:bg-gray-200 w-full"
            >
              {menu.label}
            </button>
          )
        )}
      </nav>

      <footer className="p-4 text-sm text-center text-gray-500">Version 1.0.0</footer>
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
