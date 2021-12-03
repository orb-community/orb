import docker
from steps import test_config


# def after_feature(context, feature):
#     context.execute_steps('''
#     Given that the user is logged in on orb account
#     Then cleanup agents
#     Then cleanup agent group
#     Then cleanup sinks
#     ''')
#     cleanup_container()


def cleanup_container():
    docker_client = docker.from_env()
    containers = docker_client.containers.list(filters={"name": test_config.LOCAL_AGENT_CONTAINER_NAME})
    if len(containers) == 1:
        containers[0].stop()
        containers[0].remove(force=True)
