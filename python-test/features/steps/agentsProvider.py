from behave import given, when, then
import time
import functions
import os
import docker

email = os.getenv('EMAIL')
password = os.getenv('PASSWORD')

agentname = os.getenv('AGENT_NAME')
agentTagKey = os.getenv('TAG_KEY', 'test')
agentTagValeu = os.getenv('TAG_VALUE', 'true')


@given('A valid authentication')
def get_auth(context):
    context.token = functions.generate_token(email, password)
    assert len(context.token) > 0

@when('Create an agent')
def create_agent(context):
    context.status_agent, context.agent_data = functions.new_agent(context.token, agentname, agentTagKey, agentTagValeu)
    assert context.status_agent == 201
    context.agent_key = context.agent_data['key']
    context.agent_id = context.agent_data['id']
    context.agent_channel_id = context.agent_data['channel_id']

@when('Run agent container')
def run_agent(context):
    iface = os.getenv('IFACE', 'mock')
    context.orb_address = os.getenv('ORB_ADDRESS', 'beta.orb.live')
    context.version_image = os.getenv('AGENT_VERSION')
    context.agent_image = f"ns1labs/orb-agent:{context.version_image}"
    var_env ={"ORB_CLOUD_ADDRESS":context.orb_address, "ORB_CLOUD_MQTT_ID":context.agent_id,"ORB_CLOUD_MQTT_CHANNEL_ID":context.agent_channel_id, "ORB_CLOUD_MQTT_KEY":context.agent_key,"PKTVISOR_PCAP_IFACE_DEFAULT":iface}
    client = docker.from_env()
    context.agent_container = client.containers.run(context.agent_image, detach=True, network_mode= 'host', environment=var_env)
    assert len(context.agent_container.id) > 0
    context.id_terminal =  context.agent_container.id

@then("Agent should be online")
def check_container_status(context):
    timeout = 0
    context.agent_mode = functions.check_existing_agents(context.token, '/'+context.agent_id)[agentname]['state']
    while context.agent_mode != "online" and timeout < 5:
        time.sleep(0.1)
        context.agent_mode = functions.check_existing_agents(context.token, '/'+context.agent_id)[agentname]['state']
        timeout = timeout + 0.1
    assert context.agent_mode == "online"

@then("Container logs should be sending capabilities")
def check_agent_log(context):
    context.logs = functions.orb_agent_logs(context.id_terminal)
    context.match = functions.check_logs(context.logs)
    assert context.match == True  
