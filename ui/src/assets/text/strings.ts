export const STRINGS = {
  header: {
    logotype: 'ORB',
  },
  footer: {
    disclaimer: '',
  },
  home: {
    title: 'Welcome to Orb!',
    step: {
      message: 'These steps should help you get started:',
      agent: 'Create an <a href="pages/fleet/agents">Agent</a> to receive instructions on deploying orb-agent and begin observing your infrastructure',
      agent_groups: 'Group your <a href="pages/fleet/agents">Agents</a> into <a href="pages/fleet/groups">Agent Groups</a> so that you can send your agents <a href="pages/datasets/policies">Policies</a>',
      policy: 'Create an <a href="pages/datasets/policies">Agent Policy</a> which will tell your agent how to analyze your data streams',
      sink: 'Setup a <a href="pages/sinks">Sink</a> to be able to send your agent metric output to a time series database for visualizing and alerting',
      dataset: 'Finally, make a Dataset which will send your <a href="pages/datasets/policies">Policy</a> to the selected <a href="pages/fleet/groups">Agent Group</a> and the resulting metrics to the selected <a href="pages/sinks">Sink</a>',
    },
  },
  // Login Pages strings
  login: {
    presentation: {
      description: 'An <a href="https://getorb.io" target="_blank">Open-Source</a> dynamic edge observability platform',
      action: 'Unleash the power of small data with dynamic edge observability',
    },
    form: {
      title: 'Log in or sign up',
      username: 'Username',
      password: 'Password',
      forgot: 'Forgot Password?',
      login: 'Log in',
      register: 'Register',
      nonuser: 'Don\'t have an account?',
    },
  },
  // Fleet Pages strings
  fleet: {
    title: 'Fleet Management',
  },
  // Sink Pages strings
  sink: {
    // sink statuses
    state: {
      active: 'Active',
      error: 'Error',
    },
    // sink.interface name descriptors
    propNames: {
      id: 'id',
      name: 'Name',
      description: 'Description',
      tags: 'Tags',
      state: 'Status',
      error: 'Error',
      backend: 'Service Type',
      config: 'Connection Details',
      config_remote_host: 'Remote Write URL',
      config_username: 'Username',
      config_password: 'Password',
      ts_created: 'Date Created',
    },
    // add page
    add: {
      header: 'New Sink',
    },
    // edit page
    edit: {
      header: 'Edit Sink',
    },
    // delete modal
    delete: {
      header: 'Delete Sink Confirmation',
      body: 'Are you sure you want to delete this Sink? This may cause Datasets which use this Sink to become invalid. This action cannot be undone.',
      warning: '*To confirm, type your Sink name exactly as it appears',
      delete: 'I Understand, Delete This Sink',
      close: 'Close',
    },
    // details modal
    details: {
      header: 'Sink Details',
      close: 'Close',
    },
    // dashboard page
    list: {
      header: 'All Sinks',
      none: 'There are no sinks listed.',
      sink: 'sink',
      total: ['You have', 'total.'],
      error: 'have errors.',
      create: 'New Sink',
      filters: {
        select: 'Filter',
        name: 'Name',
        description: 'Description',
        state: 'Status',
        type: 'Type',
        tags: 'Tags',
      },
    },
  },
  // agents
  agentGroups: {
    // statuses
    state: {
      active: 'Active',
      error: 'Error',
    },
    // agent.interface name descriptors
    propNames: {
      id: 'id',
      name: 'Agent Group Name',
      description: 'Agent Group Description',
      key: 'Key',
      value: 'Value',
      tags: 'Tags',
      state: 'Status',
      error: 'Error',
      ts_created: 'Date Created',
      matches: 'Matches Against',
    },
    // matches
    match: {
      matchAny: 'The Selected Qualifiers Will Match Against',
      matchNone: 'The Selected Qualifiers Do Not Match Any Agent',
      agents: 'Agent(s)' +
        '',
      updated: 'Agent Group matches updated',
      expand: 'Expand',
      collapse: 'Collapse',
    },
    // add page
    add: {
      header: 'New Agent Group',
      step: {
        title1: 'Agent Group Details',
        desc1: 'This is how you will be able to easily identify your Agent Group',
        title2: 'Agent Group Tags',
        desc2: 'Set the tags that will be used to group Agents',
        title3: 'Review & Confirm',
      },
      success: 'Agent Group successfully created',
    },
    // edit page
    edit: {
      header: 'Edit Agent Group',
    },
// delete modal
    delete: {
      header: 'Delete Agent Group Confirmation',
      body: 'Are you sure you want to delete this Agent Group? This may cause Datasets which use this Agent Group to become invalid. This action cannot be undone.',
      warning: '*To confirm, type the Agent Group name exactly as it appears',
      delete: 'I Understand, Delete This Agent Group',
      close: 'Close',
    },
    // details modal
    details: {
      header: 'Agent Group Details',
      close: 'Close',
    },
    // dashboard page
    list: {
      header: 'All Agent Groups',
      none: 'There are no Agents listed.',
      agentGroup: 'agent',
      total: ['You have', 'total.'],
      error: 'have errors.',
      create: 'New Agent Group',
      filters: {
        select: 'Filter',
        name: 'Name',
        description: 'Description',
        state: 'Status',
        type: 'Type',
        tags: 'Tags',
      },
    },
  },
  // agent groups
  agents: {
    // statuses
    state: {
      active: 'Active',
      error: 'Error',
    },
    // agent.interface name descriptors
    propNames: {
      id: 'id',
      name: 'Agent Name',
      description: 'Agent Description',
      key: 'Key',
      value: 'Value',
      orb_tags: 'Orb Tags',
      state: 'Status',
      error: 'Error',
      ts_created: 'Date Created',
    },
    // add page
    add: {
      header: 'New Agent',
      step: {
        title1: 'Agent Details',
        desc1: 'This is how you will be able to easily identify your Agent',
        title2: 'Orb Tags',
        desc2: 'Set the tags that will be used to filter your Agent',
        title3: 'Review & Confirm',
      },
      success: 'Agent successfully created',
    },
    // edit page
    edit: {
      header: 'Update Agent',
    },
// delete modal
    delete: {
      header: 'Delete Confirmation',
      body: 'Are you sure you want to delete this Agent? This action cannot be undone.',
      warning: '*To confirm, type the Agent label exactly as it appears',
      delete: 'I Understand, Delete This Agent',
      close: 'Close',
    },
    // details modal
    details: {
      header: 'Agent Details',
      close: 'Close',
    },
    // dashboard page
    list: {
      header: 'All Agents',
      none: 'There are no Agents listed.',
      agentGroup: 'agent',
      total: ['You have', 'total.'],
      error: 'have errors.',
      create: 'New Agent',
      filters: {
        select: 'Filter',
        name: 'Name',
        description: 'Description',
        state: 'Status',
        type: 'Type',
        tags: 'Tags',
      },
    },
  },
  // stepper cues
  stepper: {
    back: 'Back',
    cancel: 'Cancel',
    next: 'Next',
    save: 'Save',
  },
  // tags cues
  tags: {
    addTag: 'Add New Tag',
    key: 'Tag Key',
    value: 'Tag Value',
  },
};
