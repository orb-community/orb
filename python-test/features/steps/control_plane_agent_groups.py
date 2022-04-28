from test_config import TestConfig
from local_agent import get_orb_agent_logs
from users import get_auth_token
from utils import random_string, filter_list_by_parameter_start_with, generate_random_string_with_predefined_prefix, \
    create_tags_set, check_logs_contain_message_and_name, threading_wait_until
from behave import given, then, step
from hamcrest import *
import requests
from random import sample

configs = TestConfig.configs()
agent_group_name_prefix = 'test_group_name_'
agent_group_description = "This is an agent group"
orb_url = configs.get('orb_url')


@step("an Agent Group is created with {amount_of_tags} tags contained in the agent")
def create_agent_group_matching_agent(context, amount_of_tags, **kwargs):
    if amount_of_tags.isdigit() is False:
        assert_that(amount_of_tags, equal_to("all"), 'Unexpected value for amount of tags')
    agent_group_name = agent_group_name_prefix + random_string()
    if "group_description" in kwargs.keys():
        group_description = kwargs["group_description"]
    else:
        group_description = agent_group_description

    tags_in_agent = context.agent["orb_tags"]
    if context.agent["agent_tags"] is not None:
        tags_in_agent.update(context.agent["agent_tags"])
    tags_keys = tags_in_agent.keys()

    if amount_of_tags.isdigit() is True:
        amount_of_tags = int(amount_of_tags)
    else:
        amount_of_tags = len(tags_keys)
    assert_that(tags_keys, has_length(greater_than_or_equal_to(amount_of_tags)), "Amount of tags greater than tags"
                                                                                      "contained in agent")
    tags_to_group = {key: tags_in_agent[key] for key in sample(tags_keys, amount_of_tags)}
    context.agent_group_data = create_agent_group(context.token, agent_group_name, group_description,
                                                  tags_to_group)
    group_id = context.agent_group_data['id']
    context.agent_groups[group_id] = agent_group_name


@step("an Agent Group is created with {orb_tags} orb tag(s)")
def create_new_agent_group(context, orb_tags, **kwargs):
    agent_group_name = generate_random_string_with_predefined_prefix(agent_group_name_prefix)
    if "group_description" in kwargs.keys():
        group_description = kwargs["group_description"]
    else:
        group_description = agent_group_description
    context.orb_tags = create_tags_set(orb_tags)
    if len(context.orb_tags) == 0:
        context.agent_group_data = create_agent_group(context.token, agent_group_name, group_description,
                                                      context.orb_tags, 400)
    else:
        context.agent_group_data = create_agent_group(context.token, agent_group_name, group_description,
                                                      context.orb_tags)
        group_id = context.agent_group_data['id']
        context.agent_groups[group_id] = agent_group_name


@step("an Agent Group is created with {orb_tags} orb tag(s) and {description} description")
def create_new_agent_group_with_defined_description(context, orb_tags, description):
    if description == "without":
        create_new_agent_group(context, orb_tags, group_description=None)
    else:
        description = description.replace('"', '')
        description = description.replace(' as', '')
        create_new_agent_group(context, orb_tags, group_description=description)


@step("an Agent Group is created with same tag as the agent and {description} description")
def create_agent_group_with_defined_description_and_matching_agent(context, description):
    if description == "without":
        create_agent_group_matching_agent(context, group_description=None)
    else:
        description = description.replace('"', '')
        description = description.replace(' as', '')
        create_agent_group_matching_agent(context, group_description=description)


@step("the {edited_parameters} of Agent Group is edited using: {parameters_values}")
def edit_multiple_groups_parameters(context, edited_parameters, parameters_values):
    edited_parameters = edited_parameters.split(", ")
    for param in edited_parameters:
        assert_that(param, any_of(equal_to('name'), equal_to('description'), equal_to('tags')),
                    'Unexpected parameter to edit')
    parameters_values = parameters_values.split("/ ")

    group_editing = get_agent_group(context.token, context.agent_group_data["id"])
    group_data = {"name": group_editing["name"], "tags": group_editing["tags"]}
    if "description" in group_editing.keys():
        group_data["description"] = group_editing["description"]
    else:
        group_data["description"] = None

    editing_param_dict = dict()
    for param in parameters_values:
        param_split = param.split("=")
        if param_split[1].lower() == "none":
            param_split[1] = None
        editing_param_dict[param_split[0]] = param_split[1]

    assert_that(set(editing_param_dict.keys()), equal_to(set(edited_parameters)),
                "All parameter must have referenced value")

    if "tags" in editing_param_dict.keys() and editing_param_dict["tags"] is not None:
        editing_param_dict["tags"] = create_tags_set(editing_param_dict["tags"])
    if "name" in editing_param_dict.keys() and editing_param_dict["name"] is not None:
        editing_param_dict["name"] = agent_group_name_prefix + editing_param_dict["name"]

    for parameter, value in editing_param_dict.items():
        group_data[parameter] = value

    context.editing_response = edit_agent_group(context.token, context.agent_group_data["id"], group_data["name"],
                                                group_data["description"], group_data["tags"])


@then("agent group editing must fail")
def fail_group_editing(context):
    assert_that(list(context.editing_response.keys())[0], equal_to("error"))


@step("Agent Group creation response must be an error with message '{message}'")
def error_response_message(context, message):
    response = list(context.agent_group_data.items())[0]
    response_key, response_value = response[0], response[1]
    assert_that(response_key, equal_to('error'),
                'Response of invalid agent group creation must be an error')
    assert_that(response_value, equal_to(message), "Unexpected message for error")


