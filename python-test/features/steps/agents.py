from configs import TestConfig
from utils import *
from local_agent import run_local_agent_container, run_agent_config_file, get_orb_agent_logs, get_logs_and_check
from agent_groups import return_matching_groups, tags_to_match_k_groups
from behave import given, then, step
from hamcrest import *
from datetime import datetime
from agent_config_file import FleetAgent, ConfigFiles
import yaml
from yaml.loader import SafeLoader
import re
import json
import psutil
import os
import signal
from deepdiff import DeepDiff
from concurrent.futures import ThreadPoolExecutor


configs = TestConfig.configs()
agent_name_prefix = "test_agent_name_"
orb_url = configs.get('orb_url')
verify_ssl_bool = eval(configs.get('verify_ssl').title())


@given("that an agent with {orb_tags} orb tag(s) already exists and is {status}")
def check_if_agents_exist(context, orb_tags, status):
    context.agent_name = generate_random_string_with_predefined_prefix(agent_name_prefix)
    context.orb_tags = create_tags_set(orb_tags)
    context.agent = create_agent(context.token, context.agent_name, context.orb_tags)
    context.agent_key = context.agent["key"]
    token = context.token
    run_local_agent_container(context, "available")
    agent_id = context.agent['id']
    existing_agents = get_agent(token, agent_id)
    assert_that(len(existing_agents), greater_than(0), "Agent not created")
    timeout = 30
    logs = get_orb_agent_logs(context.container_id)
    agent_status, context.agent = wait_until_expected_agent_status(token, agent_id, status, timeout=timeout)
    assert_that(agent_status, is_(equal_to(status)),
                f"Agent did not get '{status}' after {str(timeout)} seconds, but was '{agent_status}'. \n"
                f"Agent: {json.dumps(context.agent, indent=4)}. \n Logs: {logs}")
    local_orb_path = configs.get("local_orb_path")
    agent_schema_path = local_orb_path + "/python-test/features/steps/schemas/agent_schema.json"
    is_schema_valid = validate_json(context.agent, agent_schema_path)
    assert_that(is_schema_valid, equal_to(True), f"Invalid agent json. \n Agent = {context.agent}."
                                                 f"Agent logs: {get_orb_agent_logs(context.container_id)}."
                                                 f"\nLogs: {logs}")


@step('a new agent is created with {orb_tags} orb tag(s)')
def agent_is_created(context, orb_tags):
    context.agent_name = generate_random_string_with_predefined_prefix(agent_name_prefix)
    context.orb_tags = create_tags_set(orb_tags)
    context.agent = create_agent(context.token, context.agent_name, context.orb_tags)
    context.agent_key = context.agent["key"]


@step('a new agent is created with orb tags matching {amount_of_group} existing group')
def agent_is_created_matching_group(context, amount_of_group):
    context.agent_name = agent_name_prefix + random_string(10)
    all_used_tags = tags_to_match_k_groups(context.token, amount_of_group, context.agent_groups)
    agent = create_agent(context.token, context.agent_name, all_used_tags)
    context.agent = agent
    context.agent_key = context.agent["key"]


@then('the agent status in Orb should be {status} within {seconds} seconds')
def check_agent_status_within_seconds(context, status, seconds):
    timeout = int(seconds)
    token = context.token
    agent_status, context.agent = wait_until_expected_agent_status(token, context.agent['id'], status, timeout=timeout)
    logs = get_orb_agent_logs(context.container_id)
    assert_that(agent_status, is_(equal_to(status)),
                f"Agent did not get '{status}' after {str(timeout)} seconds, but was '{agent_status}'."
                f"\n Agent: {context.agent}. \nAgent logs: {logs}")


@then('the agent status in Orb should be {status} after {seconds} seconds')
def check_agent_status_after_seconds(context, status, seconds):
    event = threading.Event()
    event.wait(int(seconds))
    event.set()
    if event.is_set() is True:
        token = context.token
        agent_status, context.agent = wait_until_expected_agent_status(token, context.agent['id'], status, timeout=120)
        try:
            logs = get_orb_agent_logs(context.container_id)
        except Exception as e:
            logs = e
        assert_that(agent_status, is_(equal_to(status)),
                    f"Agent did not get '{status}' after {str(seconds)} seconds, but was '{agent_status}'."
                    f"\n Agent: {context.agent}. \nAgent logs: {logs}")


@step('the agent status is {status}')
def check_agent_status(context, status):
    timeout = 30
    token = context.token
    agent_status, context.agent = wait_until_expected_agent_status(token, context.agent['id'], status, timeout=timeout)
    logs = get_orb_agent_logs(context.container_id)
    assert_that(agent_status, is_(equal_to(status)),
                f"Agent did not get '{status}' after {str(timeout)} seconds, but was '{agent_status}'."
                f"Agent: {json.dumps(context.agent, indent=4)}."
                f"Agent logs: {logs}.")


@step("created agent has taps: {taps}")
def verify_agent_taps(context, taps):
    agent = get_agent(context.token, context.agent['id'])
    agent_taps = agent["agent_metadata"]["backends"]["pktvisor"]["data"]["taps"]
    default_taps = {
        "DNSTAP": {
            "default_dnstap": {
                "config": {
                    "tcp": "0.0.0.0:9990"
                },
                "input_type": "dnstap",
                "interface": "visor.module.input/1.0",
                "tags": None
            }},
        "NETFLOW": {
            "default_netflow": {
                "config": {
                    "bind": "0.0.0.0",
                    "flow_type": "netflow",
                    "port": 9995
                },
                "input_type": "flow",
                "interface": "visor.module.input/1.0",
                "tags": None
            }},
        "PCAP": {
            "default_pcap": {
                "config": {
                    "iface": configs.get("orb_agent_interface", "auto")
                },
                "input_type": "pcap",
                "interface": "visor.module.input/1.0",
                "tags": None
            }},
        "SFLOW": {
            "default_sflow": {
                "config": {
                    "bind": "0.0.0.0",
                    "flow_type": "sflow",
                    "port": 9994
                },
                "input_type": "flow",
                "interface": "visor.module.input/1.0",
                "tags": None
            }}
    }

    taps = taps.split(", ")
    for tap in taps:
        default_agent_tap = default_taps.get(tap.upper())
        default_agent_tap_value = list(default_agent_tap.values())[0]
        agent_tap = agent_taps.get(list(default_agent_tap.keys())[0])
        diff = DeepDiff(default_agent_tap_value, agent_tap, exclude_paths={"root['len'], root['tags']"})
        assert_that(diff, equal_to({}), f"Tap {tap} is different from expected one. Agent: {agent}\n"
                                        f"Default tap: {default_agent_tap}")


