import { NbMenuItem } from '@nebular/theme';
import { environment } from '../../environments/environment';

export const MENU: NbMenuItem[] = [
  {
    title: 'Home',
    icon: 'layout-outline',
    link: 'home',
    home: true,
  },
  {
    title: 'Fleet Management',
    icon: 'cube-outline',
    children: [
      {
        title: 'Agents',
        link: 'fleet/agents',
      },
      {
        title: 'Agent Groups',
        link: 'fleet/groups',
      },
    ],
  },
  {
    title: 'Dataset Explorer',
    icon: 'layers-outline',
    children: [
      {
        title: 'Datasets',
      },
      {
        title: 'Policy Management',
        link: 'policies',
      },
    ],
  },
  {
    title: 'Sink Management',
    icon: 'layers-outline',
    link: 'sinks',
  },
  {
    title: 'Settings',
    icon: 'settings-2-outline',
  },
];

const DEV_ITEMS: NbMenuItem[] = [
  {
    title: 'Dev',
    icon: 'layout-outline',
    link: '/pages/dev',
  },
];

export const MENU_ITEMS = [
  ...MENU,
  ...environment.production ? [] : DEV_ITEMS,
];

