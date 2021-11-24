from behave import given, when, then
import time
import functions
import os
import docker
import string
import random
import hamcrest

email = os.getenv('EMAIL')
password = os.getenv('PASSWORD')
randomAgentName = ''.join(random.choices(string.ascii_letters, k=10)) #k sets the number of characters
agentname = os.getenv('AGENT_NAME', randomAgentName)
agentTagKey = os.getenv('TAG_KEY', 'test')
agentTagValue = os.getenv('TAG_VALUE', 'true')


@given('A valid authentication')
def get_auth(context):
    context.token = functions.generate_token(email, password)['token']
    hamcrest.assert_that(len(context.token), hamcrest.greater_than(0), "Token")
    

@when('Create an agent')
def create_agent(context):
    context.status_agent, context.agent_data = functions.new_agent(context.token, agentname, agentTagKey, agentTagValue)
    hamcrest.assert_that(context.status_agent, hamcrest.equal_to(201), 'Request Token- status code')
    context.agent_key = context.agent_data['key']
    context.agent_id = context.agent_data['id']
    context.agent_channel_id = context.agent_data['channel_id']

@when('Run agent container')
def run_agent(context):
    iface = os.getenv('IFACE', 'mock')
    context.orb_address = os.getenv('ORB_ADDRESS', 'beta.orb.live')
    context.tag_image = os.getenv('AGENT_DOCKER_TAG')
    if context.tag_image != None:
        context.agent_image = f"ns1labs/orb-agent:{context.tag_image}"
    else:
        context.agent_image = "ns1labs/orb-agent"
    var_env ={"ORB_CLOUD_ADDRESS":context.orb_address, "ORB_CLOUD_MQTT_ID":context.agent_id,"ORB_CLOUD_MQTT_CHANNEL_ID":context.agent_channel_id, "ORB_CLOUD_MQTT_KEY":context.agent_key,"PKTVISOR_PCAP_IFACE_DEFAULT":iface}
    client = docker.from_env()
    context.agent_container = client.containers.run(context.agent_image, detach=True, network_mode= 'host', environment=var_env)
    hamcrest.assert_that(len(context.agent_container.id), hamcrest.greater_than(0), 'Container ID')
    context.id_terminal =  context.agent_container.id

@then("Agent should be online")
def check_container_status(context):
    timeout = 0
    context.agent_mode = functions.check_existing_agents(context.token, '/'+context.agent_id)[agentname]['state']
    while context.agent_mode != "online" and timeout < 5:
        time.sleep(0.1)
        context.agent_mode = functions.check_existing_agents(context.token, '/'+context.agent_id)[agentname]['state']
        timeout = timeout + 0.1
    hamcrest.assert_that(context.agent_mode, hamcrest.contains_string('online'), 'Agent State')
    
@then("Container logs should be sending capabilities")
def check_agent_log(context):
    context.logs = functions.orb_agent_logs(context.id_terminal)
    context.match = functions.check_logs(context.logs)
    hamcrest.assert_that(context.match, hamcrest.is_(True), 'Sending Capabilities')
