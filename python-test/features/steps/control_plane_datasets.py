import random

from behave import then, step
from utils import random_string, filter_list_by_parameter_start_with, validate_json, threading_wait_until
from hamcrest import *
import requests
from test_config import TestConfig
from datetime import datetime
from random import choice

dataset_name_prefix = "test_dataset_name_"

orb_url = TestConfig.configs().get('orb_url')
configs = TestConfig.configs()
verify_ssl_bool = eval(configs.get('verify_ssl').title())


@step("a new dataset is requested to be created with the same name as an existent one")
def create_dataset_with_conflict_name(context):
    policy_id = context.policy['id']
    dataset_name = context.dataset['name']
    groups_to_be_used = random.sample(list(context.agent_groups.keys()), int(1))
    context.error_message = create_dataset(context.token, dataset_name, policy_id, groups_to_be_used[-1],
                                           context.used_sinks_id, 409)


@step("{amount_of_datasets} new dataset is created using the policy, {group_order} group and {amount_of_sinks}"
      " {sink_number}")
def create_new_dataset(context, amount_of_datasets, group_order, amount_of_sinks, sink_number):
    assert_that(sink_number, any_of(equal_to("sink"), equal_to("sinks")), "Unexpected value for sink")
    assert_that(group_order, any_of(equal_to("first"), equal_to("second"), equal_to("last"), equal_to("an existing")),
                "Unexpected value for group.")

    if group_order == "an existing":
        groups_to_be_used = random.sample(list(context.agent_groups.keys()), int(amount_of_datasets))
    else:
        assert_that(str(amount_of_datasets), equal_to(str(1)), "For more than one dataset, pass 'an existing' as group"
                                                               " parameter")
        order_convert = {"first": 0, "last": -1, "second": 1}
        groups_to_be_used = [list(context.agent_groups.keys())[order_convert[group_order]]]

    for i in range(int(amount_of_datasets)):
        context.considered_timestamp = datetime.now().timestamp()
        token = context.token
        if amount_of_sinks == 1:
            context.used_sinks_id = [context.sink['id']]
        else:
            # todo create scenario with multiple sinks
            context.used_sinks_id = context.existent_sinks_id[:int(amount_of_sinks)]
        policy_id = context.policy['id']
        dataset_name = dataset_name_prefix + random_string(10)
        context.dataset = create_dataset(token, dataset_name, policy_id, groups_to_be_used[i], context.used_sinks_id)
        local_orb_path = configs.get("local_orb_path")
        dataset_schema_path = local_orb_path + "/python-test/features/steps/schemas/dataset_schema.json"
        is_schema_valid = validate_json(context.dataset, dataset_schema_path)
        assert_that(is_schema_valid, equal_to(True), f"Invalid dataset json. \n Dataset = {context.dataset}. \n"
                                                     f"Policy: {context.policy}. \n Group: {groups_to_be_used[i]}. \n"
                                                     f"Sink(s): {context.used_sinks_id}")
        if 'datasets_created' in context:
            context.datasets_created[context.dataset['id']] = context.dataset['name']
        else:
            context.datasets_created = dict()
            context.datasets_created[context.dataset['id']] = context.dataset['name']


@step("the dataset is edited and {amount_of_sinks} sinks are linked")
def edit_sinks_on_dataset(context, amount_of_sinks):
    dataset = get_dataset(context.token, context.dataset['id'])
    sinks = context.existent_sinks_id[:int(amount_of_sinks)]
    edit_dataset(context.token, dataset['id'], dataset['name'], dataset['agent_policy_id'], dataset['agent_group_id'],
                 sinks)


@step("editing a dataset using a name already in use")
def edit_dataset_with_conflict_name(context):
    dataset = get_dataset(context.token, context.dataset['id'])
    sinks = context.existent_sinks_id[:int(1)]
    datasets_list = list_datasets(context.token)
    agents_filtered_list = filter_list_by_parameter_start_with(datasets_list, 'name', dataset_name_prefix)
    datasets_name = list()
    for dataset in agents_filtered_list:
        datasets_name.append(dataset['name'])
    datasets_name.remove(context.dataset['name'])
    name_to_use = random.choice(datasets_name)
    context.error_message = edit_dataset(context.token, dataset['id'], name_to_use, dataset['agent_policy_id'],
                                         dataset['agent_group_id'], sinks, 409)


@step('a dataset linked to this agent is removed')
def remove_dataset_from_agent(context):
    dataset_remove = choice(list(context.datasets_created.keys()))
    context.dataset = get_dataset(context.token, dataset_remove)
    delete_dataset(context.token, dataset_remove)
    context.datasets_created.pop(dataset_remove)
    context.policy.clear()
    context.policy = \
        {'id': context.dataset["agent_policy_id"], 'name': context.policies_created[context.dataset["agent_policy_id"]]}
    context.list_agent_policies_id.remove(context.policy["id"])
    context.policies_created.pop(context.policy["id"])