@step("{amount_agent_matching} agent must be matching on response field matching_agents")
def matching_agent(context, amount_agent_matching):
    context.agent_group_data = get_agent_group(context.token, context.agent_group_data["id"])
    matching_total_agents = context.agent_group_data['matching_agents']['total']
    assert_that(matching_total_agents, equal_to(int(amount_agent_matching)))


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


@step('the container logs contain the message "{text_to_match}" referred to each matching group within'
      '{time_to_wait} seconds')
def check_logs_for_group(context, text_to_match, time_to_wait):
    groups_matching, context.groups_matching_id = return_matching_groups(context.token, context.agent_groups, context.agent)
    text_found, groups_to_which_subscribed = check_subscription(groups_matching, text_to_match, context.container_id,
                                                                timeout=time_to_wait)
    assert_that(text_found, is_(True), f"Message {text_to_match} was not found in the agent logs for group(s)"
                                       f"{set(groups_matching).difference(groups_to_which_subscribed)}!")


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

    response = requests.post(orb_url + '/api/v1/agent_groups', json=json_request, headers=headers_request)
    assert_that(response.status_code, equal_to(expected_status_code),
                'Request to create agent group failed with status=' + str(response.status_code))

    return response.json()


def get_agent_group(token, agent_group_id):
    """
    Gets an agent group from Orb control plane

    :param (str) token: used for API authentication
    :param (str) agent_group_id: that identifies the agent group to be fetched
    :returns: (dict) the fetched agent group
    """

    get_groups_response = requests.get(orb_url + '/api/v1/agent_groups/' + agent_group_id,
                                       headers={'Authorization': token})

    assert_that(get_groups_response.status_code, equal_to(200),
                'Request to get agent group id=' + agent_group_id + ' failed with status=' + str(
                    get_groups_response.status_code))

    return get_groups_response.json()


def list_agent_groups(token, limit=100):
    """
    Lists up to 100 agent groups from Orb control plane that belong to this user

    :param (str) token: used for API authentication
    :param (int) limit: Size of the subset to retrieve (max 100). Default = 100
    :returns: (list) a list of agent groups
    """

    response = requests.get(orb_url + '/api/v1/agent_groups', headers={'Authorization': token},
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

    response = requests.delete(orb_url + '/api/v1/agent_groups/' + agent_group_id,
                               headers={'Authorization': token})

    assert_that(response.status_code, equal_to(204), 'Request to delete agent group id='
                + agent_group_id + ' failed with status=' + str(response.status_code))


@threading_wait_until
def check_subscription(agent_groups_names, expected_message, container_id, event=None):
    """

    :param (list) agent_groups_names: groups to which the agent must be subscribed
    :param (str) expected_message: message that we expect to find in the logs
    :param (str) container_id: agent container id
    :param (obj) event: threading.event
    :return: (bool) True if agent is subscribed to all matching groups, (list) names of the groups to which agent is subscribed
    """
    groups_to_which_subscribed = set()
    for name in agent_groups_names:
        logs = get_orb_agent_logs(container_id)
        text_found, log_line = check_logs_contain_message_and_name(logs, expected_message, name, "group_name")
        if text_found is True:
            groups_to_which_subscribed.add(log_line["group_name"])
            if set(groups_to_which_subscribed) == set(agent_groups_names):
                event.set()
                return event.is_set(), groups_to_which_subscribed

    return event.is_set(), groups_to_which_subscribed


def edit_agent_group(token, agent_group_id, name, description, tags, expected_status_code=200):
    """

    :param (str) token: used for API authentication
    :param (str) agent_group_id: that identifies the agent group to be edited
    :param (str) name: agent group's name
    :param (str) description: agent group's description
    :param (str) tags: orb tags that will be used to connect agents to groups
    :param (int) expected_status_code: expected request's status code. Default:200.
    :returns: (dict) the edited agent group
    """

    json_request = {"name": name, "description": description, "tags": tags,
                    "validate_only": False}
    json_request = {parameter: value for parameter, value in json_request.items() if value}

    headers_request = {'Content-type': 'application/json', 'Accept': '*/*', 'Authorization': token}

    group_edited_response = requests.put(orb_url + '/api/v1/agent_groups/' + agent_group_id, json=json_request,
                                         headers=headers_request)

    if name is None or tags is None:
        expected_status_code = 400
    assert_that(group_edited_response.status_code, equal_to(expected_status_code),
                'Request to edit agent group failed with status=' + str(group_edited_response.status_code))

    return group_edited_response.json()


def return_matching_groups(token, existing_agent_groups, agent_json):
    """

    :param (str) token: used for API authentication
    :param (dict) existing_agent_groups: dictionary with the existing groups, the id of the groups being the key and the name the values
    :param (dict) agent_json: dictionary containing all the information of the agent to which the groups must be matching

    :return (list): groups_matching, groups_matching_id
    """
    groups_matching = list()
    groups_matching_id = list()
    for group in existing_agent_groups.keys():
        group_data = get_agent_group(token, group)
        group_tags = dict(group_data["tags"])
        agent_tags = agent_json["orb_tags"]
        if all(item in agent_tags.items() for item in group_tags.items()) is True:
            groups_matching.append(existing_agent_groups[group])
            groups_matching_id.append(group)
    return groups_matching, groups_matching_id
