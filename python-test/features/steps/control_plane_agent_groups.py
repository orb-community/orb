from test_config import TestConfig
from local_agent import get_orb_agent_logs
from users import get_auth_token
from utils import random_string, filter_list_by_parameter_start_with, generate_random_string_with_predefined_prefix,\
    create_tags_set, check_logs_contain_message_and_name
from behave import given, when, then, step
from hamcrest import *
import requests
import time

configs = TestConfig.configs()
agent_group_name_prefix = 'test_group_name_'
agent_group_description = "This is an agent group"
base_orb_url = configs.get('base_orb_url')


@when("an Agent Group is created with same tag as the agent")
def create_agent_group_matching_agent(context):
    agent_group_name = agent_group_name_prefix + random_string()
    tags = context.agent["orb_tags"]
    context.agent_group_data = create_agent_group(context.token, agent_group_name, agent_group_description,
                                                  tags)
    group_id = context.agent_group_data['id']
    context.agent_groups[group_id] = agent_group_name


@step("an Agent Group is created with {tags_type} orb tag(s): {orb_tags}")
def create_new_agent_group(context, tags_type, orb_tags):
    agent_group_name = generate_random_string_with_predefined_prefix(agent_group_name_prefix)
    context.orb_tags = create_tags_set(tags_type, orb_tags)
    if len(context.orb_tags) == 0:
        context.agent_group_data = create_agent_group(context.token, agent_group_name, agent_group_description,
                                                      context.orb_tags, 400)
    else:
        context.agent_group_data = create_agent_group(context.token, agent_group_name, agent_group_description,
                                                      context.orb_tags)
        group_id = context.agent_group_data['id']
        context.agent_groups[group_id] = agent_group_name


@step("Agent Group creation response must be an error with message '{message}'")
def error_response_message(context, message):
    response = list(context.agent_group_data.items())[0]
    response_key, response_value = response[0], response[1]
    assert_that(response_key, equal_to('error'),
                'Response of invalid agent group creation must be an error')
    assert_that(response_value, equal_to(message), "Unexpected message for error")


@then("one agent must be matching on response field matching_agents")
def matching_agent(context):
    matching_total_agents = context.agent_group_data['matching_agents']['total']
    matching_online_agents = context.agent_group_data['matching_agents']['online']
    assert_that(matching_total_agents, equal_to(1))
    assert_that(matching_online_agents, equal_to(1))


@step("the group to which the agent is linked is removed")
def remove_group(context):
    delete_agent_group(context.token, context.agent_group_data['id'])


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
    agent_group_name = generate_random_string_with_predefined_prefix(agent_group_name_prefix)
    agent_tags = agent['orb_tags']
    context.agent_group_data = create_agent_group(context.token, agent_group_name, agent_group_description, agent_tags)
    group_id = context.agent_group_data['id']
    context.agent_groups[group_id] = agent_group_name
    matching_agent(context)


@step('the container logs contain the message "{text_to_match}" referred to each group within {'
      'time_to_wait} seconds')
def check_logs_for_group(context, text_to_match, time_to_wait):
    text_found, groups_to_which_subscribed = check_subscription(time_to_wait, context.agent_groups.values(),
                                                                text_to_match, context.container_id)
    assert_that(text_found, is_(True), f"Message {text_to_match} was not found in the agent logs for group(s)"
                                       f"{set(context.agent_groups.values()).difference(groups_to_which_subscribed)}!")


def create_agent_group(token, name, description, tags, expected_status_code=201):
    """
    Creates an agent group in Orb control plane

    :param (str) token: used for API authentication
    :param (str) name: of the agent to be created
    :param (str) description: description of group
    :param (dict) tags: dict with all pairs key:value that will be used as tags
    :returns: (dict) a dictionary containing the created agent group data
    :param (int) expected_status_code: expected request's status code. Default:201 (happy path).
    """

    json_request = {"name": name, "description": description, "tags": tags}
    headers_request = {'Content-type': 'application/json', 'Accept': '*/*', 'Authorization': token}

    response = requests.post(base_orb_url + '/api/v1/agent_groups', json=json_request, headers=headers_request)
    assert_that(response.status_code, equal_to(expected_status_code),
                'Request to create agent group failed with status=' + str(response.status_code))

    return response.json()


def list_agent_groups(token, limit=100):
    """
    Lists up to 100 agent groups from Orb control plane that belong to this user

    :param (str) token: used for API authentication
    :param (int) limit: Size of the subset to retrieve (max 100). Default = 100
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


def check_subscription(time_to_wait, agent_groups_names, text_to_match, container_id):
    groups_to_which_subscribed = set()
    time_waiting = 0
    sleep_time = 0.5
    timeout = int(time_to_wait)
    while time_waiting < timeout:
        for name in agent_groups_names:
            logs = get_orb_agent_logs(container_id)
            text_found, log_line = check_logs_contain_message_and_name(logs, text_to_match, name, "group_name")
            if text_found is True:
                groups_to_which_subscribed.add(log_line['group_name'])
                if set(groups_to_which_subscribed) == set(agent_groups_names):
                    return True, groups_to_which_subscribed
        time.sleep(sleep_time)
        time_waiting += sleep_time
    return False, groups_to_which_subscribed
