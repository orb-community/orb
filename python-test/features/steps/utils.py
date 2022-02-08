import random
import string
from json import loads, JSONDecodeError
import socket, errno
from hamcrest import *
from behave import step

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
    :param (str) n_random: amount of random characters
    :return: random_string_with_predefined_prefix
    """
    random_string_with_predefined_prefix = string_prefix + random_string(n_random)
    return random_string_with_predefined_prefix


def create_tags_set(tags_type, orb_tags):
    """
    Create a set of orb-tags
    :param (str) tags_type: Defines if the set of tags is previously defined or if it should be randomly generated.
                    Options: "defined" or "random".
    :param orb_tags: If defined: the defined tags that should compose the set.
                     If random: the number of tags that the set must contain.
    :return: (dict) tag_set
    """
    assert_that(tags_type, any_of(equal_to('random'), equal_to('defined')), 'Unexpect type of tags')
    tag_set = dict()
    if tags_type == 'defined':
        for tag in orb_tags.split(", "):
            key, value = tag.split(":")
            tag_set[key] = value
    else:
        assert_that(orb_tags, matches_regexp("^\d+\stag\(s\)"), "Unexpected value for tags: n tag(s) is expected")
        amount_of_tags = int(orb_tags.split()[0])
        for tag in range(amount_of_tags):
            tag_set[tag_prefix + random_string(4)] = tag_prefix + random_string(2)
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
