from test_config import TestConfig
from utils import random_string, filter_list_by_parameter_start_with, generate_random_string_with_predefined_prefix,\
    create_tags_set, find_files, threading_wait_until, check_port_is_available
from local_agent import run_local_agent_container, run_agent_config_file
from control_plane_agent_groups import return_matching_groups, tags_to_match_k_groups
from behave import given, then, step
from hamcrest import *
from datetime import datetime
import requests
import os
from agent_config_file import FleetAgent
import yaml
from yaml.loader import SafeLoader
import re

configs = TestConfig.configs()
agent_name_prefix = "test_agent_name_"
orb_url = configs.get('orb_url')


@given("that an agent with {orb_tags} orb tag(s) already exists and is {status}")
def check_if_agents_exist(context, orb_tags, status):
    context.agent_name = generate_random_string_with_predefined_prefix(agent_name_prefix)
    context.orb_tags = create_tags_set(orb_tags)
    context.agent = create_agent(context.token, context.agent_name, context.orb_tags)
    context.agent_key = context.agent["key"]
    token = context.token
    run_local_agent_container(context, "available")
    agent_id = context.agent['id']
    existing_agents = get_agent(token, agent_id)
    assert_that(len(existing_agents), greater_than(0), "Agent not created")
    timeout = 30
    agent_status = expect_container_status(token, agent_id, status, timeout=timeout)
    assert_that(agent_status, is_(equal_to(status)),
                f"Agent did not get '{status}' after {str(timeout)} seconds, but was '{agent_status}'")


@step('a new agent is created with {orb_tags} orb tag(s)')
def agent_is_created(context, orb_tags):
    context.agent_name = generate_random_string_with_predefined_prefix(agent_name_prefix)
    context.orb_tags = create_tags_set(orb_tags)
    context.agent = create_agent(context.token, context.agent_name, context.orb_tags)
    context.agent_key = context.agent["key"]


@step('a new agent is created with orb tags matching {amount_of_group} existing group')
def agent_is_created_matching_group(context, amount_of_group):
    context.agent_name = agent_name_prefix + random_string(10)
    all_used_tags = tags_to_match_k_groups(context.token, amount_of_group, context.agent_groups)
    agent = create_agent(context.token, context.agent_name, all_used_tags)
    context.agent = agent
    context.agent_key = context.agent["key"]


@then('the agent status in Orb should be {status}')
def check_agent_online(context, status):
    timeout = 10
    token = context.token
    agent_id = context.agent['id']
    agent_status = expect_container_status(token, agent_id, status, timeout=timeout)
    assert_that(agent_status, is_(equal_to(status)),
                f"Agent did not get '{status}' after {str(timeout)} seconds, but was '{agent_status}'")


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
        assert_that(len(agent['last_hb_data']['policy_state'][policy_id]['datasets']),
                    equal_to(int(amount_of_datasets)),
                    f"Amount of datasets linked with policy {policy_id} failed")


@step("this agent's heartbeat shows that {amount_of_policies} policies are successfully applied and has status {"
      "policies_status}")
def list_policies_applied_to_an_agent_and_referred_status(context, amount_of_policies, policies_status):
    list_policies_applied_to_an_agent(context, amount_of_policies)
    for policy_id in context.list_agent_policies_id:
        assert_that(context.agent['last_hb_data']['policy_state'][policy_id]["state"], equal_to(policies_status),
                    f"policy {policy_id} is not {policies_status}")


@step("this agent's heartbeat shows that {amount_of_policies} policies are successfully applied to the agent")
def list_policies_applied_to_an_agent(context, amount_of_policies):
    context.agent, context.list_agent_policies_id = get_policies_applied_to_an_agent(context.token, context.agent['id'],
                                                                                     amount_of_policies, timeout=180)

    assert_that(len(context.list_agent_policies_id), equal_to(int(amount_of_policies)),
                f"Amount of policies applied to this agent failed with {context.list_agent_policies_id} policies")


@step("this agent's heartbeat shows that {amount_of_groups} groups are matching the agent")
def list_groups_matching_an_agent(context, amount_of_groups):
    groups_matching, context.groups_matching_id = return_matching_groups(context.token, context.agent_groups,
                                                                         context.agent)
    context.list_groups_id = get_groups_to_which_agent_is_matching(context.token, context.agent['id'],
                                                                   context.groups_matching_id, timeout=180)
    assert_that(len(context.list_groups_id), equal_to(int(amount_of_groups)),
                f"Amount of groups matching the agent failed with {context.list_groups_id} groups")
    assert_that(sorted(context.list_groups_id), equal_to(sorted(context.groups_matching_id)),
                "Groups matching the agent is not the same as the created by test process")


@step("edit the agent tags and use {orb_tags} orb tag(s)")
def editing_agent_tags(context, orb_tags):
    agent = get_agent(context.token, context.agent["id"])
    context.orb_tags = create_tags_set(orb_tags)
    edit_agent(context.token, context.agent["id"], agent["name"], context.orb_tags, expected_status_code=200)
    context.agent = get_agent(context.token, context.agent["id"])


