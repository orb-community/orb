import requests
import json
from pygments import highlight, lexers, formatters
import subprocess
import shlex
import re
import time
import docker

base_orb_url = "https://beta.orb.live"
email = 'tester@email.com'
password = '12345678'


def generate_token(email, password, base_orb_url = "https://beta.orb.live"): 
    headers={'Content-type':'application/json',  'Accept':'*/*'}
    token_request = requests.post(base_orb_url+'/api/v1/tokens', json= {'email': email, 'password': password}, headers=headers)
    if token_request.status_code >= 200 and token_request.status_code <300:
        return (((token_request.text).split(":"))[1]).split('"')[1]
    else:
        return token_request.status_code


def check_existing_agents(token, id= "" , base_orb_url = "https://beta.orb.live"):
    agentDict = dict()
    check_agents = requests.get(base_orb_url+'/api/v1/agents'+id, headers={'Authorization': token})
    agents_as_json = json.loads(check_agents.text)
    formatted_json = json.dumps(agents_as_json, sort_keys=True, indent=4)
    colorful_json = highlight(formatted_json, lexers.JsonLexer(), formatters.TerminalFormatter())
    # print(colorful_json)
    if id == "":
        for i in range(len(agents_as_json['agents'])):
            agentDict[agents_as_json['agents'][i]['name']] = {'id': agents_as_json['agents'][i]['id'], 'state': agents_as_json['agents'][i]['state']}
    else:
            agentDict[agents_as_json['name']] = {'id': agents_as_json['id'], 'state': agents_as_json['state']}
    return agentDict


def delete_agent(list_of_agents, token, base_orb_url = "https://beta.orb.live"):
    agent_name = list_of_agents.keys()
    for agent in agent_name:
        delete_agent = requests.delete(base_orb_url+'/api/v1/agents/'+list_of_agents[agent]['id'], headers={'Authorization': token})
        assert delete_agent.status_code == 204

def new_agent(token, name, tagKey, tagValue, base_orb_url = "https://beta.orb.live"):
    new_agent = requests.post(base_orb_url+'/api/v1/agents', json= {"name": name,"orb_tags": {tagKey: tagValue}, "validate_only": False}, headers={'Content-type':'application/json',  'Accept':'*/*', 'Authorization': token})
    assert new_agent.status_code == 201
    return json.loads(new_agent.text)

def send_terminal_commands(command, separator=None, cwd_run=None):
    args = shlex.split(command)
    # print(args)
    docker_command_execute = subprocess.Popen(args, stdout=subprocess.PIPE, cwd=cwd_run)
    subprocess_return = docker_command_execute.stdout.read().decode()
    if separator == None:
        subprocess_return_terminal = subprocess_return.split()
    else:
        subprocess_return_terminal = re.split(separator, subprocess_return)
    return subprocess_return_terminal

def agent_provising(agentname, tagKey, tagValue, token):
    timeout = 0
    agent_data = new_agent(token, agentname, tagKey, tagValue)
    agent_key = agent_data['key']
    agent_id = agent_data['id']
    agent_channel_id = agent_data['channel_id']
    provision_agent_command = f"docker run -d --net=host -e ORB_CLOUD_ADDRESS=beta.orb.live -e ORB_CLOUD_MQTT_ID={agent_id} -e ORB_CLOUD_MQTT_CHANNEL_ID={agent_channel_id} -e ORB_CLOUD_MQTT_KEY={agent_key} -e PKTVISOR_PCAP_IFACE_DEFAULT=mock ns1labs/orb-agent"
    print(send_terminal_commands(provision_agent_command))
    agent_mode = check_existing_agents(token, '/'+agent_id)[agentname]['state']
    while agent_mode != "online" and timeout < 5:
        time.sleep(0.1)
        agent_mode = check_existing_agents(token, '/'+agent_id)[agentname]['state']
    assert agent_mode == "online"
    return agent_mode

def orb_agent_ids(name):
    idList = list()
    client = docker.from_env()
    for container in client.containers.list():
        if name in str(container.image):
            idList.append(container.id)
    return idList

def remove_docker_image(id_list):
    for i in id_list:
        command = f"docker rm -f {i}"
        print(send_terminal_commands(command))

token = generate_token(email, password)

print(agent_provising('agent1', 'test', 'true', token))

list_agents = check_existing_agents(token)

delete_agent(list_agents, token)

docker_agents = orb_agent_ids("orb-agent")

remove_docker_image(docker_agents)
