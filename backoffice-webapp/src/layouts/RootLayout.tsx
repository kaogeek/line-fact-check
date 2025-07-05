import NavigatorBar from '@/components/navigator/NavigatorBar';
import SideBar from '@/components/navigator/SideBar';
import { APP_NAME } from '@/constants/app';
import { useState } from 'react';
import { Outlet } from 'react-router';

export default function RootLayout() {
  const [isOpen, setIsOpen] = useState(false);
  return (
    <div className="flex h-screen">
      <SideBar
        brand={APP_NAME}
        menuList={[
          {
            label: 'Topic',
            link: '/topic',
          },
          {
            label: 'Dashboard',
            link: '/dashboard',
          },
          {
            label: 'Logout',
            onClick: () => {
              console.log('Logout');
            },
          },
        ]}
        isOpen={isOpen}
        setIsOpen={setIsOpen}
      />

      <main className="flex-1 overflow-auto">
        <NavigatorBar brand={APP_NAME} setIsOpen={setIsOpen} className="md:hidden" />
        <Outlet />
      </main>
    </div>
  );
}