@step("edit the agent tags and use orb tags matching {amount_of_group} existing group")
def agent_is_edited_matching_group(context, amount_of_group):
    all_used_tags = tags_to_match_k_groups(context.token, amount_of_group, context.agent_groups)
    agent = get_agent(context.token, context.agent["id"])
    edit_agent(context.token, agent["id"], agent["name"], all_used_tags, expected_status_code=200)
    context.agent = get_agent(context.token, context.agent["id"])


@step("edit the agent name")
def editing_agent_name(context):
    agent = get_agent(context.token, context.agent["id"])
    agent_new_name = generate_random_string_with_predefined_prefix(agent_name_prefix, 5)
    edit_agent(context.token, context.agent["id"], agent_new_name, agent['orb_tags'], expected_status_code=200)
    context.agent = get_agent(context.token, context.agent["id"])
    assert_that(context.agent["name"], equal_to(agent_new_name), "Agent name editing failed")


@step("edit the agent name and edit agent tags using {orb_tags} orb tag(s)")
def editing_agent_name_and_tags(context, orb_tags):
    agent_new_name = generate_random_string_with_predefined_prefix(agent_name_prefix, 5)
    context.orb_tags = create_tags_set(orb_tags)
    edit_agent(context.token, context.agent["id"], agent_new_name, context.orb_tags, expected_status_code=200)
    context.agent = get_agent(context.token, context.agent["id"])
    assert_that(context.agent["name"], equal_to(agent_new_name), "Agent name editing failed")


@step("agent must have {amount_of_tags} tags")
def check_agent_tags(context, amount_of_tags):
    agent = get_agent(context.token, context.agent["id"])
    assert_that(len(dict(agent["orb_tags"])), equal_to(int(amount_of_tags)), "Amount of orb tags failed")


@then("remove all the agents .yaml generated on test process")
def remove_agent_config_files(context):
    all_files_generated = find_files(agent_name_prefix, ".yaml", context.dir_path)
    if len(all_files_generated) > 0:
        for file in all_files_generated:
            os.remove(file)


@threading_wait_until
def check_agent_exists_on_backend(token, agent_name, event=None):
    agent = None
    all_agents = list_agents(token)
    for agent in all_agents:
        if agent_name == agent['name']:
            event.set()
            return agent, event.is_set()
    return agent, event.is_set()


@step("an agent is self-provisioned via a configuration file on port {port} with {agent_tags} agent tags and has "
      "status {status}")
def provision_agent_using_config_file(context, port, agent_tags, status):
    agent_name = f"{agent_name_prefix}{random_string(3)}"
    interface = configs.get('orb_agent_interface', 'mock')
    orb_url = configs.get('orb_url')
    base_orb_address = configs.get('orb_address')
    context.dir_path = create_agent_config_file(context.token, agent_name, interface, agent_tags, orb_url,
                                                base_orb_address, port, existing_agent_groups=context.agent_groups)
    context.container_id = run_agent_config_file(context.dir_path, agent_name)
    context.agent, is_agent_created = check_agent_exists_on_backend(context.token, agent_name, timeout=10)
    assert_that(is_agent_created, equal_to(True), f"Agent {agent_name} not found")
    assert_that(context.agent, is_not(None), f"Agent {agent_name} not correctly created")
    agent_id = context.agent['id']
    existing_agents = get_agent(context.token, agent_id)
    assert_that(len(existing_agents), greater_than(0), "Agent not created")
    expect_container_status(context.token, agent_id, status)


@step("remotely restart the agent")
def reset_agent_remotely(context):
    context.considered_timestamp_reset = datetime.now().timestamp()
    headers_request = {'Content-type': 'application/json', 'Accept': '*/*', 'Authorization': context.token}
    response = requests.post(f"{orb_url}/api/v1/agents/{context.agent['id']}/rpc/reset", headers=headers_request)
    assert_that(response.status_code, equal_to(200),
                'Request to restart agent failed with status=' + str(response.status_code))


@threading_wait_until
def expect_container_status(token, agent_id, status, event=None):
    """
    Keeps fetching agent data from Orb control plane until it gets to
    the expected agent status or this operation times out

    :param (str) token: used for API authentication
    :param (str) agent_id: whose status will be evaluated
    :param (str) status: expected agent status
    :param (obj) event: threading.event
    """

    agent = get_agent(token, agent_id)
    agent_status = agent['state']
    if agent_status == status:
        event.set()
        return agent_status
    return agent_status


def get_agent(token, agent_id):
    """
    Gets an agent from Orb control plane

    :param (str) token: used for API authentication
    :param (str) agent_id: that identifies agent to be fetched
    :returns: (dict) the fetched agent
    """

    get_agents_response = requests.get(orb_url + '/api/v1/agents/' + agent_id, headers={'Authorization': token})

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

    response = requests.get(orb_url + '/api/v1/agents', headers={'Authorization': token}, params={"limit": limit})

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

    response = requests.delete(orb_url + '/api/v1/agents/' + agent_id,
                               headers={'Authorization': token})

    assert_that(response.status_code, equal_to(204), 'Request to delete agent id='
                + agent_id + ' failed with status=' + str(response.status_code))


