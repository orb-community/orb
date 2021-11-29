import docker
from steps import test_config


def after_all(context):
    cleanup_container()


def after_feature(context, feature):
    if feature.tags[0] == 'agentGroups':
        context.execute_steps('''
        Then all agents should be deleted from orb
        Then all agent groups should be deleted from orb
        ''')


def cleanup_container():
    docker_client = docker.from_env()
    containers = docker_client.containers.list(filters={"name": test_config.LOCAL_AGENT_CONTAINER_NAME})
    if len(containers) == 1:
        containers[0].stop()
        containers[0].remove(force=True)
