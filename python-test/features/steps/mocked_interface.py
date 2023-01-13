from hamcrest import *
from behave import given, then, step
from test_config import TestConfig
from utils import random_string, send_terminal_commands
import os
import json

configs = TestConfig.configs()
orb_url = configs.get('orb_url')
verify_ssl_bool = eval(configs.get('verify_ssl').title())


@given("a mocked interface is configured with mtu: {mtu} and ip: {ip}")
def set_mock_interface(context, mtu, ip):
    context.mock_iface_name = f"dummy{random_string(10)}"
    return_0 = send_terminal_commands(f"ip link add {context.mock_iface_name} type dummy")
    return_1 = send_terminal_commands(f"ip link set {context.mock_iface_name} up")
    return_2 = send_terminal_commands(f"ifconfig {context.mock_iface_name} mtu {mtu}")
    if "Operation not permitted" in return_0[1] or return_1 != ('', ''):
        context.access_denied = True
        context.scenario.skip('Root privileges are required')
    assert_that(return_2, equal_to(('', '')), f"Error while configuring mtu of dummy iface")
    if ip.lower() != "none":
        return_4 = send_terminal_commands(f"ip addr add {ip} dev {context.mock_iface_name}")
        assert_that(return_4, equal_to(('', '')), f"Error while configuring ip of dummy iface")


@then("remove dummy interface")
def remove_dummy(context):
    send_terminal_commands(f"ip link set {context.mock_iface_name} down")
    send_terminal_commands(f"ip link delete {context.mock_iface_name} type dummy")


@step("run mocked data {files_names} for this network")
def run_mocked_data(context, files_names):
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
        run_mocked_data_command = f"tcpreplay -i {context.mock_iface_name} -tK {path_to_file}"
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
