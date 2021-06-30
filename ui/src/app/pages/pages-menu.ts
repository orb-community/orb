import { environment } from 'environments/environment';
import { NbMenuItem } from '@nebular/theme';

export const MENU_ITEMS: NbMenuItem[] = environment.production ?
  [
    {
      title: 'Home',
      icon: 'home-outline',
      link: '/pages/home',
      home: true,
    },
    {
      title: 'Things',
      icon: 'film-outline',
      link: '/pages/things',
    },
    {
      title: 'Channels',
      icon: 'flip-2-outline',
      link: '/pages/channels',
    },
    {
      title: 'Twins',
      icon: 'copy-outline',
      link: '/pages/twins',
    },
  ] :
  [
    {
      title: 'Home',
      icon: 'home-outline',
      link: '/pages/home',
      home: true,
    },
    {
      title: 'User Groups',
      icon: 'shield-outline',
      link: '/pages/users/groups',
    },
    {
      title: 'Users',
      icon: 'people-outline',
      link: '/pages/users',
    },
    {
      title: 'Things',
      icon: 'film-outline',
      link: '/pages/things',
    },
    {
      title: 'Channels',
      icon: 'flip-2-outline',
      link: '/pages/channels',
    },
    {
      title: 'Twins',
      icon: 'copy-outline',
      link: '/pages/twins',
    },
    {
      title: 'Services',
      icon: 'layers-outline',
      children: [
        {
          title: 'LoRa',
          icon: 'radio-outline',
          link: '/pages/services/lora',
        },
        {
          title: 'OPC-UA',
          icon: 'globe-outline',
          link: '/pages/services/opcua',
        },
        {
          title: 'Gateways',
          icon: 'hard-drive-outline',
          link: '/pages/services/gateways',
        },
      ],
    },
  ];

