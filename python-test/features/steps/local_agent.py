from utils import safe_load_json
from behave import when, then
from hamcrest import *
from test_config import TestConfig, LOCAL_AGENT_CONTAINER_NAME
import docker
import time
import subprocess
import shlex

configs = TestConfig.configs()
bypass_ssl_certificate_check = configs.get('bypass_ssl_certificate_check')


@when('the agent container is started')
def run_local_agent_container(context):
    orb_address = configs.get('orb_address')
    interface = configs.get('orb_agent_interface', 'mock')
    agent_docker_image = configs.get('agent_docker_image', 'ns1labs/orb-agent')
    image_tag = ':' + configs.get('agent_docker_tag', 'latest')
    agent_image = agent_docker_image + image_tag
    env_vars = {"ORB_CLOUD_ADDRESS": orb_address,
                "ORB_CLOUD_MQTT_ID": context.agent['id'],
                "ORB_CLOUD_MQTT_CHANNEL_ID": context.agent['channel_id'],
                "ORB_CLOUD_MQTT_KEY": context.agent['key'],
                "PKTVISOR_PCAP_IFACE_DEFAULT": interface}
    if bypass_ssl_certificate_check:
        env_vars["ORB_TLS_VERIFY"] = "false"

    context.container_id = run_agent_container(agent_image, env_vars)


@then('the container logs should contain the message "{text_to_match}" within {time_to_wait} seconds')
def check_agent_log(context, text_to_match, time_to_wait):
    time_waiting = 0
    sleep_time = 0.5
    timeout = int(time_to_wait)
    text_found = False

    while time_waiting < timeout:
        logs = get_orb_agent_logs(context.container_id)
        text_found = check_logs_contain_message(logs, text_to_match)
        if text_found is True:
            break
        time.sleep(sleep_time)
        time_waiting += sleep_time

    assert_that(text_found, is_(True), 'Message "' + text_to_match + '" was not found in the agent logs!')


@when("the agent container is started using the command provided by the UI")
def run_container_using_ui_command(context):
    context.container_id = run_local_agent_from_terminal(context.agent_provisioning_command,
                                                         bypass_ssl_certificate_check)
    assert_that(context.container_id, is_not((none())))
    rename_container(context.container_id, LOCAL_AGENT_CONTAINER_NAME)


def run_agent_container(container_image, env_vars):
    """
    Gets a specific agent from Orb control plane

    :param (str) container_image: that will be used for running the container
    :param (dict) env_vars: that will be passed to the container context
    :returns: (str) the container ID
    """

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


def check_logs_contain_message(logs, expected_message):
    """
    Gets the logs from Orb agent container

    :param (list) logs: list of log lines
    :param (str) expected_message: message that we expect to find in the logs
    :returns: (bool) whether expected message was found in the logs
    """

    for log_line in logs:
        log_line = safe_load_json(log_line)

        if log_line is not None and log_line['msg'] == expected_message:
            return True

    return False


def run_local_agent_from_terminal(command, bypass_check_ssl_certificate):
    """
    :param (str) command: docker command to provision an agent
    :param (bool) bypass_check_ssl_certificate: True if orb address doesn't have a valid certificate.
    :return: agent container ID
    """
    args = shlex.split(command)
    if bypass_check_ssl_certificate:
        args.insert(-1, "-e")
        args.insert(-1, "ORB_TLS_VERIFY=false")
    terminal_running = subprocess.Popen(
        args, stdout=subprocess.PIPE)
    subprocess_return = terminal_running.stdout.read().decode()
    container_id = subprocess_return.split()
    assert_that(container_id[0], is_not((none())))
    return container_id[0]


def rename_container(container_id, container_name):
    """

    :param container_id: agent container ID
    :param container_name: agent container name
    """
    rename_container_command = f"docker rename {container_id} {container_name}"
    rename_container_args = shlex.split(rename_container_command)
    subprocess.Popen(rename_container_args, stdout=subprocess.PIPE)
