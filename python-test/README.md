# AUTOMATED TESTS

### Configuration of venv

> Create a virtual Environment by running:

`$ python3 -m venv name_of_virtualenv`

> Activate your virtual Environment:

`$ source name_of_virtualenv/bin/activate`

## Install the required libraries

`$ pip install -r requirements.txt`


## Configuring the correct variables
Modified the ./steps/env_variables.env file changing the correct values


- EMAIL=`<email>`
  - Default value = None 
  - user's email
- PASSWORD=`<password>`
  - Default value =  None 
  - user's password
- AGENT_NAME=`<agent-name>`
  - Default value = Random String 
- TAG_KEY=`<tag-key>`
  - Default value = test
- TAG_VALUE=`<tag-value>`
  - Default value = true
- AGENT_DOCKER_TAG=`<tag>`
  - Default value: None (if none, it will default to "latest" by docker)
    - latest/develop
- ORB_ADDRESS=`<address>`
  - Default value = beta.orb.live
- IFACE=`<iface>`
  - Default value = mock