@then('cleanup agents')
def clean_agents(context):
    """
    Remove all agents starting with 'agent_name_prefix' from the orb

    :param context: Behave class that contains contextual information during the running of tests.
    """
    token = context.token
    agents_list = list_agents(token)
    agents_filtered_list = filter_list_by_parameter_start_with(agents_list, 'name', agent_name_prefix)
    delete_agents(token, agents_filtered_list)


@step("{amount_of_datasets} datasets are linked with each policy on agent's heartbeat within {time_to_wait} seconds")
def multiple_dataset_for_policy(context, amount_of_datasets, time_to_wait):
    datasets_ok, context.agent = check_datasets_for_policy(context.token, context.agent['id'],
                                                           context.list_agent_policies_id,
                                                           amount_of_datasets, timeout=time_to_wait)
    logs = get_orb_agent_logs(context.container_id)
    diff = datasets_ok ^ set(context.list_agent_policies_id)
    assert_that(datasets_ok, equal_to(set(context.list_agent_policies_id)),
                f"Amount of datasets linked with policy {diff} failed. Agent: {context.agent}. \nAgent logs: {logs}")


@step("this agent's heartbeat shows that {amount_of_policies} policies are applied and {amount_of_policies_with_status}"
      " has status {policies_status}")
def list_policies_applied_to_an_agent_and_referred_status(context, amount_of_policies, amount_of_policies_with_status,
                                                          policies_status):
    if amount_of_policies_with_status == "all":
        amount_of_policies_with_status = int(amount_of_policies)
    context.agent, context.list_agent_policies_id, amount_of_policies_applied_with_status = \
        get_policies_applied_to_an_agent_by_status(context.token, context.agent['id'], amount_of_policies,
                                                   amount_of_policies_with_status, policies_status, timeout=180)
    logs = get_orb_agent_logs(context.container_id)
    assert_that(amount_of_policies_applied_with_status, equal_to(int(amount_of_policies_with_status)),
                f"{amount_of_policies_with_status} policies was supposed to have status {policies_status}. \n"
                f"Agent: {context.agent}. \n Logs: {logs}")


@step("the version of policy in agent must be {policy_version}")
def check_agent_version_policy(context, policy_version):
    assert_that(policy_version.isdigit(), equal_to(True), "Policy version must be a number")
    policy_version = int(policy_version)
    policy_id = context.policy.get("id")
    agent, correct_version = wait_policy_version_in_agent(context.token, context.agent.get("id"),
                                                          policy_id, policy_version)
    assert_that(correct_version, equal_to(True), f"Policy {policy_id} version "
                                                 f"{policy_version} was not applied. Agent: {agent}")


@step("this agent's heartbeat shows that {amount_of_policies} policies are applied to the agent")
def list_policies_applied_to_an_agent(context, amount_of_policies):
    context.agent, context.list_agent_policies_id = get_policies_applied_to_an_agent(context.token, context.agent['id'],
                                                                                     amount_of_policies, timeout=180)
    context.agent = get_agent(context.token, context.agent['id'])
    logs = get_orb_agent_logs(context.container_id)
    assert_that(len(context.list_agent_policies_id), equal_to(int(amount_of_policies)),
                f"Amount of policies applied to this agent failed with {len(context.list_agent_policies_id)} policies."
                f"\n Agent: {json.dumps(context.agent, indent=4)}. \n Logs: {logs}")


@step("this agent's heartbeat shows that {amount_of_groups} groups are matching the agent")
def list_groups_matching_an_agent(context, amount_of_groups):
    groups_matching, context.groups_matching_id = return_matching_groups(context.token, context.agent_groups,
                                                                         context.agent)
    context.list_groups_id, context.agent = get_groups_to_which_agent_is_matching(context.token, context.agent['id'],
                                                                                  context.groups_matching_id,
                                                                                  timeout=180)
    logs = get_orb_agent_logs(context.container_id)
    assert_that(len(context.list_groups_id), equal_to(int(amount_of_groups)),
                f"Amount of groups matching the agent failed with {context.list_groups_id} groups. \n"
                f"Agent: {json.dumps(context.agent, indent=4)} \n\n"
                f"Agent Logs: {logs}.")
    assert_that(sorted(context.list_groups_id), equal_to(sorted(context.groups_matching_id)),
                "Groups matching the agent is not the same as the created by test process  \n"
                f"Agent: {json.dumps(context.agent, indent=4)} \n\n"
                f"Agent Logs: {logs}.")


@step("edit the orb tags on agent and use {orb_tags} orb tag(s)")
def editing_agent_tags(context, orb_tags):
    agent = get_agent(context.token, context.agent["id"])
    context.orb_tags = create_tags_set(orb_tags)
    edit_agent(context.token, context.agent["id"], agent["name"], context.orb_tags, expected_status_code=200)
    context.agent = get_agent(context.token, context.agent["id"])


@step("edit the orb tags on agent and use orb tags matching {amount_of_group} existing group")
def agent_is_edited_matching_group(context, amount_of_group):
    all_used_tags = tags_to_match_k_groups(context.token, amount_of_group, context.agent_groups)
    agent = get_agent(context.token, context.agent["id"])
    edit_agent(context.token, agent["id"], agent["name"], all_used_tags, expected_status_code=200)
    context.agent = get_agent(context.token, context.agent["id"])


