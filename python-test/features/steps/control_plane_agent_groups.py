import time
from control_plane_agents import create_agent, agent_name, agent_tag_key, agent_tag_value, get_agent, \
    expect_container_status, base_orb_url
from local_agent import run_agent_container
from behave import given, when, then
from hamcrest import *
from test_config import TestConfig
import random
import string
import requests

configs = TestConfig.configs()
random_agent_name_group = ''.join(random.choices(string.ascii_letters, k=10))  # k sets the number of characters
agent_group_name = 'test_group' + random_agent_name_group
agent_group_description = "This is an agent group"


@given("that an agent already exists and be {status}")
def check_if_agents_exist(context, status):
    agent = create_agent(context.token, agent_name, agent_tag_key, agent_tag_value)
    context.agent = agent
    token = context.token
    orb_address = configs.get('orb_address')
    interface = configs.get('orb_agent_interface', 'mock')
    agent_docker_image = configs.get('agent_docker_image', 'ns1labs/orb-agent')
    image_tag = ':' + configs.get('agent_docker_tag', 'latest')
    agent_image = agent_docker_image + image_tag
    env_vars = {"ORB_CLOUD_ADDRESS": orb_address,
                "ORB_CLOUD_MQTT_ID": context.agent['id'],
                "ORB_CLOUD_MQTT_CHANNEL_ID": context.agent['channel_id'],
                "ORB_CLOUD_MQTT_KEY": context.agent['key'],
                "PKTVISOR_PCAP_IFACE_DEFAULT": interface}
    context.container_id = run_agent_container(agent_image, env_vars)
    agent_id = context.agent['id']
    existing_agents = get_agent(token, agent_id)
    assert_that(len(existing_agents), greater_than(0), "Agent not created")
    expect_container_status(token, agent_id, status)


@when("an Agent Group is created with same tag as agent")
def creat_agent_group(context):
    context.agent_group_data = create_agent_group(context.token, agent_group_name, agent_group_description,
                                                  agent_tag_key, agent_tag_value)
    time.sleep(2)


@then("one matching agent must be seen")
def matching_agent(context):
    matching_total_agents = context.agent_group_data['matching_agents']['total']
    matching_online_agents = context.agent_group_data['matching_agents']['online']
    assert_that(matching_total_agents, equal_to(1))
    assert_that(matching_online_agents, equal_to(1))


def create_agent_group(token, name, description, tag_key, tag_value):
    """
    Creates an agent group in Orb control plane

    :param (str) token: used for API authentication
    :param (str) name: of the agent to be created
    :param (str) description: description of group
    :param (str) tag_key: the key of the tag to be added to this agent
    :param (str) tag_value: the value of the tag to be added to this agent
    :returns: (dict) a dictionary containing the created agent group data
    """

    response = requests.post(base_orb_url + '/api/v1/agent_groups',
                             json={"name": name, "description": description, "tags": {tag_key: tag_value}},
                             headers={'Content-type': 'application/json', 'Accept': '*/*',
                                      'Authorization': token})
    assert_that(response.status_code, equal_to(201),
                'Request to create agent failed with status=' + str(response.status_code))

    return response.json()


def list_agent_groups(token):
    """
    Lists all agent groups from Orb control plane that belong to this user

    :param (str) token: used for API authentication
    :returns: (list) a list of agent groups
    """

    response = requests.get(base_orb_url + '/api/v1/agent_groups', headers={'Authorization': token})

    assert_that(response.status_code, equal_to(200),
                'Request to list agent groups failed with status=' + str(response.status_code))

    agent_groups_as_json = response.json()
    return agent_groups_as_json['agentGroups']


def delete_agent_groups(token, list_of_agent_groups):
    """
    Deletes from Orb control plane the agent groups specified on the given list

    :param (str) token: used for API authentication
    :param (list) list_of_agent_groups: that will be deleted
    """

    for agent_Groups in list_of_agent_groups:
        delete_agent_group(token, agent_Groups['id'])


def delete_agent_group(token, agent_group_id):
    """
    Deletes an agent group from Orb control plane

    :param (str) token: used for API authentication
    :param (str) agent_group_id: that identifies the agent group to be deleted
    """

    response = requests.delete(base_orb_url + '/api/v1/agent_groups/' + agent_group_id,
                               headers={'Authorization': token})

    assert_that(response.status_code, equal_to(204), 'Request to delete agent group id='
                + agent_group_id + ' failed with status=' + str(response.status_code))
