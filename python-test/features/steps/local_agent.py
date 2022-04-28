from utils import safe_load_json, random_string, threading_wait_until, check_port_is_available
from behave import then, step
from hamcrest import *
from test_config import TestConfig, LOCAL_AGENT_CONTAINER_NAME
import docker
import subprocess
import shlex
from retry import retry
import threading

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
    env_vars = {"ORB_CLOUD_ADDRESS": orb_address,
                "ORB_CLOUD_MQTT_ID": context.agent['id'],
                "ORB_CLOUD_MQTT_CHANNEL_ID": context.agent['channel_id'],
                "ORB_CLOUD_MQTT_KEY": context.agent_key,
                "PKTVISOR_PCAP_IFACE_DEFAULT": interface}
    if ignore_ssl_and_certificate_errors == 'true':
        env_vars["ORB_TLS_VERIFY"] = "false"

    context.port = check_port_is_available(availability[status_port])

    if context.port != 10583:
        env_vars["ORB_BACKENDS_PKTVISOR_API_PORT"] = str(context.port)

    context.container_id = run_agent_container(agent_image, env_vars, LOCAL_AGENT_CONTAINER_NAME)
    if context.container_id not in context.containers_id.keys():
        context.containers_id[context.container_id] = str(context.port)


@step('the container logs that were output after {condition} contain the message "{text_to_match}" within'
      '{time_to_wait} seconds')
def check_agent_logs_considering_timestamp(context, condition, text_to_match, time_to_wait):
    #todo improve the logic for timestamp
    if "reset" in condition:
        considered_timestamp = context.considered_timestamp_reset
    else:
        considered_timestamp = context.considered_timestamp
    text_found = get_logs_and_check(context.container_id, text_to_match, considered_timestamp,
                                    timeout=time_to_wait)

    assert_that(text_found, is_(True), 'Message "' + text_to_match + '" was not found in the agent logs!')


@then('the container logs should contain the message "{text_to_match}" within {time_to_wait} seconds')
def check_agent_log(context, text_to_match, time_to_wait):
    text_found = get_logs_and_check(context.container_id, text_to_match, timeout=time_to_wait)

    assert_that(text_found, is_(True), 'Message "' + text_to_match + '" was not found in the agent logs!')


@step("{order} container created is {status} after {seconds} seconds")
def check_last_container_status(context, order, status, seconds):
    event = threading.Event()
    event.wait(int(seconds))
    event.set()
    if event.is_set() is True:
        order_convert = {"first": 0, "last": -1, "second": 1}
        container = list(context.containers_id.keys())[order_convert[order]]
        container_status = check_container_status(container, status, timeout=seconds)
        assert_that(container_status, equal_to(status), f"Container {context.container_id} failed with status"
                                                        f"{container_status}")


@step("the agent container is started using the command provided by the UI on an {status_port} port")
def run_container_using_ui_command(context, status_port):
    assert_that(status_port, any_of(equal_to("available"), equal_to("unavailable")), "Unexpected value for port")
    availability = {"available": True, "unavailable": False}
    context.port = check_port_is_available(availability[status_port])
    context.container_id = run_local_agent_from_terminal(context.agent_provisioning_command,
                                                         ignore_ssl_and_certificate_errors, str(context.port))
    assert_that(context.container_id, is_not((none())))
    rename_container(context.container_id, LOCAL_AGENT_CONTAINER_NAME)
    if context.container_id not in context.containers_id.keys():
        context.containers_id[context.container_id] = str(context.port)


def run_agent_container(container_image, env_vars, container_name):
    """
    Gets a specific agent from Orb control plane

    :param (str) container_image: that will be used for running the container
    :param (dict) env_vars: that will be passed to the container context
    :param (str) container_name: base of container name
    :returns: (str) the container ID
    """
    LOCAL_AGENT_CONTAINER_NAME = container_name + random_string(5)
    client = docker.from_env()
    container = client.containers.run(container_image, name=LOCAL_AGENT_CONTAINER_NAME, detach=True,
                                      network_mode='host', environment=env_vars)
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

        if log_line is not None and log_line['msg'] == expected_message and log_line['ts'] > start_time:
            event.set()
            return event.is_set()

    return event.is_set()


def run_local_agent_from_terminal(command, ignore_ssl_and_certificate_errors, pktvisor_port):
    """
    :param (str) command: docker command to provision an agent
    :param (bool) ignore_ssl_and_certificate_errors: True if orb address doesn't have a valid certificate.
    :param (str or int) pktvisor_port: Port on which pktvisor should run
    :return: agent container ID
    """
    command = command.replace("\\\n", " ")
    args = shlex.split(command)
    if ignore_ssl_and_certificate_errors == 'true':
        args.insert(-1, "-e")
        args.insert(-1, "ORB_TLS_VERIFY=false")
    if pktvisor_port != 'default':
        args.insert(-1, "-e")
        args.insert(-1, f"ORB_BACKENDS_PKTVISOR_API_PORT={pktvisor_port}")
    terminal_running = subprocess.Popen(
        args, stdout=subprocess.PIPE)
    subprocess_return = terminal_running.stdout.read().decode()
    container_id = subprocess_return.split()
    assert_that(container_id[0], is_not((none())))
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
    assert_that(container, has_length(1))
    container = container[0]
    if container.status == status:
        event.set()
    return container.status


@threading_wait_until
def get_logs_and_check(container_id, expected_message, start_time=0, event=None):
    """

    :param container_id: agent container ID
    :param (str) expected_message: message that we expect to find in the logs
    :param (int) start_time: time to be considered as initial time. Default: None
    :param (obj) event: threading.event
    :return: (bool) if the expected message is found return True, if not, False
    """
    logs = get_orb_agent_logs(container_id)
    text_found = check_logs_contain_message(logs, expected_message, event, start_time)
    return text_found


def run_agent_config_file(orb_path, agent_name):
    """
    Run an agent container using an agent config file

    :param orb_path: path to orb directory
    :param agent_name: name of the orb agent
    :return: agent container id
    """
    agent_docker_image = configs.get('agent_docker_image', 'ns1labs/orb-agent')
    agent_image = f"{agent_docker_image}:{configs.get('agent_docker_tag', 'latest')}"
    volume = f"{orb_path}:/usr/local/orb/"
    agent_command = f"/usr/local/orb/{agent_name}.yaml"
    command = f"docker run -d -v {volume} --net=host {agent_image} run -c {agent_command}"
    args = shlex.split(command)
    terminal_running = subprocess.Popen(args, stdout=subprocess.PIPE)
    subprocess_return = terminal_running.stdout.read().decode()
    container_id = subprocess_return.split()[0]
    rename_container(container_id, LOCAL_AGENT_CONTAINER_NAME)
    return container_id
