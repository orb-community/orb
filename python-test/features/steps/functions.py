import requests
import json
import os
import docker
import re
import hamcrest

base_orb_url = "https://" + os.getenv('ORB_ADDRESS', "beta.orb.live")

def generate_token(email, password, base_orb_url = base_orb_url):
    if email != None and password != None: 
        headers={'Content-type':'application/json',  'Accept':'*/*'}
        token_request = requests.post(base_orb_url+'/api/v1/tokens', json= {'email': email, 'password': password}, headers=headers)
        hamcrest.assert_that(token_request.json(), hamcrest.has_key('token'), 'Request Json - token')
        return token_request.json()
    else:
        hamcrest.assert_that(email, hamcrest.not_none(), 'Email')
        hamcrest.assert_that(password, hamcrest.not_none(), 'Password')

def check_existing_agents(token, id="", base_orb_url = base_orb_url):
    agentDict = dict()
    check_agents = requests.get(base_orb_url+'/api/v1/agents'+id, headers={'Authorization': token})
    agents_as_json = check_agents.json()
    if id == "":
        for i in range(len(agents_as_json['agents'])):
            agentDict[agents_as_json['agents'][i]['name']] = {'id': agents_as_json['agents'][i]['id'], 'state': agents_as_json['agents'][i]['state']}
    else:
            agentDict[agents_as_json['name']] = {'id': agents_as_json['id'], 'state': agents_as_json['state']}
    agentDict[agents_as_json['name']] = {'id': agents_as_json['id'], 'state': agents_as_json['state']}
    return agentDict

def delete_agent(list_of_agents, token, base_orb_url = base_orb_url):
    agent_name = list_of_agents.keys()
    for agent in agent_name:
        delete_agent = requests.delete(base_orb_url+'/api/v1/agents/'+list_of_agents[agent]['id'], headers={'Authorization': token})
        hamcrest.assert_that(delete_agent.status_code, hamcrest.equal_to(204), 'Delete Agent- status code')

def new_agent(token, name, tagKey, tagValue, base_orb_url = base_orb_url):
    new_agent = requests.post(base_orb_url+'/api/v1/agents', json= {"name": name,"orb_tags": {tagKey: tagValue}, "validate_only": False}, headers={'Content-type':'application/json',  'Accept':'*/*', 'Authorization': token})
    return new_agent.status_code, json.loads(new_agent.text)

def match_keys(json, dict_keys):
    for i in dict_keys:
        hamcrest.assert_that(json, hamcrest.has_key(i))




def check_logs(logs):
    match = False
    try:
        for i in logs:
            logline = json.loads(i)
            if logline['msg'] == 'sending capabilities':
                match_keys(logline, ['level', 'ts', 'caller', 'msg', 'value'])
                match_keys(json.loads(logline['value']), ['schema_version', 'orb_agent', 'agent_tags', 'backends'])         
                match_keys(json.loads(logline['value'])['orb_agent'], ['version'])
                match_keys(json.loads(logline['value'])['backends'], ['pktvisor'])
                match_keys(json.loads(logline['value'])['backends']['pktvisor'], ['version', 'data']) 
                match_keys(json.loads(logline['value'])['backends']['pktvisor']['data'], ['taps'])
                match_keys(json.loads(logline['value'])['backends']['pktvisor']['data']['taps'], ['default_pcap'])
                match_keys(json.loads(logline['value'])['backends']['pktvisor']['data']['taps']['default_pcap'], ['config', 'input_type', 'interface'])
                match_keys(json.loads(logline['value'])['backends']['pktvisor']['data']['taps']['default_pcap']['config'], ['iface','pcap_source'])
                hamcrest.assert_that(logline['level'], hamcrest.contains_string('info')) 

                match = True
    except ValueError:
        pass
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
