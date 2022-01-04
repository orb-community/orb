from behave import given, when, then
from utils import random_string, filter_list_by_parameter_start_with
from hamcrest import *
import requests
from test_config import TestConfig
from datetime import datetime

dataset_name_prefix = "test_dataset_name_"

base_orb_url = TestConfig.configs().get('base_orb_url')


@when("a new dataset is created using referred group, sink and policy ID")
def create_new_dataset(context):
    context.dataset_applied_timestamp = datetime.now().timestamp()
    token = context.token
    agent_groups_id = context.agent_group_data['id']
    sink_id = context.sink['id']
    policy_id = context.policy['id']
    dataset_name = dataset_name_prefix + random_string(10)
    context.dataset = create_dataset(token, dataset_name, policy_id, agent_groups_id, sink_id)


def create_dataset(token, name_label, policy_id, agent_group_id, sink_id):
    """

    :param (str) token: used for API authentication
    :param (str) name_label:  of the dataset to be created
    :param (str) policy_id: that identifies policy to be bound
    :param (str) agent_group_id: that identifies agent_group to be bound
    :param (str) sink_id: that identifies sink to be bound
    :return:
    """

    json_request = {"name": name_label, "agent_group_id": agent_group_id, "agent_policy_id": policy_id,
                    "sink_ids": [sink_id]}
    header_request = {'Content-type': 'application/json', 'Accept': '*/*', 'Authorization': token}

    response = requests.post(base_orb_url + '/api/v1/policies/dataset', json=json_request, headers=header_request)
    assert_that(response.status_code, equal_to(201),
                'Request to create dataset failed with status=' + str(response.status_code))

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


@given('that a dataset using referred group, sink and policy ID already exists')
def new_dataset(context):
    create_new_dataset(context)


def list_datasets(token, limit=100):
    """
    Lists up to 100 datasets from Orb control plane that belong to this user

    :param (str) token: used for API authentication
    :param (int) limit: Size of the subset to retrieve. (max 100). Default = 100
    :returns: (list) a list of datasets
    """

    response = requests.get(base_orb_url + '/api/v1/policies/dataset', headers={'Authorization': token},
                            params={"limit": limit})

    assert_that(response.status_code, equal_to(200),
                'Request to list datasets failed with status=' + str(response.status_code))

    datasets_as_json = response.json()
    return datasets_as_json['datasets']


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

    response = requests.delete(base_orb_url + '/api/v1/policies/dataset/' + dataset_id,
                               headers={'Authorization': token})

    assert_that(response.status_code, equal_to(204), 'Request to delete dataset id='
                + dataset_id + ' failed with status=' + str(response.status_code))
