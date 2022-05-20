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
    title: 'Agents',
    icon: 'pin-outline',
    link: 'fleet/agents',
  },
  {
    title: 'Agent Groups',
    icon: 'globe-outline',
    link: 'fleet/groups',
  },
  {
    title: 'Policy Management',
    icon: 'layers-outline',
    link: 'datasets/policies',
  },
  {
    title: 'Sink Management',
    icon: 'cloud-upload-outline',
    link: 'sinks',
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

