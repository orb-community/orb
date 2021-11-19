from behave import given, when, then
import time
import functions

base_orb_url = "https://beta.orb.live"
email = 'tester@email.com'
password = '12345678'

agentname = 'agent1'
agentTagKey = 'test'
agentTagValeu = 'true'


@given('A valid authentication')
def get_auth(context):
    context.token = functions.generate_token(email, password, base_orb_url)
    assert len(context.token) > 0

@when('Create an agent')
def create_agent(context):
    context.status_agent, context.agent_data = functions.new_agent(context.token, agentname, agentTagKey, agentTagValeu)
    assert context.status_agent == 201
    context.agent_key = context.agent_data['key']
    context.agent_id = context.agent_data['id']
    context.agent_channel_id = context.agent_data['channel_id']

@when('Build agent container')
def build_agent(context):
    iface = 'mock'
    provision_agent_command = f"docker run -d --net=host -e ORB_CLOUD_ADDRESS=beta.orb.live -e ORB_CLOUD_MQTT_ID={context.agent_id} -e ORB_CLOUD_MQTT_CHANNEL_ID={context.agent_channel_id} -e ORB_CLOUD_MQTT_KEY={context.agent_key} -e PKTVISOR_PCAP_IFACE_DEFAULT={iface} ns1labs/orb-agent"
    context.id_terminal = functions.send_terminal_commands(provision_agent_command, '[')
    assert len(context.id_terminal) > 0 

@then("Agente should be online")
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
    print(context.id_terminal)
    print(context.id_terminal)
    context.logs = functions.orb_agent_logs(context.id_terminal)
    context.match = functions.check_logs(context.logs)
    assert context.match == True
    


    # context.agent_mode = functions.agent_provising('Agent10', 'test', 'true', context.token, iface="mock")
    # assert agent_mode == 'online'
