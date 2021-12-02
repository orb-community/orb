import { environment as defaultEnvironment } from './environment.defaults';

export const environment = {
  production: false,
  GTAGID: 'G-387CGPZQF0',
  // PACTSAFE
  PS: {
    // site id
    SID: `${process.env.PS_SID}`,
    // group key
    GROUP_KEY: `${process.env.PS_GROUP_KEY}`,
  },
  ...defaultEnvironment,
};
