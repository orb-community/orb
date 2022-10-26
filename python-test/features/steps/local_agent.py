from utils import safe_load_json, random_string, threading_wait_until, return_port_to_run_docker_container
from behave import then, step
from hamcrest import *
from test_config import TestConfig, LOCAL_AGENT_CONTAINER_NAME
import docker
import subprocess
import shlex
from retry import retry
import threading
import json
from datetime import datetime
import ciso8601

configs = TestConfig.configs()
ignore_ssl_and_certificate_errors = configs.get('ignore_ssl_and_certificate_errors')


@step('the agent container is started on an {status_port} port')
def run_local_agent_container(context, status_port):
    assert_that(status_port, any_of(equal_to("available"), equal_to("unavailable")), "Unexpected value for port")
    availability = {"available": True, "unavailable": False}
    orb_address = configs.get('orb_address')
    interface = configs.get('orb_agent_interface', 'mock')
    agent_docker_image = configs.get('agent_docker_image', 'ns1labs/orb-agent')
    image_tag = ':' + configs.get('agent_docker_tag', 'latest')
    agent_image = agent_docker_image + image_tag
    include_otel_env_var = configs.get("include_otel_env_var")
    enable_otel = configs.get("enable_otel")
    env_vars = {"ORB_CLOUD_ADDRESS": orb_address,
                "ORB_CLOUD_MQTT_ID": context.agent['id'],
                "ORB_CLOUD_MQTT_CHANNEL_ID": context.agent['channel_id'],
                "ORB_CLOUD_MQTT_KEY": context.agent_key,
                "PKTVISOR_PCAP_IFACE_DEFAULT": interface}
    if ignore_ssl_and_certificate_errors == 'true':
        env_vars["ORB_TLS_VERIFY"] = "false"
    if include_otel_env_var == "true":
        env_vars["ORB_OTEL_ENABLE"] = enable_otel

    context.port = return_port_to_run_docker_container(context, availability[status_port])

    if context.port != 10583:
        env_vars["ORB_BACKENDS_PKTVISOR_API_PORT"] = str(context.port)

    context.container_id = run_agent_container(agent_image, env_vars, LOCAL_AGENT_CONTAINER_NAME +
                                               context.agent['name'][-5:])
    if context.container_id not in context.containers_id.keys():
        context.containers_id[context.container_id] = str(context.port)
    if availability[status_port]:
        log = f"web server listening on localhost:{context.port}"
    else:
        log = f"unable to bind to localhost:{context.port}"
    agent_started, logs = get_logs_and_check(context.container_id, log, element_to_check="log")
    assert_that(agent_started, equal_to(True), f"Log {log} not found on agent logs. Agent Name: {context.agent['name']}."
                                               f"\n Logs:{logs}")


@step('the container logs that were output after {condition} contain the message "{text_to_match}" within'
      '{time_to_wait} seconds')
def check_agent_logs_considering_timestamp(context, condition, text_to_match, time_to_wait):
    # todo improve the logic for timestamp
    if "reset" in condition:
        considered_timestamp = context.considered_timestamp_reset
    else:
        considered_timestamp = context.considered_timestamp
    text_found, logs = get_logs_and_check(context.container_id, text_to_match, considered_timestamp,
                                    timeout=time_to_wait)
    assert_that(text_found, is_(True), f"Message {text_to_match} was not found in the agent logs!. \n\n"
                                       f"Container logs: {json.dumps(logs, indent=4)}")


@step("the container logs should not contain any {type_of_message} message")
def check_errors_on_agent_logs(context, type_of_message):
    type_dict = {"error": '"level":"error"', "panic": 'panic:'}
    logs = get_orb_agent_logs(context.container_id)
    non_expected_logs = [log for log in logs if type_dict[type_of_message] in log]
    assert_that(len(non_expected_logs), equal_to(0), f"agents logs contain the following {type_of_message}: "
                                                     f"{non_expected_logs}. \n All logs: {logs}.")


@then('the container logs should contain the message "{text_to_match}" within {time_to_wait} seconds')
def check_agent_msg_in_logs(context, text_to_match, time_to_wait):
    text_found, logs = get_logs_and_check(context.container_id, text_to_match, timeout=time_to_wait)

    assert_that(text_found, is_(True), f"Message {text_to_match} was not found in the agent logs!. \n\n"
                                       f"Container logs: {json.dumps(logs, indent=4)}")


@then('the container logs should contain "{error_log}" as log within {time_to_wait} seconds')
def check_agent_log_in_logs(context, error_log, time_to_wait):
    error_log = error_log.replace(":port", f":{context.port}")
    text_found, logs = get_logs_and_check(context.container_id, error_log, element_to_check="log", timeout=time_to_wait)
    assert_that(text_found, is_(True), f"Log {error_log} was not found in the agent logs!. \n\n"
                                       f"Container logs: {json.dumps(logs, indent=4)}")


@step("{order} container created is {status} within {seconds} seconds")
def check_last_container_status(context, order, status, seconds):
    order_convert = {"first": 0, "last": -1, "second": 1}
    container = list(context.containers_id.keys())[order_convert[order]]
    container_status = check_container_status(container, status, timeout=seconds)
    assert_that(container_status, equal_to(status), f"Container {context.container_id} failed with status "
                                                    f"{container_status}")


