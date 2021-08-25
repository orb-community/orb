export const SINK_BACKEND_SETTINGS = {
  prometheus: [
    {
      type: 'text',
      input: 'text',
      title: 'Remote Host',
      name: 'remote_host',
      required: true,
    },
    {
      type: 'text',
      input: 'text',
      title: 'Username',
      name: 'username',
      required: true,
    },
    {
      type: 'password',
      input: 'text',
      title: 'Password',
      name: 'password',
      required: true,
    },
  ],
  aws: [
    {
      type: 'text',
      input: 'text',
      title: 'Remote Host',
      name: 'remote_host',
      required: true,
    },
    {
      type: 'text',
      input: 'text',
      title: 'Username',
      name: 'username',
      required: true,
    },
    {
      type: 'password',
      input: 'text',
      title: 'Password',
      name: 'password',
      required: true,
    },
  ],
};

export enum SINK_BACKEND_TYPES {
  prometheus = 'prometheus',
  aws = 'aws',
}

/**
 * Available sink statuses
 */
export enum SINK_STATUS {
  active = 'active',
  error = 'error',
}
