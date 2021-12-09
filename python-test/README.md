# Integration Tests

Here's what you'll need to do in order to run these tests:
- Setup your python environment
- Configure the test settings
- Run behave

## Setup your Python environment
Create a virtual environment: `python3 -m venv name_of_virtualenv`

Activate your virtual environment: `source name_of_virtualenv/bin/activate`

Install the required libraries: `pip install -r requirements.txt`

## Test settings
Create the test config file from the template: `cp test_config.ini.tpl test_config.ini`.

Then fill in the correct values:

- **email**:
  - Mandatory!
  - Orb user's email
- **password**:
  - Mandatory!
  - Orb user's password
- **orb_address**:
  - Mandatory!
  - URL of the Orb deployment. Do NOT include the protocol (`https://` or `mqtt://`).
- **agent_docker_image**:
  - Docker image of the orb agent.
  - Default value: `ns1labs/orb-agent`
- **agent_docker_tag**:
  - Tag of the Orb agent docker image.
  - Default value: `latest`
- **orb_agent_interface**:
  - Network interface that will be used by pktvisor when running the Orb agent.
  - Default value: `mock`
- **prometheus_username**
  - Mandatory for running the tests in [sinks feature](./features/sinks.feature)
  - Your Grafana Cloud Prometheus username
- **prometheus_key**
  - Mandatory for running the tests in [sinks feature](./features/sinks.feature)
  - Your Grafana Cloud API Key. Be sure to grant the key a role with metrics push privileges
- **remote_prometheus_endpoint**
  - Mandatory for running the tests in [sinks feature](./features/sinks.feature)
  - base URL to send Prometheus metrics to Grafana Cloud> `(ex. prometheus-prod-10-prod-us-central-0.grafana.net)`

## Run behave
Simply run `behave`, optionally passing the feature file as follows:

```sh
$ behave --include agentsProvider.feature
```
Output:
```text
@agentGroups
Feature: agent groups creation # features/agentGroups.feature:2

  Scenario: Create Agent Group                                                                                # features/agentGroups.feature:4
    Given the Orb user logs in                                                                                # features/steps/users.py:10 0.668s
    And that an agent already exists and is online                                                            # features/steps/control_plane_agents.py:15 4.647s
    When an Agent Group is created with same tag as the agent                                                 # features/steps/control_plane_agent_groups.py:15 0.720s
    Then one agent must be matching on response field matching_agents                                         # features/steps/control_plane_agent_groups.py:22 0.000s
    And the container logs should contain the message "completed RPC subscription to group" within 10 seconds # features/steps/local_agent.py:26 0.550s

@agents
Feature: agent provider # features/agentsProvider.feature:2

  Scenario: Provision agent                                                                    # features/agentsProvider.feature:4
    Given the Orb user logs in                                                                 # features/steps/users.py:10 0.677s
    When a new agent is created                                                                # features/steps/control_plane_agents.py:30 0.761s
    And the agent container is started                                                         # features/steps/local_agent.py:11 0.262s
    Then the agent status in Orb should be online                                              # features/steps/control_plane_agents.py:39 3.619s
    And the container logs should contain the message "sending capabilities" within 10 seconds # features/steps/local_agent.py:26 0.050s

@datasets
Feature: datasets creation # features/datasets.feature:2

  Scenario: Create Dataset                                                                                 # features/datasets.feature:4
    Given the Orb user logs in                                                                             # features/steps/users.py:10 0.648s
    And that an agent already exists and is online                                                         # features/steps/control_plane_agents.py:15 4.121s
    And referred agent is subscribed to a group                                                            # features/steps/control_plane_agent_groups.py:43 0.759s
    And that a sink already exists                                                                         # features/steps/control_plane_sink.py:60 0.617s
    And that a policy already exists                                                                       # features/steps/control_plane_policies.py:38 2.549s
    When a new dataset is created using referred group, sink and policy ID                                 # features/steps/control_plane_datasets.py:10 0.593s
    Then the container logs should contain the message "managing agent policy from core" within 10 seconds # features/steps/local_agent.py:26 0.034s
    Then the container logs should contain the message "scraped metrics for policy" within 120 seconds     # features/steps/local_agent.py:26 113.350s
    Then referred sink must have active state on response                                                  # features/steps/control_plane_sink.py:40 0.549s

@policies
Feature: policy creation # features/policies.feature:2

  Scenario: Create a policy                                        # features/policies.feature:4
    Given the Orb user logs in                                     # features/steps/users.py:10 0.633s
    And that an agent already exists and is online                 # features/steps/control_plane_agents.py:15 4.001s
    When a new policy is created                                   # features/steps/control_plane_policies.py:12 0.615s
    Then referred policy must be listened on the orb policies list # features/steps/control_plane_policies.py:18 0.612s

@sinks
Feature: sink creation # features/sinks.feature:2

  Scenario: Create Sink using Prometheus                       # features/sinks.feature:4
    Given that the user has the prometheus/grafana credentials # features/steps/control_plane_sink.py:12 0.000s
    And the Orb user logs in                                   # features/steps/users.py:10 0.921s
    When a new sink is created                                 # features/steps/control_plane_sink.py:30 0.675s
    Then referred sink must have unknown state on response     # features/steps/control_plane_sink.py:40 0.654s

5 features passed, 0 failed, 0 skipped
5 scenarios passed, 0 failed, 0 skipped
27 steps passed, 0 failed, 0 skipped, 0 undefined
Took 2m23.288s


```