@step("edit the agent name")
def editing_agent_name(context):
    agent = get_agent(context.token, context.agent["id"])
    agent_new_name = generate_random_string_with_predefined_prefix(agent_name_prefix, 5)
    edit_agent(context.token, context.agent["id"], agent_new_name, agent['orb_tags'], expected_status_code=200)
    context.agent = get_agent(context.token, context.agent["id"])
    assert_that(context.agent["name"], equal_to(agent_new_name), f"Agent name editing failed. Agent: {context.agent}")


@step("edit the agent name and edit orb tags on agent using {orb_tags} orb tag(s)")
def editing_agent_name_and_tags(context, orb_tags):
    agent_new_name = generate_random_string_with_predefined_prefix(agent_name_prefix, 5)
    context.orb_tags = create_tags_set(orb_tags)
    edit_agent(context.token, context.agent["id"], agent_new_name, context.orb_tags, expected_status_code=200)
    context.agent = get_agent(context.token, context.agent["id"])
    assert_that(context.agent["name"], equal_to(agent_new_name), f"Agent name editing failed. Agent: {context.agent}")
    assert_that(context.agent['orb_tags'], equal_to(context.orb_tags), f"Agent orb tags editing failed."
                                                                       f" Agent{context.agent}")


@step("agent must have {amount_of_tags} tags")
def check_agent_tags(context, amount_of_tags):
    agent = get_agent(context.token, context.agent["id"])
    assert_that(len(dict(agent["orb_tags"])), equal_to(int(amount_of_tags)), f"Amount of orb tags failed. "
                                                                             f"Agent: {agent}")


@then("remove all the agents .yaml generated on test process")
def remove_agent_config_files(context):
    dir_path = configs.get("local_orb_path")
    all_files_generated = find_files(agent_name_prefix, ".yaml", dir_path)
    if len(all_files_generated) > 0:
        for file in all_files_generated:
            os.remove(file)


@then("remove the agent .yaml generated on each scenario")
def remove_one_agent_config_files(context):
    dir_path = configs.get("local_orb_path")
    all_files_generated = find_files(context.agent_file_name, ".yaml", dir_path)
    if len(all_files_generated) > 0:
        for file in all_files_generated:
            os.remove(file)


@step("this agent is removed")
def remove_orb_agent(context):
    delete_agent(context.token, context.agent['id'])
    get_agent(context.token, context.agent['id'], 404)


@threading_wait_until
def check_agent_exists_on_backend(token, agent_name, event=None):
    agent = None
    all_agents = list_agents(token)
    for agent in all_agents:
        if agent_name == agent['name']:
            event.set()
            return agent, event.is_set()
    return agent, event.is_set()


@step("an agent(backend_type:{backend_type}, settings: {settings}) is {provision} via a configuration file on "
      "port {port} with {agent_tags} agent tags and has status {status}. [Overwrite default: {overwrite_default}. "
      "Paste only file: {paste_only_file}. Use specific backend config {specific_backend_config}]")
def provision_agent_using_config_file_drop_pkt_config(context, backend_type, settings, provision, port, agent_tags,
                                                      status, overwrite_default, paste_only_file,
                                                      specific_backend_config):
    specific_backend_config = json.loads(specific_backend_config)
    provision_agent_using_config_file(context, backend_type, settings, provision, port, agent_tags, status,
                                      overwrite_default, paste_only_file,
                                      specific_backend_config=specific_backend_config)


@step("an agent(backend_type:{backend_type}, settings: {settings}) is {provision} via a configuration file on port {"
      "port} with {agent_tags} agent tags and has status {status}. [Overwrite default: {overwrite_default}. Paste only "
      "file: {paste_only_file}]")
