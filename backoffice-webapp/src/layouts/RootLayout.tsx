import NavigatorBar from '@/components/navigator/NavigatorBar';
import SideBar from '@/components/navigator/SideBar';
import { APP_NAME } from '@/constants/app';
import { useState } from 'react';
import { Outlet } from 'react-router';
import { useTranslation } from 'react-i18next';

export default function RootLayout() {
  const { t } = useTranslation();
  const [isOpen, setIsOpen] = useState(false);
  return (
    <div className="flex h-screen">
      <SideBar
        brand={APP_NAME}
        menuList={[
          {
            label: t('menu.topic'),
            link: '/topic',
          },
          {
            label: t('menu.dashboard'),
            link: '/dashboard',
          },
          {
            label: t('menu.logout'),
            onClick: () => {
              console.log('Logout');
            },
          },
        ]}
        isOpen={isOpen}
        setIsOpen={setIsOpen}
      />

      <main className="flex-1 h-screen">
        <NavigatorBar brand={APP_NAME} setIsOpen={setIsOpen} />
        <Outlet />
      </main>
    </div>
  );
}
