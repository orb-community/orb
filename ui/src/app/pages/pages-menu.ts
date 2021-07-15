import { environment } from 'environments/environment';
import { NbMenuItem } from '@nebular/theme';

export const MENU_ITEMS: NbMenuItem[] = environment.production ?
  [
    {
      title: 'Agent Groups Explorer',
      icon: 'home-outline',
      link: '/pages/agents',
    },
    {
      title: 'Dataset Explorer',
      icon: 'home-outline',
      link: '/pages/datasets',
    },
    {
      title: 'Fleet Management',
      icon: 'home-outline',
      link: '/pages/fleets',
    },
    {
      title: 'Settings',
      icon: 'home-outline',
      children: [
        {
          title: 'Sink Management',
          icon: 'home-outline',
          link: '/pages/sinks',
        },
        {
          title: 'Selector Management',
          icon: 'home-outline',
          link: '/pages/selectors',
        },
        {
          title: 'Policy Management',
          icon: 'home-outline',
          link: '/pages/policies',
        },
      ],
      link: '/pages/agent-groups-explorer',
    },
  ] :
    [
      {
        title: 'Agent Groups Explorer',
        icon: 'home-outline',
        link: '/pages/agents',
      },
      {
        title: 'Dataset Explorer',
        icon: 'home-outline',
        link: '/pages/datasets',
      },
      {
        title: 'Fleet Management',
        icon: 'home-outline',
        link: '/pages/fleets',
      },
      {
        title: 'Settings',
        icon: 'home-outline',
        children: [
          {
            title: 'Sink Management',
            icon: 'home-outline',
            link: '/pages/sinks',
          },
          {
            title: 'Selector Management',
            icon: 'home-outline',
            link: '/pages/selectors',
          },
          {
            title: 'Policy Management',
            icon: 'home-outline',
            link: '/pages/policies',
          },
        ],
        link: '/pages/agent-groups-explorer',
      },
    ];

