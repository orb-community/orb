from behave import given, when, then, step
from configs import TestConfig
from utils import (random_string, filter_list_by_parameter_start_with, threading_wait_until, validate_json, \
                   remove_empty_from_json, return_api_get_response, return_api_put_response, return_api_post_response,
                   return_api_delete_response)
from hamcrest import *
import threading
import yaml

configs = TestConfig.configs()
sink_name_prefix = "test_sink_label_name_"
orb_url = configs.get('orb_url')
verify_ssl_bool = eval(configs.get('verify_ssl').title())


@given("that the user has the prometheus/grafana credentials")
def check_prometheus_grafana_credentials(context):
    if configs.get("sink_backend_type") == "otlphttp":
        context.remote_prometheus_endpoint = configs.get('otlp_publisher_endpoint')
        assert_that(context.remote_prometheus_endpoint, not_none(), 'No remote write endpoint to send otlp '
                                                                    'metrics'
                                                                    'to Grafana Cloud was provided!')
        assert_that(context.remote_prometheus_endpoint, not_(""), 'No remote write endpoint to send otlp '
                                                                  'metrics to'
                                                                  'Grafana Cloud was provided!')

        context.prometheus_username = configs.get('otlp_publisher_username')
        assert_that(context.prometheus_username, not_none(), 'No Grafana Cloud otlp username was provided!')
        assert_that(context.prometheus_username, not_(""), 'No Grafana Cloud otlp username was provided!')

        context.prometheus_key = configs.get('otlp_publisher_key')
        assert_that(context.prometheus_key, not_none(), 'No Grafana Cloud API Key was provided!')
        assert_that(context.prometheus_key, not_(""), 'No Grafana Cloud API Key was provided!')
    else:
        context.remote_prometheus_endpoint = configs.get('remote_prometheus_endpoint')
        assert_that(context.remote_prometheus_endpoint, not_none(), 'No remote write endpoint to send Prometheus '
                                                                    'metrics'
                                                                    'to Grafana Cloud was provided!')
        assert_that(context.remote_prometheus_endpoint, not_(""), 'No remote write endpoint to send Prometheus '
                                                                  'metrics to'
                                                                  'Grafana Cloud was provided!')

        context.prometheus_username = configs.get('prometheus_username')
        assert_that(context.prometheus_username, not_none(), 'No Grafana Cloud Prometheus username was provided!')
        assert_that(context.prometheus_username, not_(""), 'No Grafana Cloud Prometheus username was provided!')

        context.prometheus_key = configs.get('prometheus_key')
        assert_that(context.prometheus_key, not_none(), 'No Grafana Cloud API Key was provided!')
        assert_that(context.prometheus_key, not_(""), 'No Grafana Cloud API Key was provided!')


@step("a new sink is created")
def create_sink(context, **kwargs):
    backend_type = configs.get("sink_backend_type")
    sink_label_name = sink_name_prefix + random_string(10)
    token = context.token
    endpoint = context.remote_prometheus_endpoint
    username = context.prometheus_username
    password = context.prometheus_key
    if "configuration_type" in kwargs:
        configuration_type = kwargs["configuration_type"].lower()
    else:
        configuration_type = ""
    context.sink = create_new_sink(token, sink_label_name, endpoint, username, password, backend_type=backend_type,
                                   configuration_type=configuration_type)
    local_orb_path = configs.get("local_orb_path")
    schema_map = {"prometheus": "/sink_schema.json", "otlphttp": "/sink_otlphttp_schema.json"}
    sink_schema_path = f"{local_orb_path}/python-test/features/steps/schemas{schema_map[backend_type]}"
    is_schema_valid = validate_json(context.sink, sink_schema_path)
    assert_that(is_schema_valid, equal_to(True), f"Invalid sink json. \n Sink = {context.sink}")
    sink_get_response = get_sink(context.token, context.sink['id'])
    is_get_schema_valid = validate_json(sink_get_response, sink_schema_path)
    assert_that(is_get_schema_valid, equal_to(True), f"Invalid sink json in get response."
                                                     f" \n Sink = {sink_get_response}")
    context.existent_sinks_id.append(context.sink['id'])


@step("a new sink is created using {configuration_type} configuration type")
def create_yaml_sink(context, configuration_type):
    create_sink(context, configuration_type=configuration_type)