@step('referred dataset {condition} be listed on the orb datasets list')
def check_orb_datasets_list(context, condition='must'):
    dataset_id = context.dataset['id']
    all_existing_datasets = list_datasets(context.token)
    is_dataset_listed = any(dataset_id in dataset.values() for dataset in all_existing_datasets)
    if condition == 'must':
        assert_that(is_dataset_listed, equal_to(True), f"Dataset {dataset_id} not listed on orb datasets list."
                                                       f" {context.dataset}")
        get_dataset(context.token, dataset_id)
    elif condition == 'must not':
        assert_that(is_dataset_listed, equal_to(False), f"Dataset {dataset_id} exists in the orb datasets list. "
                                                        f"{context.dataset}")
        policy = get_dataset(context.token, dataset_id, 404)
        assert_that(policy['error'], equal_to('non-existent entity'),
                    f"Unexpected response for get dataset request. {policy}")


@step('no dataset should be linked to the removed {element_removed} anymore')
def check_dataset_status_valid(context, element_removed):
    assert_that(element_removed, any_of(equal_to('group'), equal_to('groups'), equal_to("sink"), equal_to("sinks"),
                                        equal_to("policy"), equal_to("policies")), "Unexpected removed element.")
    datasets_list = list_datasets(context.token)
    datasets_test_list = [dataset for dataset in datasets_list if dataset['id'] in context.datasets_created]
    if element_removed == "group" or element_removed == "groups":
        related_datasets = [dataset for dataset in datasets_test_list if
                            dataset['agent_group_id'] in context.removed_groups_ids]
    elif element_removed == "policy" or element_removed == "policies":
        related_datasets = [dataset for dataset in datasets_test_list if
                            dataset['agent_policy_id'] in context.removed_policies_ids]
    else:
        related_datasets = [dataset for dataset in datasets_test_list if
                            any(sink in dataset['sink_ids'] for sink in context.removed_sinks_ids)]
    assert_that(len(related_datasets), equal_to(0), f"The following datasets are still linked to removed"
                                                    f" {element_removed}: {related_datasets}")


@step("{amount_of_datasets_valid} dataset(s) have validity valid and {amount_of_datasets_invalid} have validity "
      "invalid in {time_to_wait} seconds")
def check_amount_datasets_per_status(context, amount_of_datasets_valid, amount_of_datasets_invalid, time_to_wait):
    amount_of_datasets_valid = int(amount_of_datasets_valid)
    amount_of_datasets_invalid = int(amount_of_datasets_invalid)
    valid_datasets, invalid_datasets = check_all_test_dataset_per_status(context.token, context.datasets_created,
                                                                         amount_of_datasets_valid,
                                                                         amount_of_datasets_invalid,
                                                                         timeout=int(time_to_wait))
    assert_that(len(valid_datasets), equal_to(amount_of_datasets_valid),
                f"Unexpected amount of datasets valid.\nValid: {valid_datasets}. \nInvalid: {invalid_datasets}")
    assert_that(len(invalid_datasets), equal_to(amount_of_datasets_invalid),
                f"Unexpected amount of datasets invalid.\nValid: {valid_datasets}. \nInvalid: {invalid_datasets}")


@threading_wait_until
def check_all_test_dataset_per_status(token, existing_datasets_ids_list, amount_of_datasets_valid,
                                      amount_of_datasets_invalid, event=None):
    datasets_list = list_datasets(token)
    datasets_test_list = [dataset for dataset in datasets_list if dataset['id'] in existing_datasets_ids_list]
    valid_datasets = [dataset for dataset in datasets_test_list if dataset['valid'] is True]
    invalid_datasets = [dataset for dataset in datasets_test_list if dataset['valid'] is False]
    if len(valid_datasets) == amount_of_datasets_valid and len(invalid_datasets) == amount_of_datasets_invalid:
        event.set()
    return valid_datasets, invalid_datasets


def create_dataset(token, name_label, policy_id, agent_group_id, sink_id, expected_status_code=201):
    """

    :param (str) token: used for API authentication
    :param (str) name_label:  of the dataset to be created
    :param (str) policy_id: that identifies policy to be bound
    :param (str) agent_group_id: that identifies agent_group to be bound
    :param (str) sink_id: that identifies sink to be bound
    :param (int) expected_status_code: code to be on response
    :return:
    """

    json_request = {"name": name_label, "agent_group_id": agent_group_id, "agent_policy_id": policy_id,
                    "sink_ids": sink_id}
    header_request = {'Content-type': 'application/json', 'Accept': '*/*', 'Authorization': f'Bearer {token}'}

    response = requests.post(orb_url + '/api/v1/policies/dataset', json=json_request, headers=header_request,
                             verify=verify_ssl_bool)
    try:
        response_json = response.json()
    except ValueError:
        response_json = ValueError
    assert_that(response.status_code, equal_to(expected_status_code),
                'Request to create dataset failed with status=' + str(response.status_code) + ': ' +
                str(response_json))

    return response_json


