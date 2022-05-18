import random
import string
from json import loads, JSONDecodeError
from hamcrest import *
import threading
from datetime import datetime
import socket
import os
import re

tag_prefix = "test_tag_"


def random_string(k=10):
    """
    Generates a string composed of of k (int) random letters lowercase and uppercase mixed

    :param (int) k: sets the length of the randomly generated string
    :return: (str) string consisting of k random letters lowercase and uppercase mixed. Default:10
    """
    return ''.join(random.choices(string.ascii_letters, k=k))


def safe_load_json(json_str):
    """
    Safely parses a string into a JSON object, without ever raising an error.
    :param (str) json_str: to be loaded
    :return: the JSON object, or None if string is not a valid JSON.
    """

    try:
        return loads(json_str)
    except JSONDecodeError:
        return None


def filter_list_by_parameter_start_with(list_of_elements, parameter, start_with):
    """
    :param (list) list_of_elements: a list of elements
    :param (str) parameter: key of dict whose values will be used to filter the elements
    :param (str) start_with: prefix that will be used to filter the elements that start with it
    :return: (list) a list of filtered elements
    """
    list_of_filtered_elements = list()
    for element in list_of_elements:
        if element[parameter].startswith(start_with):
            list_of_filtered_elements.append(element)
    return list_of_filtered_elements


def insert_str(str_base, str_to_insert, index):
    """

    :param (str) str_base: string in which some letter will be inserted
    :param (str) str_to_insert: letter to be inserted
    :param (int) index: position that letter should be inserted
    :return: (str) string with letter inserted on determined index
    """
    return str_base[:index] + str_to_insert + str_base[index:]


def generate_random_string_with_predefined_prefix(string_prefix, n_random=10):
    """
    :param (str) string_prefix: prefix to identify object created by tests
    :param (int) n_random: amount of random characters
    :return: random_string_with_predefined_prefix
    """
    random_string_with_predefined_prefix = string_prefix + random_string(n_random)
    return random_string_with_predefined_prefix


def create_tags_set(orb_tags):
    """
    Create a set of orb-tags
    :param orb_tags: If defined: the defined tags that should compose the set.
                     If random: the number of tags that the set must contain.
    :return: (dict) tag_set
    """
    tag_set = dict()
    if orb_tags.isdigit() is False:
        assert_that(orb_tags, any_of(matches_regexp("^.+\:.+"), matches_regexp("\d+ orb tag\(s\)"),
                                     matches_regexp("\d+ orb tag")), f"Unexpected regex for tags. Passed: {orb_tags}."
                                                                     f"Expected (examples):"
                                                                     f"If you want 1 randomized tag: 1 orb tag."
                                                                     f"If you want more than 1 randomized tags: 2 orb tags. Note that you can use any int. 2 its only an example."
                                                                     f"If you want specified tags: test_key:test_value, second_key:second_value.")
        if re.match(r"^.+\:.+", orb_tags): # We expected key values separated by a colon ":" and multiple tags separated
            # by a comma ",". Example: test_key:test_value, my_orb_key:my_orb_value
            for tag in orb_tags.split(", "):
                key, value = tag.split(":")
                tag_set[key] = value
                return tag_set
    amount_of_tags = int(orb_tags.split()[0])
    for tag in range(amount_of_tags):
        tag_set[tag_prefix + random_string(6)] = tag_prefix + random_string(4)
    return tag_set


def check_logs_contain_message_and_name(logs, expected_message, name, name_key):
    """
    Gets the logs from Orb agent container

    :param (list) logs: list of log lines
    :param (str) expected_message: message that we expect to find in the logs
    :param (str) name: element name that we expect to find in the logs
    :param (str) name_key: key to get element name on log line
    :returns: (bool) whether expected message was found in the logs
    """

    for log_line in logs:
        log_line = safe_load_json(log_line)

        if log_line is not None and log_line['msg'] == expected_message:
            if log_line is not None and log_line[name_key] == name:
                return True, log_line

    return False, "Logs doesn't contain the message and name expected"


def remove_empty_from_json(json_file):
    """
    Delete keys with the value "None" in a dictionary, recursively.

    """
    for key, value in list(json_file.items()):
        if value is None:
            del json_file[key]
        elif isinstance(value, dict):
            remove_empty_from_json(value)
    return json_file


def remove_key_from_json(json_file, key_to_be_removed):
    """

    :param json_file: json object
    :param key_to_be_removed: key that need to be removed
    :return: json object without keys removed
    """
    for key, value in list(json_file.items()):
        if key == key_to_be_removed:
            del json_file[key]
        elif isinstance(value, dict):
            remove_key_from_json(value, key_to_be_removed)
    return json_file


def threading_wait_until(func):
    def wait_event(*args, wait_time=0.5, timeout=10, start_func_value=False, **kwargs):
        event = threading.Event()
        func_value = start_func_value
        start = datetime.now().timestamp()
        time_running = 0
        while not event.is_set() and time_running < int(timeout):
            func_value = func(*args, event=event, **kwargs)
            event.wait(wait_time)
            time_running = datetime.now().timestamp() - start
        return func_value

    return wait_event


def check_port_is_available(availability=True):
    """

    :param (str) availability: Status of the port on which agent must try to run. Default: available.
    :return: (int) port number
    """
    assert_that(availability, any_of(equal_to(True), equal_to(False)), "Unexpected value for availability")
    available_port = None
    port_options = range(10853, 10900)
    for port in port_options:
        sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        result = sock.connect_ex(('127.0.0.1', port))
        sock.close()
        if result == 0:
            available_port = port
            if availability is True:
                continue
            else:
                return available_port
        else:
            available_port = port
            break
    assert_that(available_port, is_not(equal_to(None)), "No available ports to bind")
    return available_port


def find_files(prefix, suffix, path):
    """
    Find files that match with prefix and suffix condition

    :param prefix: string with which the file should start with
    :param suffix: string with which the file should end with
    :param path: directory where the files will be searched
    :return: (list) path to all files that matches
    """
    result = list()
    for root, dirs, files in os.walk(path):
        for name in files:
            if name.startswith(prefix) and name.endswith(suffix):
                result.append(os.path.join(root, name))
    return result