@step("a new sink is is requested to be created with the same name as an existent one")
def create_sink_with_conflict_name(context):
    token = context.token
    endpoint = context.remote_prometheus_endpoint
    username = context.prometheus_username
    password = context.prometheus_key
    context.error_message = create_new_sink(token, context.sink['name'], endpoint, username, password,
                                            backend_type=backend_type, expected_status_code=409)


@step("the name of last Sink is edited using an already existent one")
def edit_sink_name_with_conflict(context):
    first_sink = get_sink(context.token, context.existent_sinks_id[0])
    context.sink['name'] = first_sink['name']
    context.sink.pop("config")
    context.error_message = edit_sink(context.token, context.sink['id'], context.sink, expected_status_code=409)


@step("{amount_of_sinks} new sinks are created")
def create_multiple_sinks(context, amount_of_sinks):
    check_prometheus_grafana_credentials(context)
    for sink in range(int(amount_of_sinks)):
        create_sink(context)


@given("that a sink with {configuration_type} configuration type already exists")
def new_sink(context, configuration_type):
    check_prometheus_grafana_credentials(context)
    if configuration_type.lower() == "default":
        create_sink(context)
    else:
        assert_that(configuration_type.lower(), any_of(equal_to("yaml"), equal_to("json")),
                    f"Unexpected value for sink configuration type")
        create_sink(context, configuration_type=configuration_type.lower())


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
    sink_label_name = sink_name_prefix + random_string(10)
    token = context.token
    prometheus_credentials = {'endpoint': context.remote_prometheus_endpoint, 'username': context.prometheus_username,
                              'password': context.prometheus_key}
    prometheus_credentials[credential] = prometheus_credentials[credential][:-2]
    context.sink = create_new_sink(token, sink_label_name, prometheus_credentials['endpoint'],
                                   prometheus_credentials['username'], prometheus_credentials['password'],
                                   backend_type=backend_type)
    local_orb_path = configs.get("local_orb_path")
    sink_schema_path = local_orb_path + "/python-test/features/steps/schemas/sink_schema.json"
    is_schema_valid = validate_json(context.sink, sink_schema_path)
    assert_that(is_schema_valid, equal_to(True), f"Invalid sink json. \n Sink = {context.sink}")
    context.existent_sinks_id.append(context.sink['id'])


@step("referred sink must have {status} state on response within {time_to_wait} seconds")
def check_sink_status(context, status, time_to_wait):
    sink_id = context.sink["id"]
    context.sink = get_sink_status_and_check(context.token, sink_id, status, timeout=time_to_wait, wait_time=10)

    assert_that(context.sink['state'], equal_to(status), f"Sink {context.sink} state failed")


@step("referred sink must have {status} state on response after {time_to_wait} seconds")
def check_sink_status(context, status, time_to_wait):
    sink_id = context.sink["id"]
    assert_that(time_to_wait.isdigit(), is_(True), f"Invalid type: 'time_to_wait' must be an int and is {time_to_wait}")
    time_to_wait = int(time_to_wait)
    threading.Event().wait(time_to_wait)
    context.sink = get_sink_status_and_check(context.token, sink_id, status)

    assert_that(context.sink['state'], equal_to(status), f"Sink {context.sink} state failed")


@step("the sink {field_to_edit} is edited and an {type_of_field} one is used")
def edit_sink_field(context, field_to_edit, type_of_field):
    assert_that(field_to_edit, any_of("remote host", "endpoint", "username", "password", "name", "description", "tags"),
                "Unexpected field to be edited.")
    assert_that(type_of_field, any_of("valid", "invalid"), "Invalid type of sink field")
    sink = get_sink(context.token, context.sink['id'])
    sink['config']['authentication']['password'] = configs.get('prometheus_key')
    if field_to_edit == "remote host" or field_to_edit == "endpoint":
        exporter_key = {"prometheus": "remote_host", "otlphttp": "endpoint"}
        sink_backend = context.sink.get("backend")
        assert_that(sink_backend, any_of("prometheus", "otlphttp"), f"Invalid sink backend: '{sink_backend}'")
        if type_of_field == "invalid":
            sink['config']['exporter'][exporter_key[sink_backend]] = context.remote_prometheus_endpoint[:-2]
        else:
            sink['config']['exporter'][exporter_key[sink_backend]] = configs.get('remote_prometheus_endpoint')
    if field_to_edit == "password":
        if type_of_field == "invalid":
            sink['config']['authentication']['password'] = context.prometheus_key[:-2]
        else:
            sink['config']['authentication']['password'] = configs.get('prometheus_key')
    if field_to_edit == "username":
        if type_of_field == "invalid":
            sink['config']['authentication']['username'] = context.prometheus_username[:-2]
        else:
            sink['config']['authentication']['username'] = configs.get('prometheus_username')

    context.sink = edit_sink(context.token, context.sink['id'], sink)


