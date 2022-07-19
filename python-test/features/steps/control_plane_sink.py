from behave import given, when, then, step
from test_config import TestConfig
from utils import random_string, filter_list_by_parameter_start_with, threading_wait_until, validate_json
from hamcrest import *
import requests

configs = TestConfig.configs()
sink_label_name_prefix = "test_sink_label_name_"
orb_url = configs.get('orb_url')


@given("that the user has the prometheus/grafana credentials")
def check_prometheus_grafana_credentials(context):
    context.remote_prometheus_endpoint = configs.get('remote_prometheus_endpoint')
    assert_that(context.remote_prometheus_endpoint, not_none(), 'No remote write endpoint to send Prometheus metrics '
                                                                'to Grafana Cloud was provided!')
    assert_that(context.remote_prometheus_endpoint, not_(""), 'No remote write endpoint to send Prometheus metrics to '
                                                              'Grafana Cloud was provided!')

    context.prometheus_username = configs.get('prometheus_username')
    assert_that(context.prometheus_username, not_none(), 'No Grafana Cloud Prometheus username was provided!')
    assert_that(context.prometheus_username, not_(""), 'No Grafana Cloud Prometheus username was provided!')

    context.prometheus_key = configs.get('prometheus_key')
    assert_that(context.prometheus_key, not_none(), 'No Grafana Cloud API Key was provided!')
    assert_that(context.prometheus_key, not_(""), 'No Grafana Cloud API Key was provided!')


@when("a new sink is created")
def create_sink(context):
    sink_label_name = sink_label_name_prefix + random_string(10)
    token = context.token
    endpoint = context.remote_prometheus_endpoint
    username = context.prometheus_username
    password = context.prometheus_key
    context.sink = create_new_sink(token, sink_label_name, endpoint, username, password)
    local_orb_path = configs.get("local_orb_path")
    sink_schema_path = local_orb_path + "/python-test/features/steps/schemas/sink_schema.json"
    is_schema_valid = validate_json(context.sink, sink_schema_path)
    assert_that(is_schema_valid, equal_to(True), f"Invalid sink json. \n Sink = {context.sink}")
    context.existent_sinks_id.append(context.sink['id'])


@step("{amount_of_sinks} new sinks are created")
def create_multiple_sinks(context, amount_of_sinks):
    check_prometheus_grafana_credentials(context)
    for sink in range(int(amount_of_sinks)):
        create_sink(context)
        
        
@given("that a sink already exists")
def new_sink(context):
    check_prometheus_grafana_credentials(context)
    create_sink(context)
    
    
@step("that {amount_of_sinks} sinks already exists")
def new_multiple_sinks(context, amount_of_sinks):
    check_prometheus_grafana_credentials(context)
    for sink in range(int(amount_of_sinks)):
        create_sink(context)


@step("remove {amount_of_sinks} of the linked sinks from orb")
def remove_sink_from_orb(context, amount_of_sinks):
    for i in range(int(amount_of_sinks)):
        delete_sink(context.token, context.used_sinks_id[i])
        if 'removed_sinks_ids' in context:
            context.removed_sinks_ids.append(context.used_sinks_id[i])
        else:
            context.removed_sinks_ids = list()
            context.removed_sinks_ids.append(context.used_sinks_id[i])
        context.existent_sinks_id.remove(context.used_sinks_id[i])
        context.used_sinks_id.remove(context.used_sinks_id[i])


@step("that a sink with invalid {credential} already exists")
def create_invalid_sink(context, credential):
    assert_that(credential, any_of(equal_to('endpoint'), equal_to('username'), equal_to('password')),
                "Invalid prometheus field")
    check_prometheus_grafana_credentials(context)
    sink_label_name = sink_label_name_prefix + random_string(10)
    token = context.token
    prometheus_credentials = {'endpoint': context.remote_prometheus_endpoint, 'username': context.prometheus_username,
                              'password': context.prometheus_key}
    prometheus_credentials[credential] = prometheus_credentials[credential][:-2]
    context.sink = create_new_sink(token, sink_label_name, prometheus_credentials['endpoint'],
                                   prometheus_credentials['username'], prometheus_credentials['password'])
    local_orb_path = configs.get("local_orb_path")
    sink_schema_path = local_orb_path + "/python-test/features/steps/schemas/sink_schema.json"
    is_schema_valid = validate_json(context.sink, sink_schema_path)
    assert_that(is_schema_valid, equal_to(True), f"Invalid sink json. \n Sink = {context.sink}")
    context.existent_sinks_id.append(context.sink['id'])


@step("referred sink must have {status} state on response within {time_to_wait} seconds")
def check_sink_status(context, status, time_to_wait):
    sink_id = context.sink["id"]
    get_sink_response = get_sink_status_and_check(context.token, sink_id, status, timeout=time_to_wait)

    assert_that(get_sink_response['state'], equal_to(status), f"Sink {context.sink} state failed")