def provision_agent_using_config_file(context, backend_type, settings, provision, port, agent_tags, status,
                                      overwrite_default, paste_only_file, **kwargs):
    assert_that(provision, any_of(equal_to("self-provisioned"), equal_to("provisioned")), "Unexpected provision "
                                                                                          "attribute")
    overwrite_default = overwrite_default.title()
    paste_only_file = paste_only_file.title()
    assert_that(overwrite_default, any_of("True", "False"), "Unexpected value for overwrite_default parameter.")
    assert_that(paste_only_file, any_of("True", "False"), "Unexpected value for overwrite_default parameter.")
    overwrite_default = eval(overwrite_default)
    paste_only_file = eval(paste_only_file)
    settings_is_json, settings = is_json(settings)
    assert_that(settings_is_json, is_(True), f"settings must be written in json format. Current settings: "
                                             f"{settings}")
    if ("tcp" in settings.keys() and settings["tcp"].split(":")[1] == "available_port") or (
            "port" in settings.keys() and settings["port"] == "available_port"):
        port_to_attach = return_port_by_availability(context)
        if "tcp" in settings.keys():
            ip = settings["tcp"].split(":")[0]
            tcp = f"{ip}:{port_to_attach}"
            settings["tcp"] = tcp
        else:
            settings["port"] = port_to_attach
    if "port" in settings.keys() and settings["port"] == "switch":
        settings["port"] = context.switch_port
    if provision == "provisioned":
        auto_provision = "false"
        orb_cloud_mqtt_id = context.agent['id']
        orb_cloud_mqtt_key = context.agent['key']
        orb_cloud_mqtt_channel_id = context.agent['channel_id']
        agent_name = context.agent['name']
    else:
        auto_provision = "true"
        orb_cloud_mqtt_id = None
        orb_cloud_mqtt_key = None
        orb_cloud_mqtt_channel_id = None
        agent_name = f"{agent_name_prefix}{random_string(10)}"

    if "iface" in settings.keys() and isinstance(settings["iface"], str) and settings["iface"].lower() == "mocked":
        interface = context.mock_iface_name
    else:
        interface = configs.get('orb_agent_interface', 'auto')
    settings["iface"] = interface
    orb_url = configs.get('orb_url')
    context.port = return_port_by_availability(context, True)
    if "tap_name" in context:
        tap_name = context.tap_name
    else:
        tap_name = agent_name
    settings["tap_name"] = tap_name
    context.agent_file_name, tags_on_agent, context.tap, safe_config_file = \
        create_agent_config_file(backend_type, context.token, agent_name, agent_tags, orb_url,
                                 context.port,
                                 context.agent_groups, auto_provision,
                                 orb_cloud_mqtt_id, orb_cloud_mqtt_key, orb_cloud_mqtt_channel_id, settings,
                                 overwrite_default, paste_only_file, **kwargs)
    if backend_type == "pktvisor":
        for key, value in context.tap.items():
            if 'tags' in value.keys():
                context.tap_tags.update({key: value['tags']})
    else:
        context.tap_tags = {}
    context.container_id = run_agent_config_file(context.agent_file_name, overwrite_default, paste_only_file)
    if context.container_id not in context.containers_id.keys():
        context.containers_id[context.container_id] = str(context.port)
    if backend_type == "pktvisor":
        log_message = f"web server listening on localhost:{context.port}"
        agent_started, logs, log_line = get_logs_and_check(context.container_id, log_message, element_to_check="log")
        assert_that(agent_started, equal_to(True), f"Log {log_message} not found on agent logs."
                                                   f" Agent Name: {agent_name}. Logs:{logs}")
    else:
        log_message = f"Started receiver for OTLP in orb-agent"
        agent_started, logs, log_line = get_logs_and_check(context.container_id, log_message, element_to_check="msg")
        assert_that(agent_started, equal_to(True), f"Log {log_message} not found on agent logs."
                                                   f" Agent Name: {agent_name}. Logs:{logs}")
        assert_that(str(log_line.get("port", "")), equal_to(str(context.port)), f"Log {log_message} related to port "
                                                                                f"{context.port} not found on agent "
                                                                                f"logs. Logs: {logs}")

    context.agent, is_agent_created = check_agent_exists_on_backend(context.token, agent_name, timeout=60)
    logs = get_orb_agent_logs(context.container_id)
    assert_that(is_agent_created, equal_to(True), f"Agent {agent_name} not found in /agents route."
                                                  f"\n Config File (json converted): {safe_config_file}."
                                                  f"\nLogs: {logs}.")
    context.agent, are_tags_correct = get_agent_tags(context.token, context.agent['id'], tags_on_agent)
    assert_that(are_tags_correct, equal_to(True), f"Agent tags created does not match with the required ones. Agent:"
                                                  f"{context.agent}. Tags that would be present: {tags_on_agent}.\n"
                                                  f"Agent Logs: {logs}")
    assert_that(context.agent, is_not(None), f"Agent {agent_name} not correctly created. Logs: {logs}")
    agent_id = context.agent['id']
    existing_agents = get_agent(context.token, agent_id)
    assert_that(len(existing_agents), greater_than(0), f"Agent not created. Logs: {logs}")
    agent_status, context.agent = wait_until_expected_agent_status(context.token, agent_id, status)
    assert_that(agent_status, is_(equal_to(status)),
                f"Agent did not get '{status}' after 30 seconds, but was '{agent_status}'. \n"
                f"Agent: {json.dumps(context.agent, indent=4)}. \n Logs: {logs}")


@step("remotely restart the agent")
def reset_agent_remotely(context):
    context.considered_timestamp_reset = datetime.now().timestamp()
    status_code, repsonse = (
        return_api_post_response(f"{orb_url}/api/v1/agents/{context.agent['id']}/rpc/reset",
                                 token=context.token, verify=verify_ssl_bool))
    logs = get_orb_agent_logs(context.container_id)
    assert_that(status_code, equal_to(200),
                f"Request to restart agent failed with status= {str(status_code)}. \n Agent: {context.agent}\n"
                f" Logs: {logs}")


@step("{route} route must be enabled")
def check_agent_backend_pktvisor_routes(context, route):
    assert_that(route, any_of(equal_to("taps"), equal_to("handlers"), equal_to("inputs"), equal_to("backends")),
                "Invalid agent route")

    agent_backend_routes = {"backends": "backends", "taps": "backends/pktvisor/taps",
                            "inputs": "backends/pktvisor/inputs",
                            "handlers": "backends/pktvisor/handlers"}

    status_code, response = return_api_get_response(f"{orb_url}/api/v1/agents/{agent_backend_routes[route]}",
                                                    token=context.token, verify=verify_ssl_bool)
    assert_that(status_code, equal_to(200),
                f"Request to get {route} route failed with status =" + str(status_code))
    local_orb_path = configs.get("local_orb_path")
    route_schema_path = local_orb_path + f"/python-test/features/steps/schemas/{route}_schema.json"
    is_schema_valid = validate_json(response, route_schema_path)
    assert_that(is_schema_valid, equal_to(True), f"Invalid route json. \n Route = {route}")


@step("edit the agent name using an already existent one")
def edit_agent_using_name_with_conflict(context):
    agents_list = list_agents(context.token)
    agents_filtered_list = filter_list_by_parameter_start_with(agents_list, 'name', agent_name_prefix)
    agents_name = list()
    for agent in agents_filtered_list:
        agents_name.append(agent['name'])
    agents_name.remove(context.agent['name'])
    name_to_use = random.choice(agents_name)
    context.error_message = edit_agent(context.token, context.agent['id'], name_to_use, context.agent['orb_tags'], 409)


@step("a new agent is requested to be created with the same name as an existent one")
def create_agent_with_name_conflict(context):
    tag = create_tags_set('1')
    context.error_message = create_agent(context.token, context.agent['name'], tag, 409)


@step("the error message on response is {message}")
def check_error_message(context, message):
    assert_that(context.error_message['error'], equal_to(message), "Unexpected error message")


