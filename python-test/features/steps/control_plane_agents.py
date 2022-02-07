from test_config import TestConfig
from utils import random_string, filter_list_by_parameter_start_with, generate_random_string_with_predefined_prefix, create_tags_set
from local_agent import run_local_agent_container
from behave import given, when, then, step
from hamcrest import *
import time
import requests

configs = TestConfig.configs()
agent_name_prefix = "test_agent_name_"
base_orb_url = configs.get('base_orb_url')


@given("that an agent with {tags_type} orb tag(s): {orb_tags} already exists and is {status}")
def check_if_agents_exist(context, tags_type, orb_tags, status):
    context.agent_name = generate_random_string_with_predefined_prefix(agent_name_prefix)
    context.orb_tags = create_tags_set(tags_type, orb_tags)
    context.agent = create_agent(context.token, context.agent_name, context.orb_tags)
    token = context.token
    run_local_agent_container(context, "default")
    agent_id = context.agent['id']
    existing_agents = get_agent(token, agent_id)
    assert_that(len(existing_agents), greater_than(0), "Agent not created")
    expect_container_status(token, agent_id, status)


@step('a new agent is created with {tags_type} orb tag(s): {orb_tags}')
def agent_is_created(context, tags_type, orb_tags):
    context.agent_name = generate_random_string_with_predefined_prefix(agent_name_prefix)
    context.orb_tags = create_tags_set(tags_type, orb_tags)
    context.agent = create_agent(context.token, context.agent_name, context.orb_tags)


@when('a new agent is created with tags matching an existing group')
def agent_is_created_matching_group(context):
    context.agent_name = agent_name_prefix + random_string(10)
    agent = create_agent(context.token, context.agent_name, context.orb_tags)
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


@step("{amount_of_datasets} datasets are linked with each policy on agent's heartbeat")
def multiple_dataset_for_policy(context, amount_of_datasets):
    agent = get_agent(context.token, context.agent['id'])
    for policy_id in context.list_agent_policies_id:
        assert_that(len(agent['last_hb_data']['policy_state'][policy_id]['datasets']), equal_to(int(amount_of_datasets)),
                    f"Amount of datasets linked with policy {policy_id} failed")


@step("this agent's heartbeat shows that {amount_of_policies} policies are successfully applied")
def list_policies_applied_to_an_agent(context, amount_of_policies):
    time_waiting = 0
    sleep_time = 0.5
    timeout = 180
    context.list_agent_policies_id = list()
    while time_waiting < timeout:
        agent = get_agent(context.token, context.agent['id'])
        if 'policy_state' in agent['last_hb_data'].keys():
            context.list_agent_policies_id = list(agent['last_hb_data']['policy_state'].keys())
            if len(context.list_agent_policies_id) == int(amount_of_policies):
                break
        time.sleep(sleep_time)
        time_waiting += sleep_time

    assert_that(len(context.list_agent_policies_id), equal_to(int(amount_of_policies)),
                f"Amount of policies applied to this agent failed with {context.list_agent_policies_id} policies")
    assert_that(sorted(context.list_agent_policies_id), equal_to(sorted(context.policies_created.keys())),
                "Policies linked with the agent is not the same as the created by test process")
    for policy_id in context.list_agent_policies_id:
        assert_that(agent['last_hb_data']['policy_state'][policy_id]["state"], equal_to('running'),
                    f"policy {policy_id} is not running")


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

    assert_that(agent_status, is_(equal_to(status)),
                f"Agent did not get '{status}' after {str(timeout)} seconds, but was '{agent_status}'")


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


def list_agents(token, limit=100):
    """
    Lists up to 100 agents from Orb control plane that belong to this user

    :param (str) token: used for API authentication
    :param (int) limit: Size of the subset to retrieve. (max 100). Default = 100
    :returns: (list) a list of agents
    """

    response = requests.get(base_orb_url + '/api/v1/agents', headers={'Authorization': token}, params={"limit": limit})

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


def create_agent(token, name, tags):
    """
    Creates an agent in Orb control plane

    :param (str) token: used for API authentication
    :param (str) name: of the agent to be created
    :param (str) tag_key: the key of the tag to be added to this agent
    :param (str) tag_value: the value of the tag to be added to this agent
    :returns: (dict) a dictionary containing the created agent data
    """

    json_request = {"name": name, "orb_tags": tags, "validate_only": False}
    headers_request = {'Content-type': 'application/json', 'Accept': '*/*',
                       'Authorization': token}

    response = requests.post(base_orb_url + '/api/v1/agents', json=json_request, headers=headers_request)
    assert_that(response.status_code, equal_to(201),
                'Request to create agent failed with status=' + str(response.status_code))

    return response.json()
