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
    title: 'Fleet Management',
    icon: 'cube-outline',
    pathMatch: 'prefix',
    children: [
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
    ],
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

export function updateMenuItems(pageName: string) {
  MENU_ITEMS.forEach(item => {
    if (item.children) {
      item.children.forEach(child => {
        child.selected = child.title === pageName;
      });
      item.selected = item.children.some(child => child.selected);
    } else {
      item.selected = item.title === pageName;
    }
  });
}
