import { NbMenuItem } from '@nebular/theme';
import { environment } from '../../environments/environment';

export const MENU: NbMenuItem[] = [
  {
    title: 'Home',
    icon: 'layout-outline',
    link: '/pages/home',
    pathMatch: 'full',
    home: true,
  },
  {
    title: 'Agents',
    icon: 'pin-outline',
    link: '/pages/fleet/agents',
    pathMatch: 'full',
  },
  {
    title: 'Agent Groups',
    icon: 'globe-outline',
    link: '/pages/fleet/groups',
    pathMatch: 'full',
  },
  {
    title: 'Policy Management',
    icon: 'layers-outline',
    link: '/pages/datasets/policies',
    pathMatch: 'full',
  },
  {
    title: 'Sink Management',
    icon: 'cloud-upload-outline',
    link: '/pages/sinks',
    pathMatch: 'full',
  },
];

const DEV_ITEMS: NbMenuItem[] = [
  {
    title: 'Dev',
    icon: 'layout-outline',
    link: '/pages/dev',
    pathMatch: 'full',
  },
];

export const MENU_ITEMS = [
  ...MENU,
  ...environment.production ? [] : DEV_ITEMS,
];
