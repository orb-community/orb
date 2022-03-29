import { environment } from '../../environments/environment';
import { MenuItem } from '../shared/interfaces/menu-item.interface';

const MENU: MenuItem[] = [
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
    title: 'Sink Management',
    icon: 'cloud-upload-outline',
    link: 'sinks',
  },
  {
    title: 'Dataset Explorer',
    icon: 'layers-outline',
    children: [
      {
        title: 'Policy Management',
        link: 'datasets/policies',
      },
      {
        title: 'Datasets',
        link: 'datasets/list',
      },
    ],
  },
  {
    title: 'Settings',
    icon: 'settings-2-outline',
  },
];

const DEV_ITEMS: MenuItem[] = [
  {
    title: 'Dev',
    icon: 'layout-outline',
    link: '/pages/dev',
  },
];

export const MENU_ITEMS: MenuItem[] = [
  ...MENU,
  ...environment.production ? [] : DEV_ITEMS,
];
