import docker
from steps import test_config


def before_scenario(context, scenario):
    cleanup_container()
    context.containers_id = dict()
    context.agent_groups = dict()


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
    containers = docker_client.containers.list(all=True)
    for container in containers:
        test_container = container.name.startswith(test_config.LOCAL_AGENT_CONTAINER_NAME)
        if test_container is True:
            container.stop()
            container.remove()