def create_agent(token, name, tags):
    """
    Creates an agent in Orb control plane

    :param (str) token: used for API authentication
    :param (str) name: of the agent to be created
    :param (dict) tags: orb agent tags
    :returns: (dict) a dictionary containing the created agent data
    """

    json_request = {"name": name, "orb_tags": tags, "validate_only": False}
    headers_request = {'Content-type': 'application/json', 'Accept': '*/*',
                       'Authorization': token}

    response = requests.post(orb_url + '/api/v1/agents', json=json_request, headers=headers_request)
    assert_that(response.status_code, equal_to(201),
                'Request to create agent failed with status=' + str(response.status_code))

    return response.json()


def edit_agent(token, agent_id, name, tags, expected_status_code=200):
    """
    :param (str) token: used for API authentication
    :param (str) agent_id: that identifies the agent to be deleted
    :param (str) name: of the agent to be created
    :param (dict) tags: orb agent tags
    :param (int) expected_status_code: expected request's status code. Default:200 (happy path).
    :return: (dict) a dictionary containing the edited agent data
    """

    json_request = {"name": name, "orb_tags": tags, "validate_only": False}
    headers_request = {'Content-type': 'application/json', 'Accept': '*/*',
                       'Authorization': token}
    response = requests.put(orb_url + '/api/v1/agents/' + agent_id, json=json_request, headers=headers_request)
    assert_that(response.status_code, equal_to(expected_status_code),
                'Request to edit agent failed with status=' + str(response.status_code))

    return response.json()


@threading_wait_until
def get_policies_applied_to_an_agent(token, agent_id, amount_of_policies, event=None):
    """

    :param (str) token: used for API authentication
    :param (str) agent_id: that identifies the agent to be deleted
    :param (int) amount_of_policies: amount of policies that is expected to be applied to the agents
    :param (obj) event: threading.event
    :return:  (dict) agent -> the fetched agent and (list) list_agent_policies_id -> list with the ids of the policies
    that are applied to the agent
    """
    list_agent_policies_id = list()
    agent = get_agent(token, agent_id)
    if 'policy_state' in agent['last_hb_data'].keys():
        list_agent_policies_id = list(agent['last_hb_data']['policy_state'].keys())
        if len(list_agent_policies_id) == int(amount_of_policies):
            event.set()
            return agent, list_agent_policies_id
    return agent, list_agent_policies_id


@threading_wait_until
def get_groups_to_which_agent_is_matching(token, agent_id, groups_matching_ids, event=None):
    """

    :param (str) token: used for API authentication
    :param (str) agent_id: that identifies the agent to be deleted
    :param (list) groups_matching_ids: list with the ids of the groups to with the agent should be subscribed
    :param (obj) event: threading.event
    :return: (list) list_groups_id -> list with the ids of the groups to with the agent is subscribed
    """
    list_groups_id = list()
    agent = get_agent(token, agent_id)
    if 'group_state' in agent['last_hb_data'].keys():
        list_groups_id = list(agent['last_hb_data']['group_state'].keys())
        if sorted(list_groups_id) == sorted(groups_matching_ids):
            event.set()
            return list_groups_id
    return list_groups_id


def create_agent_config_file(token, agent_name, iface, agent_tags, orb_url, base_orb_address, status_port='available',
                             existing_agent_groups=None):
    """
    Create a file .yaml with configs of the agent that will be provisioned

    :param (str) token: used for API authentication
    :param agent_name: name of the agent that will be created
    :param (str) iface: network interface
    :param (str) agent_tags: agent tags
    :param (str) orb_url: entire orb url
    :param (str) base_orb_address: base orb url address
    :param (str) status_port: status of the port on which agent must try to run. Default: available
    :return: path to the directory where the agent config file was created
    """
    assert_that(status_port, any_of(equal_to("available"), equal_to("unavailable")), "Unexpected value for port")
    availability = {"available": True, "unavailable": False}
    if re.match(r"matching (\d+|all|the) group*", agent_tags):
        amount_of_group = re.search(r"(\d+|all|the)", agent_tags).groups()[0]
        all_used_tags = tags_to_match_k_groups(token, amount_of_group, existing_agent_groups)
        tags = {"tags": all_used_tags}
    else:
        tags = {"tags": create_tags_set(agent_tags)}
    agent_config_file = FleetAgent.config_file_of_agent_tap_pcap(agent_name, token, iface, orb_url, base_orb_address)
    agent_config_file = yaml.load(agent_config_file, Loader=SafeLoader)
    agent_config_file['orb'].update(tags)
    port = check_port_is_available(availability[status_port])
    agent_config_file['orb']['backends']['pktvisor'].update({"api_port": f"{port}"})
    agent_config_file = yaml.dump(agent_config_file)
    cwd = os.getcwd()
    dir_path = os.path.dirname(cwd)
    with open(f"{dir_path}/{agent_name}.yaml", "w+") as f:
        f.write(agent_config_file)
    return dir_path
