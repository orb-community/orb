## Login

* Register an unregistered account
* Check if user is not enabled to register an account previously registered
* Register an account only with valid email address and password
* Register an account with full name, company name, valid email address and password
* Register an account with invalid email address
* Register an account with invalid password
* Check if email and password is required fields
* Login with valid credentials
* Login with invalid credentials
* Request password with registered email address
* Request password with unregistered email address


## Agents

* Check if total agent on agent page is correct
* Check if regex matches
* Create agent without tags
* Create agent with one tag
* Create agent with multiple tags (test whether tag removal is enabled)
* Create agent with duplicate name
* Provision agent before agent group (check status and last activity)
* Provision agent after agent group
* Provision agent before dataset
* Provision agent after dataset
* Test filters (name, status, tags, search)
* Check if agent details is correctly displayed (agent name, channel ID, created on, status and tags)
* Check if is possible to edit an agent through the details modal
* Edit agent name
* Edit agent tag
* Save agent without tag
* Insert tags in agents created without tags
* Check if is possible cancel operations with no change
* Stop/remove agent container (check status)
* Remove agent (check if removal changes matching agents only on referred group)
* Check if removal is not allowed without correct name

## Agent Groups

* Check if total agent groups on agent groups page is correct
* Check if regex matches
* Create agent group with duplicate name
* Create agent group with description
* Create agent group without description
* Check if user is not enabled to create group without tag
* Create agent group with one tag (validate if only match agents are linked to group)
* Create agent group with multiple tags (validate if only match agents are linked to group)
* Test filters (name, description, agents, tags, search)
* Test if is possible visualize matching agents
* Check if agent groups details is correctly displayed (agent group name, agent group description, matches against, and tags)
* Check if is possible to edit an agent group through the details modal
* Check if is possible cancel operations with no change
* Edit agent name
* Edit agent description
* Edit agent group tag
* Insert/remove tags in existing agent group (check if correctly matching agents remains)
* Remove agent group
* Remove agent group to which an agent is linked (check container logs)
* Check if removal is not allowed without correct name

## Sinks

* Check if total sinks on sinks page is correct
* Check if regex matches
* Create sink with duplicate name
* Create sink with description
* Create sink without description
* Check if remote host, username and password are required to create a sink
* Check is user is allowed to create a sink without tags (create sink without tags)
* Create sink with one tag
* Create sink with multiple tags
* Test filters (name, description, type, status, tags, search)
* Check if sink details is correctly displayed (name, description, service type, remote host, username and tags)
* Check if is possible to edit a sink through the details modal
* Edit sink name
* Edit sink description
* Edit sink remote host
* Edit sink username
* Edit sink password
* Edit sink tags
* Check if is possible cancel operations with no change
* Remove sink
* Check if removal is not allowed without correct name

## Policies

* Check if policy creation when no agent is provisioned has an alert message (backend)
* Check if total policies on policies page is correct
* Check if regex matches
* Create policy with duplicate name
* Create policy with description
* Create policy without description
* Create policy with dhcp handler
* Create policy with dns handler
* Create policy with net handler
* Create policy with multiple handlers
* Create broken policy
* Test filters (name, description, version, last modified and search)
* Check if policies details is correctly displayed (name, description, policy backend, version)
* Check if is possible to edit a policy through the details modal
* Edit policy name
* Edit policy description
* Edit policy handler
* Check if is possible cancel operations with no change
* Remove policy
* Check if removal is not allowed without correct name

## Datasets

* Check if total datasets on datasets page is correct
* Check if regex matches
* Create dataset (check if group, policy and sink are required)
* Check if datasets details is correctly displayed (name, validity, group, policy and sink)
* Check if is possible cancel operations with no change
* Test filter
* Check if is possible to edit a dataset through the details modal
* Edit dataset name
* Edit dataset sink
* Remove dataset
* Check if removal is not allowed without correct name

## Integration tests

 - [Check if sink is active while scraping metrics](integration/sink_active_while_scraping_metrics.md)
 - [Check if sink with invalid credentials becomes active](integration/sink_error_invalid_credentials.md)
 - [Check if after 30 minutes without data sink becomes idle](integration/sink_idle_30_minutes.md)
 - [Provision agent before group (check if agent subscribes to the group)](integration/provision_agent_before_group.md)
 - [Provision agent after group (check if agent subscribes to the group)](integration/provision_agent_after_group.md)
 - [Create agent with tag matching existing group linked to a valid dataset](integration/multiple_agents_subscribed_to_a_group.md)
 - [Apply multiple policies to a group](integration/apply_multiple_policies.md)
 - [Apply multiple policies to a group and remove one policy](integration/remove_one_of_multiple_policies.md)
 - [Apply multiple policies to a group and remove all of them](integration/remove_all_policies.md)
 - [Apply multiple policies to a group and remove one dataset](integration/remove_one_of_multiple_datasets.md)
 - [Apply multiple policies to a group and remove all datasets](integration/remove_all_datasets.md)
 - [Apply the same policy twice to the agent](integration/apply_policy_twice.md)
 - [Delete sink linked to a dataset, create another one and edit dataset using new sink](integration/change_sink_on_dataset.md)
 - [Remove one of multiples datasets that apply the same policy to the agent](integration/remove_one_dataset_of_multiples_with_same_policy.md)
 - [Remove group (invalid dataset, agent logs)](integration/remove_group.md)
 - [Remove sink (invalid dataset, agent logs)](integration/remove_sink.md)
 - [Remove policy (invalid dataset, agent logs, heartbeat)](integration/remove_policy.md)
 - [Remove dataset (check agent logs, heartbeat)](integration/remove_dataset.md)
 - [Remove agent container (logs, agent groups matches)](integration/remove_agent_container.md)
 - [Remove agent container force (logs, agent groups matches)](integration/remove_agent_container_force.md)
 - [Remove agent (logs, agent groups matches)](integration/remove_agent.md)
## Pktvisor Agent

* Providing Orb-agent with sample commands
* Providing Orb-agent with configuration files
* Providing Orb-agent with advanced auto-provisioning setup
* Providing more than one Orb-agent with different ports
* Providing Orb-agent using mocking interface
* Providing a Orb-agent with a wrong interface
* Pull the latest orb-agent image, build and run the agent