@step("agent backend (pktvisor) stops running")
def kill_pktvisor_on_agent(context):
    try:
        current_proc_pid = None
        for process in psutil.process_iter():
            if "pkt" in process.name():
                proc_con = process.connections()
                proc_port = proc_con[0].laddr.port
                if proc_port == context.port:
                    current_proc_pid = process.pid
                    process.send_signal(signal.SIGKILL)
                    break
        assert_that(current_proc_pid, is_not(None), "Unable to find pid of pktvisor process")
    except psutil.AccessDenied:
        context.access_denied = True
        context.scenario.skip(f"You are not allowed to run this scenario without root permissions.")
    except Exception as exception:
        raise exception


@step("{backend} state is {state}")
def check_back_state(context, backend, state):
    backend_state, agent = wait_until_expected_backend_state(context.token, context.agent['id'], backend, state,
                                                             timeout=180)
    logs = get_orb_agent_logs(context.container_id)
    assert_that(backend_state, equal_to(state), f"Unexpected backend state on agent: {agent}. Logs: {logs}")


@step("{backend} error is {error}")
def check_back_error(context, backend, error):
    backend_state, agent = wait_until_expected_backend_error(context.token, context.agent['id'], backend, error,
                                                             timeout=180)
    logs = get_orb_agent_logs(context.container_id)
    assert_that(backend_state, equal_to(error), f"Unexpected backend error on agent: {agent}. Logs: {logs}")


@step("agent backend {backend} restart_count is {restart_count}")
def check_auto_reset(context, backend, restart_count):
    amount, agent = wait_until_expected_amount_of_restart_count(context.token, context.agent['id'], backend,
                                                                restart_count, timeout=400)
    logs = get_orb_agent_logs(context.container_id)
    assert_that(int(amount), equal_to(int(restart_count)), f"Unexpected restart count for backend {backend} on agent: "
                                                           f"{agent}. Logs: {logs}")


@threading_wait_until
def wait_until_expected_agent_status(token, agent_id, status, event=None):
    """
    Keeps fetching agent data from Orb control plane until it gets to
    the expected agent status or this operation times out

    :param (str) token: used for API authentication
    :param (str) agent_id: whose status will be evaluated
    :param (str) status: expected agent status
    :param (obj) event: threading.event
    """

    agent = get_agent(token, agent_id)
    agent_status = agent['state']
    if agent_status == status:
        event.set()
        return agent_status, agent
    return agent_status, agent


@threading_wait_until
def wait_until_expected_amount_of_restart_count(token, agent_id, backend, amount_of_restart, event=None):
    """
    Keeps fetching agent data from Orb control plane until it gets to
    the expected agent status or this operation times out

    :param (str) token: used for API authentication
    :param (str) agent_id: whose backend state will be evaluated
    :param (str) amount_of_restart: expected amount of restart state
    :param (str) backend: backend to check state
    :param (obj) event: threading.event
    """

    agent = get_agent(token, agent_id)
    if 'restart_count' in agent["last_hb_data"]["backend_state"][backend].keys():
        amount = agent["last_hb_data"]["backend_state"][backend]['restart_count']
        if int(amount_of_restart) == int(amount):
            event.set()
            return amount, agent
        return amount, agent
    else:
        return None, agent


@threading_wait_until
def wait_until_expected_backend_state(token, agent_id, backend, state, event=None):
    """
    Keeps fetching agent data from Orb control plane until it gets to
    the expected agent status or this operation times out

    :param (str) token: used for API authentication
    :param (str) agent_id: whose backend state will be evaluated
    :param (str) state: expected backend state
    :param (str) backend: backend to check state
    :param (obj) event: threading.event
    """

    backend_state, agent = get_backend_info(token, agent_id, backend, "state")
    if backend_state == state:
        event.set()
        return backend_state, agent
    return backend_state, agent


@threading_wait_until
def wait_until_expected_backend_error(token, agent_id, backend, error, event=None):
    """
    Keeps fetching agent data from Orb control plane until it gets to
    the expected agent status or this operation times out

    :param (str) token: used for API authentication
    :param (str) agent_id: whose backend error will be evaluated
    :param (str) error: expected backend error
    :param (str) backend: backend to check error
    :param (obj) event: threading.event
    """

    backend_error, agent = get_backend_info(token, agent_id, backend, "error")
    if backend_error == error:
        event.set()
        return backend_error, agent
    else:
        return None, agent


def get_backend_info(token, agent_id, backend, info):
    """
    Get the backend state for a specific agent.

    :param str token: Access token for authentication
    :param str agent_id: ID of the agent
    :param str backend: Name of the agent backend
    :param str info: Info requested
    :return: The agent backend required info
    :rtype: str or None
    """
    agent = get_agent(token, agent_id)
    try:
        last_hb_data = agent.get("last_hb_data", {})
        backend_info = last_hb_data.get("backend_state", {}).get(backend, {}).get(info)
        return backend_info, agent
    except Exception as e:
        log.error(f"Error getting backend {info}: {e}. Agent: {agent}")
        return None, agent


def get_agent(token, agent_id, expected_status_code=200):
    """
    Gets an agent from Orb control plane

    :param (str) token: used for API authentication
    :param (str) agent_id: that identifies agent to be fetched
    :param (int) expected_status_code: status code that must be returned on response
    :returns: (dict) the fetched agent
    """
    status_code, response = return_api_get_response(f"{orb_url}/api/v1/agents/{agent_id}", token=token,
                                                    verify=verify_ssl_bool)
    assert_that(status_code, equal_to(expected_status_code),
                f"Request to get agent id= {agent_id} failed with status= {status_code}:"
                f"{response}")

    return response


def list_agents(token, limit=100, offset=0):
    """
    Lists all agents from Orb control plane that belong to this user

    :param (str) token: used for API authentication
    :param (int) limit: Size of the subset to retrieve. (max 100). Default = 100
    :param (int) offset: Number of items to skip during retrieval. Default = 0.
    :returns: (list) a list of agents
    """

    all_agents, total, offset = list_up_to_limit_agents(token, limit, offset)

    new_offset = limit + offset

    while new_offset < total:
        agents_from_offset, total, offset = list_up_to_limit_agents(token, limit, new_offset)
        all_agents = all_agents + agents_from_offset
        new_offset = limit + offset

    return all_agents


