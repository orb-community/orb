import {environment as defaultEnvironment} from './environment.defaults';
import {environment as envVars} from './environment.env';

const {orbApi: {apiUrl, version, urlKeys, servicesUrls}} = defaultEnvironment;

export const environment = {
    production: true,
    GTAGID: 'G-387CGPZQF0',
    ...defaultEnvironment,
    ...envVars,
    // ORB api --prod
    // override all urls prepend /api/v<#>/<service_url>
    ...urlKeys.reduce(
        (acc, cur) => {
            acc[cur] = `${apiUrl}${version}${servicesUrls[cur]}`;
            return acc;
        },
        {}),
};
