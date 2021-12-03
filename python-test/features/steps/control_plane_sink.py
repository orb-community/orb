from behave import given, when, then
from test_config import TestConfig
from control_plane_agents import base_orb_url
from utils import random_string, filter_list_by_parameter_start_with
from hamcrest import *
import requests

configs = TestConfig.configs()
sink_label_name_prefix = "test_sink_label_name_"
sink_label_name = sink_label_name_prefix + random_string(10)


@given("that the user has the prometheus/grafana credentials")
def check_prometheus_grafana_credentials(context):
    remote_prometheus_endpoint = configs.get('remote_prometheus_endpoint')
    assert_that(remote_prometheus_endpoint, not_none(), 'No remote write endpoint to send Prometheus metrics to '
                                                        'Grafana Cloud was provided!')
    assert_that(remote_prometheus_endpoint, not_(""), 'No remote write endpoint to send Prometheus metrics to Grafana '
                                                      'Cloud was provided!')
    context.remote_prometheus_endpoint = "https://" + remote_prometheus_endpoint + "/api/prom/push"

    context.prometheus_username = configs.get('prometheus_username')
    assert_that(context.prometheus_username, not_none(), 'No Grafana Cloud Prometheus username was provided!')
    assert_that(context.prometheus_username, not_(""), 'No Grafana Cloud Prometheus username was provided!')

    context.prometheus_key = configs.get('prometheus_key')
    assert_that(context.prometheus_key, not_none(), 'No Grafana Cloud API Key was provided!')
    assert_that(context.prometheus_key, not_(""), 'No Grafana Cloud API Key was provided!')


@when("a new sink is created")
def create_sink(context):
    token = context.token
    endpoint = context.remote_prometheus_endpoint
    username = context.prometheus_username
    password = context.prometheus_key
    context.sink = create_new_sink(token, sink_label_name, endpoint, username, password)


@then("referred sink must have {status} state on response")
def check_sink_status(context, status):
    sink_id = context.sink['id']
    assert_that(sink_id['state'], equal_to(status), f"Sink {sink_id} state failed")
    get_sink_response = get_sink(context.token, sink_id)
    assert_that(sink_id['state'], equal_to(status), f"Sink {sink_id} state failed")


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
    :param name_label:  of the sink to be created
    :param remote_host: base url to send metrics to a dashboard
    :param username: user that enables access to the dashboard
    :param password: key that enables access to the dashboard
    :param (str) description: description of sink
    :param tag_key: the key of the tag to be added to this sink. Default: ''
    :param tag_value: the value of the tag to be added to this sink. Default: None
    :param backend_type: type of backend used to send metrics. Default: prometheus
    :return: (dict) a dictionary containing the created sink data
    """
    response = requests.post(base_orb_url + '/api/v1/sinks',
                             json={"name": name_label, "description": description, "tags": {tag_key: tag_value},
                                   "backend": backend_type, "validate_only": False,
                                   "config": {"remote_host": remote_host, "username": username, "password": password}},
                             headers={'Content-type': 'application/json', 'Accept': '*/*',
                                      'Authorization': token})
    assert_that(response.status_code, equal_to(201),
                'Request to create sink failed with status=' + str(response.status_code))

    return response.json()


def get_sink(token, sink_id):
    """
    Gets a sink from Orb control plane

    :param (str) token: used for API authentication
    :param (str) sink_id: that identifies sink to be fetched
    :returns: (dict) the fetched sink
    """

    get_sink_response = requests.get(base_orb_url + '/api/v1/sinks/' + sink_id, headers={'Authorization': token})

    assert_that(get_sink_response.status_code, equal_to(200),
                'Request to get sink id=' + sink_id + ' failed with status=' + str(get_sink_response.status_code))

    return get_sink_response.json()


def list_sinks(token):
    """
    Lists all sinks from Orb control plane that belong to this user

    :param (str) token: used for API authentication
    :returns: (list) a list of sinks
    """

    response = requests.get(base_orb_url + '/api/v1/sinks', headers={'Authorization': token})

    assert_that(response.status_code, equal_to(200),
                'Request to list sinks failed with status=' + str(response.status_code))

    sinks_as_json = response.json()
    return sinks_as_json['sinks']


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

    response = requests.delete(base_orb_url + '/api/v1/sinks/' + sink_id,
                               headers={'Authorization': token})

    assert_that(response.status_code, equal_to(204), 'Request to delete sink id='
                + sink_id + ' failed with status=' + str(response.status_code))
