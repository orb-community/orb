import requests
import json
# from pygments import highlight, lexers, formatters
import os
import docker
import re

base_orb_url = os.getenv('ORB_URL', "https://beta.orb.live")


def generate_token(email, password, base_orb_url = base_orb_url):
    if email != None and password != None: 
        headers={'Content-type':'application/json',  'Accept':'*/*'}
        token_request = requests.post(base_orb_url+'/api/v1/tokens', json= {'email': email, 'password': password}, headers=headers)
        if token_request.status_code >= 200 and token_request.status_code <300:
            return (((token_request.text).split(":"))[1]).split('"')[1]
        else:
            return token_request.status_code
    else:
        assert email != None
        assert password != None




def check_existing_agents(token, id= "" , base_orb_url = base_orb_url):
    agentDict = dict()
    check_agents = requests.get(base_orb_url+'/api/v1/agents'+id, headers={'Authorization': token})
    agents_as_json = json.loads(check_agents.text)
    formatted_json = json.dumps(agents_as_json, sort_keys=True, indent=4)
    # colorful_json = highlight(formatted_json, lexers.JsonLexer(), formatters.TerminalFormatter())
    # print(colorful_json)
    if id == "":
        for i in range(len(agents_as_json['agents'])):
            agentDict[agents_as_json['agents'][i]['name']] = {'id': agents_as_json['agents'][i]['id'], 'state': agents_as_json['agents'][i]['state']}
    else:
            agentDict[agents_as_json['name']] = {'id': agents_as_json['id'], 'state': agents_as_json['state']}
    return agentDict


def delete_agent(list_of_agents, token, base_orb_url = base_orb_url):
    agent_name = list_of_agents.keys()
    for agent in agent_name:
        delete_agent = requests.delete(base_orb_url+'/api/v1/agents/'+list_of_agents[agent]['id'], headers={'Authorization': token})
        assert delete_agent.status_code == 204

def new_agent(token, name, tagKey, tagValue, base_orb_url = base_orb_url):
    new_agent = requests.post(base_orb_url+'/api/v1/agents', json= {"name": name,"orb_tags": {tagKey: tagValue}, "validate_only": False}, headers={'Content-type':'application/json',  'Accept':'*/*', 'Authorization': token})
    return new_agent.status_code, json.loads(new_agent.text)


def check_logs(logs):
    match = False
    providing_log_regex = "'level':'info','ts':.+,'caller':'agent/rpc_to.go:.+','msg':'sending capabilities','value':'{.'schema_version.':.+,.'orb_agent.':{.'version.':.+},.'agent_tags.'.+.'backends.':{.'pktvisor.':{.'version.':.+.',.'data.':{.'taps.':{.'default_pcap.':{.'config.':{.'iface.':.+,.'pcap_source.':.+},.'input_type.':.'pcap.',.'interface.':.+.'"
    try:
        for i in logs:
            i_replaced = i.replace('"',"'")
            match_agent = re.findall(providing_log_regex, str(i_replaced))
            if match_agent:
                match = True
    except ValueError:
        print("Invalid container ID")
    return match

def orb_agent_ids(name):
    idList = list()
    client = docker.from_env()
    for container in client.containers.list():
        if name in str(container.image):
            idList.append(container.id)
    return idList

def orb_agent_logs(id):
    client = docker.from_env()
    for container in client.containers.list():
        if id in str(container.id):
            return container.logs().decode("utf-8").split("\n") 
