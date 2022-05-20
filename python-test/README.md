# Integration Tests
This directory contains integration tests that can be run through the UI and/or API


Here's what you'll need to do in order to run these tests:
- Setup your python environment
- Configure the test settings
- Run behave

## Setup your Python environment

<b>You need to have python >= 3.8 installed</b> 

Create a virtual environment: `python3 -m venv name_of_virtualenv`

Activate your virtual environment: `source name_of_virtualenv/bin/activate`

Install the required libraries: `pip install -r requirements.txt`


### Additional configuration of your Python environment for UI tests
- Install Google Chrome :
```
$ sudo curl -sS -o - https://dl-ssl.google.com/linux/linux_signing_key.pub | apt-key add
$ sudo echo "deb [arch=amd64]  http://dl.google.com/linux/chrome/deb/ stable main" >> /etc/apt/sources.list.d/google-chrome.list
$ sudo apt-get -y update
$ sudo apt-get -y install google-chrome-stable
```

- Install ChromeDriver
```
$ wget https://chromedriver.storage.googleapis.com/$(google-chrome --product-version)/chromedriver_linux64.zip
$ unzip chromedriver_linux64.zip
$ sudo mv chromedriver /usr/bin/chromedriver
$ sudo chown root:root /usr/bin/chromedriver
$ sudo chmod +x /usr/bin/chromedriver
```

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
  - URL of the Orb deployment. Do NOT include the protocol (`https://`, `http://` or `mqtt://`).
- **prometheus_username**
  - Mandatory!
  - Your Grafana Cloud Prometheus username
- **prometheus_key**
  - Mandatory!
  - Your Grafana Cloud API Key. Be sure to grant the key a role with metrics push privileges
- **remote_prometheus_endpoint**
  - Mandatory!
  - base URL to send Prometheus metrics to Grafana Cloud> `(ex. prometheus-prod-10-prod-us-central-0.grafana.net)`
- **ignore_ssl_and_certificate_errors**:
  - Bool
  - Replaces HTTPS connections with HTTP and disables SSL certificate validation in MQTT connections.
  - Default value: `False`
- **is_credentials_registered**:
  - Bool
  - If false, register an account with credentials (email and password) used
  - Default value: `False`
- **agent_docker_image**:
  - Docker image of the orb agent.
  - Default value: `ns1labs/orb-agent`
- **agent_docker_tag**:
  - Tag of the Orb agent docker image.
  - Default value: `latest`
- **orb_agent_interface**:
  - Network interface that will be used by pktvisor when running the Orb agent.
  - Default value: `mock`

## Run behave
Simply run `behave`, optionally passing the feature file as follows:

```sh
$ behave --include agentsProvider.feature
```
Output:
```text
@agents
Feature: agent provider # features/agentsProvider.feature:2
  Scenario: Provision agent                                                  # features/agentsProvider.feature:4
    Given that the user is logged in on orb account                                         # features/steps/users.py:10 1.031s
    When a new agent is created                                              # features/steps/control_plane_agents.py:18 1.032s
    And the agent container is started                                       # features/steps/local_agent.py:10 0.217s
    Then the agent status in Orb should be online                            # features/steps/control_plane_agents.py:24 2.556s
    And the container logs should contain the message "sending capabilities" # features/steps/local_agent.py:26 0.023s
1 feature passed, 0 failed, 0 skipped
1 scenario passed, 0 failed, 0 skipped
5 steps passed, 0 failed, 0 skipped, 0 undefined
Took 0m4.858s

```