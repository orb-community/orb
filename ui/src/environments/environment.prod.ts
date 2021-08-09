import {environment as defaultEnvironment} from './environment.defaults';

const {sinksUrl, agentsUrl, agentGroupsUrl, orbApi: {apiUrl, version}} = defaultEnvironment;

export const environment = {
  production: true,

  ...defaultEnvironment,
  // ORB api --prod
  // override all urls prepend /api/v<#>/<service_url>
  sinksUrl: `${apiUrl}${version}${sinksUrl}`,
  agentsUrl: `${apiUrl}${version}${agentsUrl}`,
  agentGroupsUrl: `${apiUrl}${version}${agentGroupsUrl}`,
};
