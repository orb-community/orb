import { environment as defaultEnvironment } from './environment.defaults';

export const environment = {
  production: false,
  GTMID: '',
  ...defaultEnvironment,
};