@step("{order} container created is {status} after {seconds} seconds")
def check_last_container_status_after_time(context, order, status, seconds):
    event = threading.Event()
    event.wait(int(seconds))
    event.set()
    if event.is_set() is True:
        check_last_container_status(context, order, status, seconds)


@step("the agent container is started using the command provided by the UI on an {status_port} port")
def run_container_using_ui_command(context, status_port):
    assert_that(status_port, any_of(equal_to("available"), equal_to("unavailable")), "Unexpected value for port")
    availability = {"available": True, "unavailable": False}
    context.port = return_port_to_run_docker_container(context, availability[status_port])
    include_otel_env_var = configs.get("include_otel_env_var")
    enable_otel = configs.get("enable_otel")
    context.container_id = run_local_agent_from_terminal(context.agent_provisioning_command,
                                                         ignore_ssl_and_certificate_errors, str(context.port),
                                                         include_otel_env_var, enable_otel)
    assert_that(context.container_id, is_not((none())), f"Agent container was not run")
    rename_container(context.container_id, LOCAL_AGENT_CONTAINER_NAME + context.agent['name'][-5:])
    if context.container_id not in context.containers_id.keys():
        context.containers_id[context.container_id] = str(context.port)


@step(
    "the agent container is started using the command provided by the UI without {parameter_to_remove} on an {"
    "status_port} port")
def run_container_using_ui_command_without_restart(context, parameter_to_remove, status_port):
    context.agent_provisioning_command = context.agent_provisioning_command.replace(f"{parameter_to_remove} ", "")
    run_container_using_ui_command(context, status_port)


@step("stop the orb-agent container")
def stop_orb_agent_container(context):
    for container_id in context.containers_id.keys():
        stop_container(container_id)


@step("remove the orb-agent container")
def remove_orb_agent_container(context):
    for container_id in context.containers_id.keys():
        remove_container(container_id)
    context.containers_id = {}


@step("forced remove the orb-agent container")
def remove_orb_agent_container(context):
    for container_id in context.containers_id.keys():
        remove_container(container_id, force_remove=True)
    context.containers_id = {}


@step("force remove of all agent containers whose names start with the test prefix")
def remove_all_orb_agent_test_containers(context):
    docker_client = docker.from_env()
    containers = docker_client.containers.list(all=True)
    for container in containers:
        test_container = container.name.startswith(LOCAL_AGENT_CONTAINER_NAME)
        if test_container is True:
            container.remove(force=True)


def run_agent_container(container_image, env_vars, container_name, time_to_wait=5):
    """
    Gets a specific agent from Orb control plane

    :param (str) container_image: that will be used for running the container
    :param (dict) env_vars: that will be passed to the container context
    :param (str) container_name: base of container name
    :param (int) time_to_wait: seconds that threading must wait after run the agent
    :returns: (str) the container ID
    """
    client = docker.from_env()
    container = client.containers.run(container_image, name=container_name, detach=True,
                                      network_mode='host', environment=env_vars)
    threading.Event().wait(time_to_wait)
    return container.id


def get_orb_agent_logs(container_id):
    """
    Gets the logs from Orb agent container

    :param (str) container_id: specifying the orb agent container
    :returns: (list) of log lines
    """
    docker_client = docker.from_env()
    container = docker_client.containers.get(container_id)
    return container.logs().decode("utf-8").split("\n")


def check_logs_contain_message(logs, expected_message, event, start_time=0):
    """
    Gets the logs from Orb agent container

    :param (list) logs: list of log lines
    :param (str) expected_message: message that we expect to find in the logs
    :param (obj) event: threading.event
    :param (int) start_time: time to be considered as initial time. Default: None
    :returns: (bool) whether expected message was found in the logs
    """

    for log_line in logs:
        log_line = safe_load_json(log_line)

        if log_line is not None and log_line['msg'] == expected_message and isinstance(log_line['ts'], int) and \
                log_line['ts'] > start_time:
            event.set()
            return event.is_set()
        elif log_line is not None and log_line['msg'] == expected_message and isinstance(log_line['ts'], str) and \
                datetime.timestamp(ciso8601.parse_datetime(log_line['ts'])) > start_time:
            event.set()
            return event.is_set()

    return event.is_set()


def check_logs_contain_log(logs, expected_log, event, start_time=0):
    """
    Check if the logs from Orb agent container contain specific log

    :param (list) logs: list of log lines
    :param (str) expected_log: log that we expect to find in the logs
    :param (obj) event: threading.event
    :param (int) start_time: time to be considered as initial time. Default: None
    :returns: (bool) whether expected message was found in the logs
    """

    for log_line in logs:
        log_line = safe_load_json(log_line)

        if log_line is not None and "log" in log_line.keys() and expected_log in log_line['log'] and isinstance(
                log_line['ts'], int) and \
                log_line['ts'] > start_time:
            event.set()
            return event.is_set()
        elif log_line is not None and "log" in log_line.keys() and expected_log in log_line['log'] and isinstance(
                log_line['ts'], str) and \
                datetime.timestamp(ciso8601.parse_datetime(log_line['ts'])) > start_time:
            event.set()
            return event.is_set()

    return event.is_set()


