import { environment as defaultEnvironment } from './environment.defaults';

const {orbApi: {apiUrl, version, urlKeys}} = defaultEnvironment;

export const environment = {
  production: true,

  ...defaultEnvironment,
  // ORB api --prod
  // override all urls prepend /api/v<#>/<service_url>
  ...urlKeys.map(key => ({[key]: `${apiUrl}${version}${defaultEnvironment[key]}`})),
};
