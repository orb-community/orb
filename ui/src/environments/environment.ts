import { environment as defaultEnvironment } from './environment.defaults';
import { environment as envVars } from './environment.env';

export const environment = {
  production: false,
  appPrefix: '',
  GTAGID: 'G-387CGPZQF0',
  ...defaultEnvironment,
  ...envVars,
};