def list_up_to_limit_agents(token, limit=100, offset=0):
    """
    Lists up to 100 agents from Orb control plane that belong to this user

    :param (str) token: used for API authentication
    :param (int) limit: Size of the subset to retrieve. (max 100). Default = 100
    :param (int) offset: Number of items to skip during retrieval. Default = 0.
    :returns: (list) a list of agents, (int) total agents on orb, (int) offset
    """

    status_code, response = return_api_get_response(f"{orb_url}/api/v1/agents", token=token,
                                                    params={"limit": limit, "offset": offset}, verify=verify_ssl_bool)

    assert_that(status_code, equal_to(200),
                f"Request to list agents failed with status= {str(status_code)}:{str(response)}")
    return response['agents'], response['total'], response['offset']


def delete_agents(token, list_of_agents):
    """
    Deletes from Orb control plane the agents specified on the given list

    :param (str) token: used for API authentication
    :param (list) list_of_agents: that will be deleted
    """
    log.debug(f"Deleting {len(list_of_agents)} agents")
    with ThreadPoolExecutor() as executor:
        futures = [executor.submit(delete_agent, token, agent.get('id')) for agent in list_of_agents]
        results = [future.result() for future in futures]
    log.debug(f"Finishing deleting agents")
    return results


def delete_agent(token, agent_id):
    """
    Deletes an agent from Orb control plane

    :param (str) token: used for API authentication
    :param (str) agent_id: that identifies the agent to be deleted
    """
    status_code, response = return_api_delete_response(f"{orb_url}/api/v1/agents/{agent_id}", token=token,
                                                       verify=verify_ssl_bool)

    assert_that(status_code, equal_to(204), f"Request to delete agent id= {agent_id} failed with status= {status_code}")


@threading_wait_until
def wait_until_agent_being_created(token, name, tags, expected_status_code=201, event=None):
    json_request = {"name": name, "orb_tags": tags, "validate_only": False}

    status_code, response = return_api_post_response(f"{orb_url}/api/v1/agents", token=token,
                                                     request_body=json_request,
                                                     verify=verify_ssl_bool)
    if status_code == expected_status_code:
        event.set()
        return status_code, response
    return status_code, response


def create_agent(token, name, tags, expected_status_code=201):
    """
    Creates an agent in Orb control plane

    :param (str) token: used for API authentication
    :param (str) name: of the agent to be created
    :param (dict) tags: orb agent tags
    :param expected_status_code: status code to be returned on response
    :returns: (dict) a dictionary containing the created agent data
    """
    status_code, response = wait_until_agent_being_created(token, name, tags, expected_status_code)
    assert_that(status_code, equal_to(expected_status_code),
                'Request to create agent failed with status=' + str(status_code) + ":" + str(response))

    return response


def edit_agent(token, agent_id, name, tags, expected_status_code=200):
    """
    :param (str) token: used for API authentication
    :param (str) agent_id: that identifies the agent to be deleted
    :param (str) name: of the agent to be created
    :param (dict) tags: orb agent tags
    :param (int) expected_status_code: expected request's status code. Default:200 (happy path).
    :return: (dict) a dictionary containing the edited agent data
    """

    json_request = {"name": name, "orb_tags": tags, "validate_only": False}
    status_code, response = return_api_put_response(orb_url + '/api/v1/agents/' + agent_id, request_body=json_request,
                                                    token=token, verify=verify_ssl_bool)

    assert_that(status_code, equal_to(expected_status_code),
                'Request to edit agent failed with status=' + str(status_code) + ":" + str(response))

    return response


@threading_wait_until
def get_policies_applied_to_an_agent(token, agent_id, amount_of_policies, event=None):
    """

    :param (str) token: used for API authentication
    :param (str) agent_id: that identifies the agent to be deleted
    :param (int) amount_of_policies: amount of policies that is expected to be applied to the agents
    :param (obj) event: threading.event
    :return:  (dict) agent -> the fetched agent and (list) list_agent_policies_id -> list with the ids of the policies
    that are applied to the agent
    """
    list_agent_policies_id = list()
    agent = get_agent(token, agent_id)
    if 'policy_state' in agent['last_hb_data'].keys():
        list_agent_policies_id = list(agent['last_hb_data']['policy_state'].keys())
        if len(list_agent_policies_id) == int(amount_of_policies):
            event.set()
            return agent, list_agent_policies_id
    return agent, list_agent_policies_id


@threading_wait_until
def get_policies_applied_to_an_agent_by_status(token, agent_id, amount_of_policies, amount_of_policies_with_status,
                                               status, event=None):
    agent, list_agent_policies_id = get_policies_applied_to_an_agent(token, agent_id, amount_of_policies, timeout=180)
    list_of_policies_status = list()
    for policy_id in list_agent_policies_id:
        list_of_policies_status.append(agent['last_hb_data']['policy_state'][policy_id]["state"])
    if amount_of_policies_with_status == "all":
        amount_of_policies_with_status = int(amount_of_policies)
    amount_of_policies_applied_with_status = list_of_policies_status.count(status)
    if amount_of_policies_applied_with_status == amount_of_policies_with_status:
        event.set()
    else:
        event.wait(5)
    return agent, list_agent_policies_id, amount_of_policies_applied_with_status


@threading_wait_until
def get_groups_to_which_agent_is_matching(token, agent_id, groups_matching_ids, event=None):
    """

    :param (str) token: used for API authentication
    :param (str) agent_id: that identifies the agent to be deleted
    :param (list) groups_matching_ids: list with the ids of the groups to with the agent should be subscribed
    :param (obj) event: threading.event
    :return: (list) list_groups_id -> list with the ids of the groups to with the agent is subscribed
    """
    list_groups_id = list()
    agent = get_agent(token, agent_id)
    if 'group_state' in agent['last_hb_data'].keys():
        list_groups_id = list(agent['last_hb_data']['group_state'].keys())
        if sorted(list_groups_id) == sorted(groups_matching_ids):
            event.set()
            return list_groups_id, agent
    return list_groups_id, agent


