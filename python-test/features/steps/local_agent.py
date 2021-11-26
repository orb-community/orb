from json import loads, JSONDecodeError
from behave import when, then
from hamcrest import *
from test_config import TestConfig, LOCAL_AGENT_CONTAINER_NAME
import docker

configs = TestConfig.configs()


@when('the agent container is started')
def run_local_agent_container(context):
    orb_address = configs.get('orb_address')
    interface = configs.get('orb_agent_interface', 'mock')
    image_tag = ':' + configs.get('agent_image_tag', 'latest')
    agent_image = "ns1labs/orb-agent" + image_tag

    env_vars = {"ORB_CLOUD_ADDRESS": orb_address,
                "ORB_CLOUD_MQTT_ID": context.agent['id'],
                "ORB_CLOUD_MQTT_CHANNEL_ID": context.agent['channel_id'],
                "ORB_CLOUD_MQTT_KEY": context.agent['key'],
                "PKTVISOR_PCAP_IFACE_DEFAULT": interface}

    context.container_id = run_agent_container(agent_image, env_vars)


@then('the container logs should contain the message "{text_to_match}"')
def check_agent_log(context, text_to_match):
    logs = get_orb_agent_logs(context.container_id)
    contains_text = check_logs_contains_message(logs, text_to_match)
    assert_that(contains_text, is_(True), 'Message "' + text_to_match + '" was not found in the agent logs!')


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


def check_logs_contains_message(logs, expected_message):
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


# TODO: move this function to some sort of utils.py module
def safe_load_json(json_str):
    """
    Safely parses a string into a JSON object, without ever raising an error.
    :param (str) json_str: to be loaded
    :return: the JSON object, or None if string is not a valid JSON.
    """

    try:
        return loads(json_str)
    except JSONDecodeError:
        return None
