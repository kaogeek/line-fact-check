import { cn } from '@/lib/utils';
import type { SideBarMenu } from './SideBar';
import { NavLink } from 'react-router';

interface SideBarMenuProps {
  menu: SideBarMenu;
}

const btnStyle = 'px-4 py-2 block rounded-lg transition duration-300 hover:text-muted-foreground hover:bg-muted w-full';
const activeBtnStyle = 'bg-primary text-primary-foreground';

export default function SideBarMenuBtn({ menu }: SideBarMenuProps) {
  return (
    <div>
      {menu.link ? (
        <NavLink to={menu.link} className={({ isActive }) => cn(btnStyle, isActive && activeBtnStyle)}>
          {menu.label}
        </NavLink>
      ) : (
        <button onClick={menu.onClick} className={cn(btnStyle, 'text-left')}>
          {menu.label}
        </button>
      )}
    </div>
  );
}