def create_agent_config_file(backend_type, token, agent_name, agent_tags, orb_url, port,
                             existing_agent_groups, auto_provision="true",
                             orb_cloud_mqtt_id=None, orb_cloud_mqtt_key=None, orb_cloud_mqtt_channel_id=None,
                             settings=None, overwrite_default=False, only_file=False, **kwargs):
    """
    Create a file .yaml with configs of the agent that will be provisioned

    :param backend_type: agent backend type
    :param (str) token: used for API authentication
    :param (str) agent_name: name of the agent that will be created
    :param (str) agent_tags: agent tags
    :param orb_url: url address of ORB
    :param (str) port: port on which agent must run.
    :param (dict) existing_agent_groups: all agent groups available
    :param (str) auto_provision: if true auto_provision the agent. If false, provision an agent already existent on orb
    :param (str) orb_cloud_mqtt_id: agent mqtt id.
    :param (str) orb_cloud_mqtt_key: agent mqtt key.
    :param (str) orb_cloud_mqtt_channel_id: agent mqtt channel id.
    :param (dict) settings: settings of input
    :param (bool) overwrite_default: if True and only_file is False saves the agent as "agent.yaml". Else, save it with
    agent name
    :param (bool) only_file: is true copy only the file. If false, copy the directory.
    :return: path to the directory where the agent config file was created
    """
    assert_that(auto_provision, any_of(equal_to("true"), equal_to("false")), "Unexpected value for auto_provision "
                                                                             "on agent config file creation")
    verify_ssl_str = configs.get('verify_ssl')
    verify_ssl = str_to_bool(verify_ssl_str)
    auto_provision = str_to_bool(auto_provision)

    if re.match(r"matching (\d+|all|the) group*", agent_tags):
        amount_of_group = re.search(r"(\d+|all|the)", agent_tags).groups()[0]
        all_used_tags = tags_to_match_k_groups(token, amount_of_group, existing_agent_groups)
        tags = {"tags": all_used_tags}
    else:
        tags = {"tags": create_tags_set(agent_tags)}
    mqtt_url = configs.get('mqtt_url')
    cloud_settings = {"orb_cloud_mqtt_id": orb_cloud_mqtt_id,
                      "orb_cloud_mqtt_key": orb_cloud_mqtt_key,
                      "orb_cloud_mqtt_channel_id": orb_cloud_mqtt_channel_id,
                      "name": agent_name,
                      "token": token}
    if backend_type == "pktvisor":
        assert_that(settings, has_key("iface"), f"Missing iface on agent config file creation")
        assert_that(settings, has_key("tap_name"), f"Missing tap_name on agent config file creation")
        tap_name = settings.get('tap_name')
        settings.pop('tap_name', None)
        input_type = settings.get('input_type', 'pcap')
        settings.pop('input_type', None)
        input_tags = settings.get('input_tags', '3')
        settings.pop('input_tags', None)
        backend_file, taps = ConfigFiles(backend_type).pktvisor_config_file(tap_name, input_type, input_tags,
                                                                            settings)
        if "specific_backend_config" in kwargs.keys() and "binary" in kwargs["specific_backend_config"].keys():
            binary = kwargs["specific_backend_config"]["binary"]
            if isinstance(binary, str) and binary.lower() == "none":
                binary = None
        else:
            binary = "/usr/local/sbin/pktvisord"
        backend_settings = {"binary": binary}
        backend_settings = remove_empty_from_json(backend_settings)
    elif backend_type == "otel":
        backend_file = ConfigFiles(backend_type).otel_config_file()
        taps = {}
        backend_settings = {}
    else:
        raise Exception(f"Unexpected value for backend_type on agent config file creation: {backend_type}")
    if "specific_backend_config" in kwargs.keys() and "config_file" in kwargs["specific_backend_config"].keys():
        config_file = kwargs["specific_backend_config"]["config_file"]
        if isinstance(config_file, str) and config_file.lower() == "none":
            config_file = None
    else:
        config_file = "default"

    agent_config_file = FleetAgent.config_file_of_orb_agent(backend_type, auto_provision,
                                                            f"{port}", backend_settings, cloud_settings,
                                                            orb_url, mqtt_url, backend_file=backend_file,
                                                            tls_verify=verify_ssl, overwrite_default=overwrite_default,
                                                            config_file=config_file)
    agent_config_file = yaml.load(agent_config_file, Loader=SafeLoader)
    agent_config_file['orb'].update(tags)
    agent_config_file_yaml = yaml.dump(agent_config_file)
    safe_agent_config_file = agent_config_file.copy()
    if "token" in safe_agent_config_file['orb']['cloud']['api'].keys():
        safe_agent_config_file['orb']['cloud']['api']['token'] = "token omitted for security reason"
    log.debug(f"Agent file: {safe_agent_config_file}")
    dir_path = configs.get("local_orb_path")
    if overwrite_default is True and only_file is False:
        agent_name = "agent"
    with open(f"{dir_path}/{agent_name}.yaml", "w+") as f:
        f.write(agent_config_file_yaml)
    return agent_name, tags, taps, safe_agent_config_file


