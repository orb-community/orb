# AUTOMATED TESTS

### Configuration of venv

> Create a virtual Environment by running:

`$ python3 -m venv name_of_virtualenv`

> Activate your virtual Environment:

`$ source name_of_virtualenv/bin/activate`

## Install the required libraries

`$ pip install -r requirements.txt`

## Configuring the correct variables

Create the file `test_config.ini` by making a copy of the template `test_config.ini.tpl`, then fill it with the correct
values.

- email:
  - Mandatory!
  - Orb user's email
- password:
  - Mandatory!
  - Orb user's password
- orb_address:
  - Mandatory!
  - URL of the Orb deployment. Do NOT include the protocol (`https://` or `mqtt://`)
- agent_docker_tag:
  - Tag of the Orb agent docker image
  - Default value: `latest`
- orb_agent_interface:
  - Network interface that will be used by pktvisor when running the Orb agent
  - Default value: `mock`