@when("the {sink_keys} of this sink is updated")
def sink_partial_update(context, sink_keys):
    sink_keys = sink_keys.replace(" and ", ", ")
    keys_to_update = sink_keys.split(", ")
    update_sink_body = get_sink(context.token, context.sink['id'])
    keys_to_not_update = list(set(update_sink_body.keys()).difference(set(keys_to_update)))
    all_sink_keys = list(set(update_sink_body.keys()).union(set(keys_to_update)))
    for sink_key in all_sink_keys:
        if sink_key in keys_to_not_update:
            del update_sink_body[sink_key]
        else:
            assert_that(sink_key, any_of("name", "description", "tags", "config"), f"Unexpected key for sink")
            remote_host = context.remote_prometheus_endpoint
            username = context.prometheus_username
            password = context.prometheus_key
            exporter_key = {"prometheus": "remote_host", "otlphttp": "endpoint"}
            sink_backend = context.sink.get("backend")
            assert_that(sink_backend, any_of("prometheus", "otlphttp"), f"Invalid sink backend: '{sink_backend}'")
            sink_configs = {"authentication": {"type": "basicauth", "password": password, "username": username},
                            "exporter": {exporter_key[sink_backend]: remote_host}}
            context.values_to_use_to_update_sink = {"name": f"{context.sink['name']}_updated",
                                                    "description": "this sink has been updated",
                                                    "tags": {"sink": "updated", "new": "tag"}, "config": sink_configs}
            update_sink_body[sink_key] = context.values_to_use_to_update_sink[sink_key]
    context.sink_before_update = get_sink(context.token, context.sink['id'])
    context.sink = edit_sink(context.token, context.sink['id'], update_sink_body, 200)


@then("the {sink_keys} updates to the new value and other fields remains the same")
def verify_sink_after_update(context, sink_keys):
    sink_keys = sink_keys.replace(" and ", ", ")
    updated_keys = sink_keys.split(", ")
    all_sink_keys = list(context.sink.keys())
    assert_that(set(context.sink.keys()), equal_to(set(context.sink_before_update).union(set(updated_keys))),
                f"Sink keys are not the same after sink partial update:"
                f"Sink before update: {context.sink_before_update}. Sink after update: {context.sink}")
    for sink_key in all_sink_keys:
        if sink_key in updated_keys and sink_key != "config":  # config returns empty as a security action
            assert_that(context.sink[sink_key], equal_to(context.values_to_use_to_update_sink[sink_key]),
                        f"{sink_key} was not correctly updated on sink. Sink before update: "
                        f"{context.sink_before_update}. Sink after update: {context.sink}")
        elif sink_key != "ts_created":
            assert_that(context.sink[sink_key], equal_to(context.sink_before_update[sink_key]),
                        f"Unexpected value for {sink_key} after sink partial update."
                        f"Sink before update: {context.sink_before_update}. Sink after update: {context.sink}")
        else:
            pass  # insert validation to ts_created > must be higher than before


@then('cleanup sinks')
def clean_sinks(context):
    """
    Remove all sinks starting with 'sink_name_prefix' from the orb

    :param context: Behave class that contains contextual information during the running of tests.
    """
    token = context.token
    sinks_list = list_sinks(token)
    sinks_filtered_list = filter_list_by_parameter_start_with(sinks_list, 'name', sink_name_prefix)
    delete_sinks(token, sinks_filtered_list)


