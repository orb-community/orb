import { environment as defaultEnvironment } from './environment.defaults';

const {orbApi: {apiUrl, version, urlKeys, servicesUrls}} = defaultEnvironment;

export const environment = {
  production: true,
  GTMID: 'G-387CGPZQF0',
  ...defaultEnvironment,
  // ORB api --prod
  // override all urls prepend /api/v<#>/<service_url>
  ...urlKeys.reduce(
    (acc, cur) => {
      acc[cur] = `${apiUrl}${version}${servicesUrls[cur]}`;
      return acc;
    },
    {}),
};
