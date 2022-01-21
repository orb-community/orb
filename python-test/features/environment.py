import docker
from steps import test_config


def before_scenario(context, scenario):
    cleanup_container()


def after_feature(context, feature):
    cleanup_container()
    context.execute_steps('''
    Given the Orb user logs in
    Then cleanup agents
    Then cleanup agent group
    Then cleanup sinks
    Then cleanup policies
    Then cleanup datasets
    ''')


def cleanup_container():
    docker_client = docker.from_env()
    containers = docker_client.containers.list(filters={"name": test_config.LOCAL_AGENT_CONTAINER_NAME})
    if len(containers) == 1:
        containers[0].stop()
        containers[0].remove()
