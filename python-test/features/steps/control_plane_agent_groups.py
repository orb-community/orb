from test_config import TestConfig
from control_plane_agents import create_agent, get_agent, expect_container_status, base_orb_url, agent_name_prefix
from local_agent import run_local_agent_container
from users import get_auth_token
from utils import random_string, filter_list_by_parameter_start_with
from behave import given, when, then
from hamcrest import *
import requests

configs = TestConfig.configs()
agent_group_name_prefix = 'test_group_name_'
agent_group_description = "This is an agent group"


@when("an Agent Group is created with same tag as the agent")
def creat_agent_group(context):
    agent_group_name = agent_group_name_prefix + random_string()
    context.agent_group_data = create_agent_group(context.token, agent_group_name, agent_group_description,
                                                  context.agent_tag_key, context.agent_tag_value)


@then("one agent must be matching on response field matching_agents")
def matching_agent(context):
    matching_total_agents = context.agent_group_data['matching_agents']['total']
    matching_online_agents = context.agent_group_data['matching_agents']['online']
    assert_that(matching_total_agents, equal_to(1))
    assert_that(matching_online_agents, equal_to(1))


@then('cleanup agent group')
def clean_agent_groups(context):
    """
    Remove all agent groups starting with 'agent_group_name_prefix' from the orb

    :param context: Behave object that contains contextual information during the running of tests.
    """
    token = context.token
    agent_groups_list = list_agent_groups(token)
    agent_groups_filtered_list = filter_list_by_parameter_start_with(agent_groups_list, 'name', agent_group_name_prefix)
    delete_agent_groups(token, agent_groups_filtered_list)


@given("referred agent is subscribed to a group")
def subscribe_agent_to_a_group(context):
    agent = context.agent
    agent_group_name = agent_group_name_prefix + random_string(4)
    agent_tag_key = list(agent['orb_tags'].keys())[0]
    agent_tag_value = agent['orb_tags'][agent_tag_key]
    context.agent_group_data = create_agent_group(context.token, agent_group_name, agent_group_description,
                                                  agent_tag_key, agent_tag_value)
    matching_agent(context)


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


def list_agent_groups(token, limit=100):
    """
    Lists up to 100 agent groups from Orb control plane that belong to this user

    :param (int) limit: Size of the subset to retrieve.
    :param (str) token: used for API authentication
    :returns: (list) a list of agent groups
    """

    response = requests.get(base_orb_url + '/api/v1/agent_groups', headers={'Authorization': token},
                            params={"limit": limit})

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