@then('cleanup sinks')
def clean_sinks(context):
    """
    Remove all sinks starting with 'sink_label_name_prefix' from the orb

    :param context: Behave class that contains contextual information during the running of tests.
    """
    token = context.token
    sinks_list = list_sinks(token)
    sinks_filtered_list = filter_list_by_parameter_start_with(sinks_list, 'name', sink_label_name_prefix)
    delete_sinks(token, sinks_filtered_list)


def create_new_sink(token, name_label, remote_host, username, password, description=None, tag_key='',
                    tag_value=None, backend_type="prometheus"):
    """

    Creates a new sink in Orb control plane

    :param (str) token: used for API authentication
    :param (str) name_label:  of the sink to be created
    :param (str) remote_host: base url to send metrics to a dashboard
    :param (str) username: user that enables access to the dashboard
    :param (str) password: key that enables access to the dashboard
    :param (str) description: description of sink
    :param (str) tag_key: the key of the tag to be added to this sink. Default: ''
    :param (str) tag_value: the value of the tag to be added to this sink. Default: None
    :param (str) backend_type: type of backend used to send metrics. Default: prometheus
    :return: (dict) a dictionary containing the created sink data
    """
    json_request = {"name": name_label, "description": description, "tags": {tag_key: tag_value},
                    "backend": backend_type, "validate_only": False,
                    "config": {"remote_host": remote_host, "username": username, "password": password}}
    headers_request = {'Content-type': 'application/json', 'Accept': '*/*',
                       'Authorization': f'Bearer {token}'}

    response = requests.post(orb_url + '/api/v1/sinks', json=json_request, headers=headers_request)
    assert_that(response.status_code, equal_to(201),
                'Request to create sink failed with status=' + str(response.status_code) + ': ' + str(response.json()))

    return response.json()


def get_sink(token, sink_id):
    """
    Gets a sink from Orb control plane

    :param (str) token: used for API authentication
    :param (str) sink_id: that identifies sink to be fetched
    :returns: (dict) the fetched sink
    """

    get_sink_response = requests.get(orb_url + '/api/v1/sinks/' + sink_id, headers={'Authorization': f'Bearer {token}'})

    assert_that(get_sink_response.status_code, equal_to(200),
                'Request to get sink id=' + sink_id + ' failed with status=' + str(get_sink_response.status_code) + ': '
                + str(get_sink_response.json()))

    return get_sink_response.json()


def list_sinks(token, limit=100, offset=0):
    """
    Lists all sinks from Orb control plane that belong to this user

    :param (str) token: used for API authentication
    :param (int) limit: Size of the subset to retrieve. (max 100). Default = 100
    :param (int) offset: Number of items to skip during retrieval. Default = 0.
    :returns: (list) a list of sinks
    """
    all_sinks, total, offset = list_up_to_limit_sinks(token, limit, offset)

    new_offset = limit + offset

    while new_offset < total:
        sinks_from_offset, total, offset = list_up_to_limit_sinks(token, limit, new_offset)
        all_sinks = all_sinks + sinks_from_offset
        new_offset = limit + offset

    return all_sinks


def list_up_to_limit_sinks(token, limit=100, offset=0):
    """
    Lists up to 100 sinks from Orb control plane that belong to this user

    :param (str) token: used for API authentication
    :param (int) limit: Size of the subset to retrieve. (max 100). Default = 100
    :param (int) offset: Number of items to skip during retrieval. Default = 0.
    :returns: (list) a list of sinks, (int) total sinks on orb, (int) offset
    """

    response = requests.get(orb_url + '/api/v1/sinks', headers={'Authorization': f'Bearer {token}'},
                            params={'limit': limit, 'offset': offset})

    assert_that(response.status_code, equal_to(200),
                'Request to list sinks failed with status=' + str(response.status_code) + ': '
                + str(response.json()))

    sinks_as_json = response.json()
    return sinks_as_json['sinks'], sinks_as_json['total'], sinks_as_json['offset']


def delete_sinks(token, list_of_sinks):
    """
    Deletes from Orb control plane the sinks specified on the given list

    :param (str) token: used for API authentication
    :param (list) list_of_sinks: that will be deleted
    """

    for sink in list_of_sinks:
        delete_sink(token, sink['id'])


def delete_sink(token, sink_id):
    """
    Deletes a sink from Orb control plane

    :param (str) token: used for API authentication
    :param (str) sink_id: that identifies the sink to be deleted
    """

    response = requests.delete(orb_url + '/api/v1/sinks/' + sink_id,
                               headers={'Authorization': f'Bearer {token}'})

    assert_that(response.status_code, equal_to(204), 'Request to delete sink id='
                + sink_id + ' failed with status=' + str(response.status_code))


@threading_wait_until
def get_sink_status_and_check(token, sink_id, status, event=None):
    """

    :param (str) token: used for API authentication
    :param (str) sink_id: that identifies the sink to be deleted
    :param status: expected status for referred sink
    :param (obj) event: threading.event
    :returns: (dict) data of the fetched sink
    """
    get_sink_response = get_sink(token, sink_id)
    if get_sink_response['state'] == status:
        event.set()
        return get_sink_response
    return get_sink_response