def edit_dataset(token, dataset_id, name_label, policy_id, agent_group_id, sink_id, expected_status_code=200):
    """

    :param (str) dataset_id: that identifies dataset to be edited
    :param (str) token: used for API authentication
    :param (str) name_label:  of the dataset to be created
    :param (str) policy_id: that identifies policy to be bound
    :param (str) agent_group_id: that identifies agent_group to be bound
    :param (str) sink_id: that identifies sink to be bound
    :param (int) expected_status_code: code to be on response
    :return: edited dataset json
    """

    json_request = {"name": name_label, "agent_group_id": agent_group_id, "agent_policy_id": policy_id,
                    "sink_ids": sink_id}
    header_request = {'Content-type': 'application/json', 'Accept': '*/*', 'Authorization': f'Bearer {token}'}

    response = requests.put(f"{orb_url}/api/v1/policies/dataset/{dataset_id}", json=json_request,
                            headers=header_request, verify=verify_ssl_bool)
    assert_that(response.status_code, equal_to(expected_status_code),
                'Request to edit dataset failed with status=' + str(response.status_code) + ': ' + str(response.json()))

    return response.json()


@then('cleanup datasets')
def clean_datasets(context):
    """
    Remove all datasets starting with 'test_dataset_name_' from orb

    :param context: Behave class that contains contextual information during the running of tests.
    """
    token = context.token
    datasets_list = list_datasets(token)
    datasets_filtered_list = filter_list_by_parameter_start_with(datasets_list, 'name', dataset_name_prefix)
    delete_datasets(token, datasets_filtered_list)


def list_datasets(token, limit=100, offset=0):
    """
    Lists all datasets from Orb control plane that belong to this user

    :param (str) token: used for API authentication
    :param (int) limit: Size of the subset to retrieve. (max 100). Default = 100
    :param (int) offset: Number of items to skip during retrieval. Default = 0.
    :returns: (list) a list of datasets
    """

    all_datasets, total, offset = list_up_to_limit_datasets(token, limit, offset)

    new_offset = limit + offset

    while new_offset < total:
        datasets_from_offset, total, offset = list_up_to_limit_datasets(token, limit, new_offset)
        all_datasets = all_datasets + datasets_from_offset
        new_offset = limit + offset

    return all_datasets


def list_up_to_limit_datasets(token, limit=100, offset=0):
    """
    Lists up to 100 datasets from Orb control plane that belong to this user

    :param (str) token: used for API authentication
    :param (int) limit: Size of the subset to retrieve. (max 100). Default = 100
    :returns: (list) a list of datasets, (int) total datasets on orb, (int) offset
    """

    response = requests.get(orb_url + '/api/v1/policies/dataset', headers={'Authorization': f'Bearer {token}'},
                            params={"limit": limit, "offset": offset}, verify=verify_ssl_bool)

    assert_that(response.status_code, equal_to(200),
                'Request to list datasets failed with status=' + str(response.status_code) + ':' + str(response.json()))

    datasets_as_json = response.json()
    return datasets_as_json['datasets'], datasets_as_json['total'], datasets_as_json['offset']


def delete_datasets(token, list_of_datasets):
    """
    Deletes from Orb control plane the datasets specified on the given list

    :param (str) token: used for API authentication
    :param (list) list_of_datasets: that will be deleted
    """

    for dataset in list_of_datasets:
        delete_dataset(token, dataset['id'])


def delete_dataset(token, dataset_id):
    """
    Deletes a dataset from Orb control plane

    :param (str) token: used for API authentication
    :param (str) dataset_id: that identifies the dataset to be deleted
    """

    response = requests.delete(orb_url + '/api/v1/policies/dataset/' + dataset_id,
                               headers={'Authorization': f'Bearer {token}'}, verify=verify_ssl_bool)

    assert_that(response.status_code, equal_to(204), 'Request to delete dataset id='
                + dataset_id + ' failed with status=' + str(response.status_code))


def get_dataset(token, dataset_id, expected_status_code=200):
    """
    Gets a dataset from Orb control plane

    :param (str) token: used for API authentication
    :param (str) dataset_id: that identifies dataset to be fetched
    :param (int) expected_status_code: expected request's status code. Default:200.
    :returns: (dict) the fetched dataset
    """

    get_dataset_response = requests.get(orb_url + '/api/v1/policies/dataset/' + dataset_id,
                                        headers={'Authorization': f'Bearer {token}'}, verify=verify_ssl_bool)

    try:
        response_json = get_dataset_response.json()
    except ValueError:
        response_json = ValueError

    assert_that(get_dataset_response.status_code, equal_to(expected_status_code),
                'Request to get policy id=' + dataset_id + ' failed with status=' +
                str(get_dataset_response.status_code) + "response=" + str(response_json))

    return response_json
