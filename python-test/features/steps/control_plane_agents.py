from test_config import TestConfig
from utils import random_string, filter_list_by_parameter_start_with
from behave import when, then
from hamcrest import *
import time
import requests

configs = TestConfig.configs()
base_orb_url = "https://" + configs.get('orb_address')
agent_name_prefix = "test_agent_name_"
agent_name = agent_name_prefix + random_string(10)
agent_tag_key = "test_tag_key_" + random_string(4)
agent_tag_value = "test_tag_value_" + random_string(4)


@when('a new agent is created')
def agent_is_created(context):
    agent = create_agent(context.token, agent_name, agent_tag_key, agent_tag_value)
    context.agent = agent


@then('the agent status in Orb should be {status}')
def check_agent_online(context, status):
    token = context.token
    agent_id = context.agent['id']
    expect_container_status(token, agent_id, status)


@then('cleanup agents')
def clean_agents(context):
    """
    Remove all agents starting with 'agent_name_prefix' from the orb

    :param context: Behave class that contains contextual information during the running of tests.
    """
    token = context.token
    agents_list = list_agents(token)
    agents_filtered_list = filter_list_by_parameter_start_with(agents_list, 'name', agent_name_prefix)
    delete_agents(token, agents_filtered_list)


def expect_container_status(token, agent_id, status):
    """
    Keeps fetching agent data from Orb control plane until it gets to
    the expected agent status or this operation times out

    :param (str) token: used for API authentication
    :param (str) agent_id: whose status will be evaluated
    :param (str) status: expected agent status
    """

    time_waiting = 0
    sleep_time = 0.5
    timeout = 10

    while time_waiting < timeout:
        agent = get_agent(token, agent_id)
        agent_status = agent['state']
        if agent_status == status:
            break

        time.sleep(sleep_time)
        time_waiting += sleep_time

    assert_that(time_waiting, is_not(equal_to(timeout)),
                'Agent did not get "' + status + '" after ' + str(timeout) + ' seconds')


def get_agent(token, agent_id):
    """
    Gets an agent from Orb control plane

    :param (str) token: used for API authentication
    :param (str) agent_id: that identifies agent to be fetched
    :returns: (dict) the fetched agent
    """

    get_agents_response = requests.get(base_orb_url + '/api/v1/agents/' + agent_id, headers={'Authorization': token})

    assert_that(get_agents_response.status_code, equal_to(200),
                'Request to get agent id=' + agent_id + ' failed with status=' + str(get_agents_response.status_code))

    return get_agents_response.json()


def list_agents(token):
    """
    Lists all agents from Orb control plane that belong to this user

    :param (str) token: used for API authentication
    :returns: (list) a list of agents
    """

    response = requests.get(base_orb_url + '/api/v1/agents', headers={'Authorization': token})

    assert_that(response.status_code, equal_to(200),
                'Request to list agents failed with status=' + str(response.status_code))

    agents_as_json = response.json()
    return agents_as_json['agents']


def delete_agents(token, list_of_agents):
    """
    Deletes from Orb control plane the agents specified on the given list

    :param (str) token: used for API authentication
    :param (list) list_of_agents: that will be deleted
    """

    for agent in list_of_agents:
        delete_agent(token, agent['id'])


def delete_agent(token, agent_id):
    """
    Deletes an agent from Orb control plane

    :param (str) token: used for API authentication
    :param (str) agent_id: that identifies the agent to be deleted
    """

    response = requests.delete(base_orb_url + '/api/v1/agents/' + agent_id,
                               headers={'Authorization': token})

    assert_that(response.status_code, equal_to(204), 'Request to delete agent id='
                + agent_id + ' failed with status=' + str(response.status_code))


def create_agent(token, name, tag_key, tag_value):
    """
    Creates an agent in Orb control plane

    :param (str) token: used for API authentication
    :param (str) name: of the agent to be created
    :param (str) tag_key: the key of the tag to be added to this agent
    :param (str) tag_value: the value of the tag to be added to this agent
    :returns: (dict) a dictionary containing the created agent data
    """

    response = requests.post(base_orb_url + '/api/v1/agents',
                             json={"name": name, "orb_tags": {tag_key: tag_value}, "validate_only": False},
                             headers={'Content-type': 'application/json', 'Accept': '*/*',
                                      'Authorization': token})
    assert_that(response.status_code, equal_to(201),
                'Request to create agent failed with status=' + str(response.status_code))

    return response.json()

