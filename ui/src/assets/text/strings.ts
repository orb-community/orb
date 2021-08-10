export const STRINGS = {
  header: {
    logotype: 'ORB',
  },
  footer: {
    disclaimer: '',
  },
  home: {
    title: 'Orb Observation Overview',
  },
  // Fleet Pages strings
  fleet: {
    title: 'Fleet Management',
  },
  // Sink Pages strings
  sink: {
    // sink statuses
    status: {
      active: 'Active',
      error: 'Error',
    },
    // sink.interface name descriptors
    propNames: {
      id: 'id',
      name: 'Name',
      description: 'Description',
      tags: 'Tags',
      status: 'Status',
      error: 'Error',
      backend: 'Service Type',
      config: 'Connection Details',
      config_remote_host: 'Remote Host',
      config_username: 'Username',
      config_password: 'Password',
      ts_created: 'Date Created',
    },
    // add page
    add: {
      header: 'Create New Sink',
    },
    // edit page
    edit: {
      header: 'Update Sink',
    },
    // delete modal
    delete: {
      header: 'Delete Sink Confirmation',
      body: 'Are you sure you want to delete this sink? This may cause policies which use this sink to become invalid. This action cannot be undone.',
      warning: '*To confirm, type your sink label exactly as it appears',
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
        status: 'Status',
        type: 'Type',
        tags: 'Tags',
      },
    },
  },
  // stepper cues
  stepper: {
    back: 'Back',
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
