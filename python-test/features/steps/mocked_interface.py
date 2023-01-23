from hamcrest import *
from behave import given, then, step
from test_config import TestConfig
from utils import random_string, send_terminal_commands, return_port_by_availability
import os
import json

configs = TestConfig.configs()
orb_url = configs.get('orb_url')
verify_ssl_bool = eval(configs.get('verify_ssl').title())


@given("a mocked interface is configured with mtu: {mtu} and ip: {ip}")
def set_mock_interface(context, mtu, ip):
    context.mock_iface_name = f"dummyqa{random_string(8)}"
    return_0 = send_terminal_commands(f"ip link add {context.mock_iface_name} type dummy")
    return_1 = send_terminal_commands(f"ip link set {context.mock_iface_name} up")
    return_2 = send_terminal_commands(f"ifconfig {context.mock_iface_name} mtu {mtu}")
    if "Operation not permitted" in return_0[1]:
        context.access_denied = True
        context.scenario.skip('Root privileges are required')
    assert_that(return_1, equal_to(('', '')), f"Error while set mocked iface up")
    assert_that(return_2, equal_to(('', '')), f"Error while configuring mtu of dummy iface")
    if ip.lower() != "none":
        return_4 = send_terminal_commands(f"ip addr add {ip} dev {context.mock_iface_name}")
        assert_that(return_4, equal_to(('', '')), f"Error while configuring ip of dummy iface")


@given("a virtual switch is configured and is up with type {switch_type} and target {target}")
def set_virtual_switch(context, switch_type, target):
    if "available" in target:
        context.switch_port = return_port_by_availability(context)
        target = target.replace("available", f"{context.switch_port}")
    context.virtual_switch_name = f"brqa{random_string(10)}"
    return_1 = send_terminal_commands(f"ovs-vsctl add-br {context.virtual_switch_name}")
    assert_that(return_1, equal_to(('', '')), "error to add virtual switch")
    return_2 = send_terminal_commands(f"ip link set {context.virtual_switch_name} up")
    assert_that(return_2, equal_to(('', '')), "error to set virtual switch up")
    assert_that(switch_type.lower(), any_of("netflow", "sflow"), f"Invalid switch type")
    if switch_type.lower() == "netflow":
        command = f"ovs-vsctl set bridge {context.virtual_switch_name} netflow=@netflow -- --id=@netflow create NetFlow targets={target}"
    else:
        command = f"ovs-vsctl -- --id=@sflow create sflow agent={context.mock_iface_name} target=\"{target}\" header=128 sampling=64 polling=10 -- set bridge {context.virtual_switch_name} sflow=@sflow"
    return_3 = send_terminal_commands(command)
    assert_that(return_3[1], equal_to(''), f"error to configure {switch_type} switch")


@then("remove dummy interface")
def remove_dummy(context):
    send_terminal_commands(f"ip link set {context.mock_iface_name} down")
    send_terminal_commands(f"ip link delete {context.mock_iface_name} type dummy")


@then("remove virtual switch")
def remove_switch(context):
    if "mock_iface_name" in context:
        send_terminal_commands(f"ovs-vsctl del-br {context.mock_iface_name}")


@step("run mocked data {files_names} on the created virtual {virtual_env}")
def run_mocked_data(context, files_names, virtual_env):
    assert_that(virtual_env, any_of("interface", "switch"), "Invalid env to run mocked data")
    if virtual_env == "interface":
        where_to_run = context.mock_iface_name
    elif virtual_env == "switch":
        where_to_run = context.virtual_switch_name
    else:
        raise "Invalid env to run mocked data"
    cwd = os.getcwd()
    dir_path = os.path.dirname(cwd)
    directory_of_network_data_files = f"{dir_path}/python-test/features/steps/pcap_files/"
    is_valid_dir = os.path.isdir(directory_of_network_data_files)
    assert_that(is_valid_dir, equal_to(True), f"Invalid directory {is_valid_dir}")
    network_data_files = files_names.split(", ")

    for network_file in network_data_files:
        path_to_file = f"{directory_of_network_data_files}{network_file}"
        valid_file_path = os.path.isfile(path_to_file)
        assert_that(valid_file_path, equal_to(True), f"Invalid path {valid_file_path}")
        run_mocked_data_command = f"tcpreplay -i {where_to_run} -tK {path_to_file}"
        tcpreplay_return = send_terminal_commands(run_mocked_data_command)
        check_successful_packets(tcpreplay_return[0])

        assert_that(tcpreplay_return[1], not_(contains_string("command not found")), f"{tcpreplay_return[1]}."
                                                                                     f"Please, install tcpreplay.")


def check_successful_packets(return_command_tcpreplay):
    return_command_tcpreplay = return_command_tcpreplay.replace("\t", "")
    return_command_tcpreplay = return_command_tcpreplay.replace("\n", "\",\"")
    return_command_tcpreplay = return_command_tcpreplay.replace(":", "\":\"")
    return_command_tcpreplay = return_command_tcpreplay.replace(" ", "")
    return_command_tcpreplay = return_command_tcpreplay[:-2]
    return_command_tcpreplay = "{\"" + return_command_tcpreplay + "}"
    return_command_tcpreplay = json.loads(return_command_tcpreplay)

    assert_that(int(return_command_tcpreplay['Actual'].split("packets")[0]),
                equal_to(int(return_command_tcpreplay['Successfulpackets'])), "Some packet may have failure")

    return return_command_tcpreplay


@then("remove all dummy ifaces and virtual switches generated on test process")
def cleanup_mocked_ifaces_and_switches(context):
    remove_all_virtual_switches()
    remove_all_dummys()


def remove_all_virtual_switches():
    command_1 = "ovs-vsctl list-br"
    all_switches = send_terminal_commands(command_1)
    all_switches = all_switches[0].split("\n")
    for switch in all_switches:
        if "brqa" in switch:
            command_2 = f"ovs-vsctl del-br {switch}"
            send_terminal_commands(command_2)
    all_switches = send_terminal_commands(command_1)
    all_switches = all_switches[0].split("\n")
    switch_test_remain = any('brqa' in switch for switch in all_switches)
    assert_that(switch_test_remain, equal_to(False), f"Unable to remove switches {switch_test_remain}")


def remove_all_dummys():
    command_1 = "ifconfig"
    dummys = send_terminal_commands(command_1)
    dummys = dummys[0].split("\n")
    dummys_test = [item for item in dummys if 'dummy' in item]
    for dummy in dummys_test:
        dummy_name = dummy.split(": ")[0]
        send_terminal_commands(f"ip link set {dummy_name} down")
        send_terminal_commands(f"ip link delete {dummy_name} type dummy")
    dummys = send_terminal_commands(command_1)
    dummys = dummys[0].split("\n")
    dummys_test = [item for item in dummys if 'dummyqa' in item]
    assert_that(dummys_test, equal_to([]), f"Unable to remove dummys {dummys_test}")