def create_agent_with_otel_backend_config_file(token, agent_name, agent_tags, orb_url, port,
                                               existing_agent_groups, auto_provision="true",
                                               orb_cloud_mqtt_id=None, orb_cloud_mqtt_key=None,
                                               orb_cloud_mqtt_channel_id=None,
                                               overwrite_default=False, only_file=False, otel_config_file=None):
    """
    Create a file .yaml with configs of the agent that will be provisioned

    :param (str) token: used for API authentication
    :param (str) agent_name: name of the agent that will be created
    :param (str) agent_tags: agent tags
    :param (str) orb_url: entire orb url
    :param (str) port: port on which agent must run.
    :param (dict) existing_agent_groups: all agent groups available
    :param (str) auto_provision: if true auto_provision the agent. If false, provision an agent already existent on orb
    :param (str) orb_cloud_mqtt_id: agent mqtt id.
    :param (str) orb_cloud_mqtt_key: agent mqtt key.
    :param (str) orb_cloud_mqtt_channel_id: agent mqtt channel id.
    :param (bool) overwrite_default: if True and only_file is False saves the agent as "agent.yaml". Else, save it with
    agent name
    :param (bool) only_file: is true copy only the file. If false, copy the directory.
    :param (str) otel_config_file: path to otel config file.
    :return: path to the directory where the agent config file was created
    """
    assert_that(auto_provision, any_of(equal_to("true"), equal_to("false")), "Unexpected value for auto_provision "
                                                                             "on agent config file creation")
    convert_type = {"true": True, "false": False}
    auto_provision = convert_type[auto_provision]

    if re.match(r"matching (\d+|all|the) group*", agent_tags):
        amount_of_group = re.search(r"(\d+|all|the)", agent_tags).groups()[0]
        all_used_tags = tags_to_match_k_groups(token, amount_of_group, existing_agent_groups)
        tags = {"tags": all_used_tags}
    else:
        tags = {"tags": create_tags_set(agent_tags)}
    mqtt_url = configs.get('mqtt_url')
    if configs.get('verify_ssl') == 'false':
        agent_config_file = (
            FleetAgent.config_file_of_orb_agent_with_otel_backend(agent_name, token, orb_url,
                                                                  mqtt_url, tls_verify=True,
                                                                  auto_provision=auto_provision,
                                                                  orb_cloud_mqtt_id=orb_cloud_mqtt_id,
                                                                  orb_cloud_mqtt_key=orb_cloud_mqtt_key,
                                                                  orb_cloud_mqtt_channel_id=orb_cloud_mqtt_channel_id,
                                                                  overwrite_default=overwrite_default))
    else:
        agent_config_file = (
            FleetAgent.config_file_of_orb_agent_with_otel_backend(agent_name, token, orb_url,
                                                                  mqtt_url, tls_verify=True,
                                                                  auto_provision=auto_provision,
                                                                  orb_cloud_mqtt_id=orb_cloud_mqtt_id,
                                                                  orb_cloud_mqtt_key=orb_cloud_mqtt_key,
                                                                  orb_cloud_mqtt_channel_id=orb_cloud_mqtt_channel_id,
                                                                  overwrite_default=overwrite_default))
    agent_config_file = yaml.load(agent_config_file, Loader=SafeLoader)
    if otel_config_file is None or otel_config_file == "None":
        agent_config_file['orb']['backends']['otel'].pop("config_file", None)
    elif otel_config_file == "default":
        pass
    else:
        agent_config_file['orb']['backends']['otel']["config_file"] = otel_config_file
    agent_config_file['orb'].update(tags)
    agent_config_file['orb']['backends']['otel'].update({"otlp_port": int(f"{port}")})
    agent_config_file_yaml = yaml.dump(agent_config_file)
    safe_agent_config_file = agent_config_file.copy()
    if "token" in safe_agent_config_file['orb']['cloud']['api'].keys():
        safe_agent_config_file['orb']['cloud']['api']['token'] = "token omitted for security reason"
    dir_path = configs.get("local_orb_path")
    if overwrite_default is True and only_file is False:
        agent_name = "agent"
    with open(f"{dir_path}/{agent_name}.yaml", "w+") as f:
        f.write(agent_config_file_yaml)
    return agent_name, tags, safe_agent_config_file


@threading_wait_until
def check_datasets_for_policy(token, agent_id, list_agent_policies_id, amount_of_datasets, event=None):
    """

    :param (str) token: used for API authentication
    :param (str) agent_id: that identifies the agent to be checked
    :param (list) list_agent_policies_id: list containing all policy ids created by the scenario
    :param (str) amount_of_datasets: amount of dataset that is expected to be applied to each policy
    :param (obj) event: threading.event
    :return: (dict) the set of policies with correct amounts of datasets, (dict) agent data

    """
    dataset_ok = set()
    agent = get_agent(token, agent_id)
    for policy_id in list_agent_policies_id:
        if len(agent['last_hb_data']['policy_state'][policy_id]['datasets']) == int(amount_of_datasets):
            dataset_ok.add(policy_id)
    if len(dataset_ok) == len(list_agent_policies_id):
        event.set()
        return dataset_ok, agent
    return dataset_ok, agent


@threading_wait_until
def get_agent_tags(token, agent_id, expected_tags, event=None):
    """

    :param (str) token: used for API authentication
    :param (str) agent_id: that identifies the agent to be checked
    :param (dict) expected_tags: agent tags expected to be on agent
    :param (obj) event: threading.event
    :return: agent data, if the tags were found
    """
    agent = get_agent(token, agent_id)
    expected_tags_insensitive = {k.lower(): v for k, v in expected_tags['tags'].items()}
    if set(agent['agent_tags']) == set(expected_tags_insensitive):
        event.set()
    else:
        event.wait(1)
    return agent, event.is_set()


@threading_wait_until
def wait_policy_version_in_agent(token, agent_id, policy_id, policy_version, event=None):
    """

    :param (str) token: used for API authentication
    :param (str) agent_id: that identifies the agent to be checked
    :param (str) policy_id: that identifies the policy to be checked
    :param (int) policy_version: that identifies the version to be checked
    :param (obj) event: threading.event
    :return: agent data, if the version was found
    """
    agent = get_agent(token, agent_id)
    if agent.get("last_hb_data", '').get("policy_state", '').get(policy_id, '').get("version", '') == policy_version:
        event.set()
    return agent, event.is_set()
