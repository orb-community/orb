import { environment } from 'environments/environment';
import { NbMenuItem } from '@nebular/theme';

export const MENU_ITEMS: NbMenuItem[] = environment.production ?
  [
    {
      title: 'Fleet Management',
      icon: 'home-outline',
      link: '/pages/fleet-management',
    },
    {
      title: 'DataSet Explorer',
      icon: 'home-outline',
      link: '/pages/dataset-explorer',
    },
  ] :
  [
    {
      title: 'Fleet Management',
      icon: 'home-outline',
      link: '/pages/fleet-management',
    },
    {
      title: 'DataSet Explorer',
      icon: 'home-outline',
      link: '/pages/dataset-explorer',
    },
  ];

