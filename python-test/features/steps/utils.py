import random
import string
from json import loads, JSONDecodeError

tag_key_prefix = "test_tag_key_"
tag_value_prefix = "test_tag_value_"


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


def generate_name_and_tag(name_prefix, tag_key_prefix, tag_value_prefix):
    """
    :param (str) name_prefix: prefix to identify object created by tests
    :param (str) tag_key_prefix: prefix to identify tag_key created by tests
    :param (str) tag_value_prefix: prefix to identify tag_value created by tests
    :return: random name, tag_key and tag_value
    """
    name = name_prefix + random_string(10)
    tag_key = tag_key_prefix + random_string(4)
    tag_value = tag_value_prefix + random_string(4)
    return name, tag_key, tag_value
