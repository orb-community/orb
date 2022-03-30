export const STRINGS = {
  header: {
    title: 'ORB',
  },
  dashboard: {
    title: 'Welcome to Orb!',
    step: {
      message: 'These steps should help you get started:',
      agent: 'Create an <a href="pages/fleet/agents">Agent</a> to receive instructions on deploying orb-agent and begin observing your infrastructure',
      agent_groups: 'Group your <a href="pages/fleet/agents">Agents</a> into <a href="pages/fleet/groups">Agent Groups</a> so that you can send your agents <a href="pages/datasets/policies">Policies</a>',
      policy: 'Create an <a href="pages/datasets/policies">Agent Policy</a> which will tell your agent how to analyze your data streams',
      sink: 'Setup a <a href="pages/sinks">Sink</a> to be able to send your agent metric output to a time series database for visualizing and alerting',
      dataset: 'Finally, make a <a href="pages/datasets/list">Dataset</a> which will send your <a href="pages/datasets/policies">Policy</a> to the selected <a href="pages/fleet/groups">Agent Group</a> and the resulting metrics to the selected <a href="pages/sinks">Sink</a>',
    },
  },
}
