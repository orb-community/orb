import { environment } from 'environments/environment';
import { NbMenuItem } from '@nebular/theme';

export const MENU_ITEMS: NbMenuItem[] = environment.production ?
  [
    {
      title: 'Agent Groups Explorer',
      icon: 'home-outline',
      link: '/pages/agent-groups-explorer',
    },
    {
      title: 'Dataset Explorer',
      icon: 'home-outline',
      link: '/pages/dataset-explorer',
    },
    {
      title: 'Fleet Management',
      icon: 'home-outline',
      link: '/pages/fleet-management',
    },
    {
      title: 'Sink Management',
      icon: 'home-outline',
      link: '/pages/agent-groups-explorer',
    },
  ] :
  [
    {
      title: 'Agent Groups Explorer',
      icon: 'home-outline',
      link: '/pages/agent-groups-explorer',
    },
    {
      title: 'Dataset Explorer',
      icon: 'home-outline',
      link: '/pages/dataset-explorer',
    },
    {
      title: 'Fleet Management',
      icon: 'home-outline',
      link: '/pages/fleet-management',
    },
    {
      title: 'Sink Management',
      icon: 'home-outline',
      link: '/pages/agent-groups-explorer',
    },
  ];

