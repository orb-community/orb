const ORB = {
  // introduce primitive ORB api versioning '/api/v1/sinks'
  orbApi: {
    version: '1', // ORB api version
    apiUrl: '/api/v', // ORB api url prefix
  },
  servicesUrls: {
    sinksUrl: '/sinks',
    agentsUrl: '/agents',
    agentGroupsUrl: '/agent_groups',
    loginUrl: '/tokens',
  },
};

export const environment = {
  usersUrl: '/users',
  groupsUrl: '/groups',
  membersUrl: '/members',
  usersVersionUrl: '/version',
  requestPassUrl: '/password/reset-request',
  resetPassUrl: '/password/reset',
  changePassUrl: '/password',
  thingsUrl: '/things',
  twinsUrl: '/twins',
  statesUrl: '/states',
  channelsUrl: '/channels',
  bootstrapConfigsUrl: '/bootstrap/things/configs',
  bootstrapUrl: '/bootstrap/things/bootstrap',
  connectUrl: '/connect',
  browseUrl: '/browse',

  httpAdapterUrl: '/http',
  readerUrl: '/reader',
  readerPrefix: 'channels',
  readerSuffix: 'messages',

  mqttWsUrl: window['env']['mqttWsUrl'] || 'ws://localhost/mqtt',
  exportConfigFile: '/configs/export/config.toml',
  // expose ORB routes and api versioning
  orbApi: {urlKeys: Object.keys(ORB.servicesUrls), ...ORB.orbApi, servicesUrls: ORB.servicesUrls},
  ...ORB.servicesUrls,
};