def create_new_sink(token, name_label, remote_host, username, password, backend_type=None,
                    description=None, tag_key=None, tag_value=None, configuration_type="", expected_status_code=201):
    """

    Creates a new sink in Orb control plane

    :param (str) token: used for API authentication
    :param (str) name_label:  of the sink to be created
    :param (str) remote_host: base url to send metrics to a dashboard
    :param (str) username: user that enables access to the dashboard
    :param (str) password: key that enables access to the dashboard
    :param (str) backend_type: type of backend used to send metrics.
    :param (str) description: description of sink
    :param (str) tag_key: the key of the tag to be added to this sink. Default: ''
    :param (str) tag_value: the value of the tag to be added to this sink. Default: None
    :param (str) configuration_type: define in which type the configuration will be written
    :param (int) expected_status_code: status code expected as response
    :return: (dict) a dictionary containing the created sink data
    """
    assert_that(backend_type, any_of("prometheus", equal_to("otlphttp")), f"Unexpected type of backend.")
    exporter_key = {"prometheus": "remote_host", "otlphttp": "endpoint"}
    sink_configs = {"authentication": {"type": "basicauth", "password": password, "username": username},
                    "exporter": {exporter_key[backend_type]: remote_host}}
    if configuration_type == "yaml":
        sink_configs_yaml = yaml.dump(sink_configs)
        json_request = {"name": name_label, "description": description, "tags": {tag_key: tag_value},
                        "backend": backend_type, "config_data": sink_configs_yaml, "format": configuration_type}
    elif configuration_type == "json":
        json_request = {"name": name_label, "description": description, "tags": {tag_key: tag_value},
                        "backend": backend_type, "config": sink_configs, "format": configuration_type}
    else:
        json_request = {"name": name_label, "description": description, "tags": {tag_key: tag_value},
                        "backend": backend_type, "config": sink_configs}

    json_request = remove_empty_from_json(json_request)

    status_code, response = return_api_post_response(f"{orb_url}/api/v1/sinks", request_body=json_request,
                                                     token=token, verify=verify_ssl_bool)

    assert_that(status_code, equal_to(expected_status_code),
                'Request to create sink failed with status=' + str(status_code) + ': ' + str(response))
    return response


def get_sink(token, sink_id):
    """
    Gets a sink from Orb control plane

    :param (str) token: used for API authentication
    :param (str) sink_id: that identifies sink to be fetched
    :returns: (dict) the fetched sink
    """

    status_code, response = return_api_get_response(f"{orb_url}/api/v1/sinks/{sink_id}", token=token,
                                                    verify=verify_ssl_bool)

    assert_that(status_code, equal_to(200),
                'Request to get sink id=' + sink_id + ' failed with status=' + str(
                    status_code) + ': ' + str(response))

    return response


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

    status_code, response = return_api_get_response(f"{orb_url}/api/v1/sinks", token=token,
                                                    params={'limit': limit, 'offset': offset}, verify=verify_ssl_bool)

    assert_that(status_code, equal_to(200),
                'Request to list sinks failed with status=' + str(status_code) + ': ' + str(response))

    assert_that(response, has_key('sinks'), f"Response does not contain 'sinks' key: {response}")
    assert_that(response, has_key('total'), f"Response does not contain 'total' key: {response}")
    assert_that(response, has_key('offset'), f"Response does not contain 'offset' key: {response}")

    return response['sinks'], response['total'], response['offset']


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

    status_code, response = return_api_delete_response(f"{orb_url}/api/v1/sinks/{sink_id}", token=token,
                                                       verify=verify_ssl_bool)

    assert_that(status_code, equal_to(204),
                f"Request to delete sink id= {sink_id} failed with status= {status_code}. Response: {response}.")


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


def edit_sink(token, sink_id, sink_body, expected_status_code=200):
    """

    :param (str) token: used for API authentication
    :param (str) sink_id: that identifies the sink to be edited
    :param (str) sink_body: sink's existent data
    :param (int) expected_status_code: expected request's status code. Default:200.
    :returns: (dict) the edited agent group
    """

    status_code, response = return_api_put_response(f"{orb_url}/api/v1/sinks/{sink_id}",
                                                    request_body=sink_body, token=token, verify=verify_ssl_bool)

    assert_that(status_code, equal_to(expected_status_code),
                'Request to edit sink failed with status=' + str(status_code) + ":" + str(response))

    return response
