const ORB = {
  // introduce primitive ORB api versioning '/api/v1/sinks'
  // TODO not needed at the moment - nginx listens to anything @80
  // orb-ui app proxy routes * to @80/api/v1/*
  orbApi: {
    version: '1', // ORB api version
    apiUrl: '/api/v', // ORB api url prefix
  },
  servicesUrls: {
    sinksUrl: '/sinks',
    sinkBackendsUrl: '/features/sinks',
    agentsUrl: '/agents',
    validateAgentsUrl: '/agents/validate',
    agentGroupsUrl: '/agent_groups',
    validateAgentGroupsUrl: '/agent_groups/validate',
    agentPoliciesUrl: '/policies/agent',
    agentsBackendUrl: '/agents/backends',
    baseBackendComposeUrl: '/agents/backends',
    pktvisorTapsUrl: '/agents/backends/pktvisor/taps',
    pktvisorInputsUrl: '/agents/backends/pktvisor/inputs',
    pktvisorHandlersUrl: '/agents/backends/pktvisor/handlers',
    datasetPoliciesUrl: '/policies/dataset',
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
    loginUrl: '/tokens',
    httpAdapterUrl: '/http',
    readerUrl: '/reader',
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
  loginUrl: '/tokens',
  httpAdapterUrl: '/http',
  readerUrl: '/reader',
  readerPrefix: 'channels',
  readerSuffix: 'messages',
  exportConfigFile: '/configs/export/config.toml',

  // PACTSAFE
  PS: {
    // site id
    SID: '',
    // group key
    GROUP_KEY: '',
  },

  // MAINTENANCE
  MAINTENANCE: '',

  // expose ORB routes and api versioning
  orbApi: {urlKeys: Object.keys(ORB.servicesUrls), ...ORB.orbApi, servicesUrls: ORB.servicesUrls},
  ...ORB.servicesUrls,
};