def run_local_agent_from_terminal(command, ignore_ssl_and_certificate_errors, pktvisor_port,
                                  include_otel_env_var="false", enable_otel="false"):
    """
    :param (str) command: docker command to provision an agent
    :param (bool) ignore_ssl_and_certificate_errors: True if orb address doesn't have a valid certificate.
    :param (str or int) pktvisor_port: Port on which pktvisor should run
    :param (str): if 'true', ORB_OTEL_ENABLE env ver is included on command provisioning of the agent
    :return: agent container ID
    """
    command = command.replace("\\\n", " ")
    args = shlex.split(command)
    if ignore_ssl_and_certificate_errors == 'true':
        args.insert(-1, "-e")
        args.insert(-1, "ORB_TLS_VERIFY=false")
    if include_otel_env_var == "true":
        args.insert(-1, f"ORB_OTEL_ENABLE={enable_otel}")
    if pktvisor_port != 'default':
        args.insert(-1, "-e")
        args.insert(-1, f"ORB_BACKENDS_PKTVISOR_API_PORT={pktvisor_port}")
    terminal_running = subprocess.Popen(
        args, stdout=subprocess.PIPE)
    subprocess_return = terminal_running.stdout.read().decode()
    container_id = subprocess_return.split()
    assert_that(container_id[0], is_not((none())), f"Failed to run the agent. Command used: {args}.")
    return container_id[0]


@retry(tries=5, delay=0.2)
def rename_container(container_id, container_name):
    """

    :param container_id: agent container ID
    :param container_name: base of agent container name
    """
    docker_client = docker.from_env()
    containers = docker_client.containers.list(all=True)
    is_container_up = any(container_id in container.id for container in containers)
    assert_that(is_container_up, equal_to(True), f"Container {container_id} not found")
    container_name = container_name + random_string(5)
    rename_container_command = f"docker rename {container_id} {container_name}"
    rename_container_args = shlex.split(rename_container_command)
    subprocess.Popen(rename_container_args, stdout=subprocess.PIPE)


@threading_wait_until
def check_container_status(container_id, status, event=None):
    """

    :param container_id: agent container ID
    :param status: status that we expect to find in the container
    :param event: threading.event
    :return status of the container
    """
    docker_client = docker.from_env()
    container = docker_client.containers.list(all=True, filters={'id': container_id})
    assert_that(container, has_length(1), f"unable to find container {container_id}.")
    container = container[0]
    if container.status == status:
        event.set()
    return container.status


@threading_wait_until
def get_logs_and_check(container_id, expected_message, start_time=0, element_to_check="msg", event=None):
    """

    :param container_id: agent container ID
    :param (str) expected_message: message that we expect to find in the logs
    :param (int) start_time: time to be considered as initial time. Default: None
    :param element_to_check: Part of the log to be validated. Options: "msg" and "log". Default: "msg".
    :param (obj) event: threading.event
    :return: (bool) if the expected message is found return True, if not, False
    """
    assert_that(element_to_check, any_of(equal_to("msg"), equal_to("log")), "Unexpected value for element to check.")
    logs = get_orb_agent_logs(container_id)
    if element_to_check == "msg":
        text_found = check_logs_contain_message(logs, expected_message, event, start_time)
    else:
        text_found = check_logs_contain_log(logs, expected_message, event, start_time)
    return text_found, logs


def run_agent_config_file(agent_name, time_to_wait=5):
    """
    Run an agent container using an agent config file

    :param agent_name: name of the orb agent
    :param time_to_wait: seconds that threading must wait after run the agent
    :return: agent container id
    """
    agent_docker_image = configs.get('agent_docker_image', 'ns1labs/orb-agent')
    agent_image = f"{agent_docker_image}:{configs.get('agent_docker_tag', 'latest')}"
    local_orb_path = configs.get("local_orb_path")
    volume = f"{local_orb_path}:/usr/local/orb/"
    agent_command = f"/usr/local/orb/{agent_name}.yaml"
    command = f"docker run -d -v {volume} --net=host {agent_image} run -c {agent_command}"
    args = shlex.split(command)
    terminal_running = subprocess.Popen(args, stdout=subprocess.PIPE)
    subprocess_return = terminal_running.stdout.read().decode()
    container_id = subprocess_return.split()[0]
    rename_container(container_id, LOCAL_AGENT_CONTAINER_NAME + agent_name[-5:])
    threading.Event().wait(time_to_wait)
    return container_id


def stop_container(container_id):
    """

    :param container_id: agent container ID
    """
    docker_client = docker.from_env()
    container = docker_client.containers.get(container_id)
    container.stop()


def remove_container(container_id, force_remove=False):
    """

    :param container_id: agent container ID
    :param force_remove: if True, similar to docker rm -f. Default: False
    """
    docker_client = docker.from_env()
    container = docker_client.containers.get(container_id)
    container.remove(force=force_remove)
