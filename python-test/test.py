import docker

def orb_agent_id(name):
    id_list = list()
    client = docker.from_env()
    for container in client.containers.list():
        if name in str(container.image):
        # print(container.logs())
            id_list.append(container.id)
    return id_list

def remove_docker_image(id):
    pass

print(orb_agent_id('orb-agent'))
