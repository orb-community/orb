import docker
from steps import test_config


def before_scenario(context, scenario):
    context.containers_id = dict()
    context.agent_groups = dict()


def after_scenario(context, feature):
    # cleanup_container()
    context.execute_steps('''
    Then remove the container
    ''')


def cleanup_container():
    docker_client = docker.from_env()
    containers = docker_client.containers.list(all=True)
    for container in containers:
        test_container = container.name.startswith(test_config.LOCAL_AGENT_CONTAINER_NAME)
        if test_container is True:
            container.stop()
            container.remove()
