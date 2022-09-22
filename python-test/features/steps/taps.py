from hamcrest import *


class Taps:
    def __init__(self):
        pass

    @classmethod
    def pcap(cls, name, input_type="pcap", settings=None):
        available_options = ["pcap_source", "iface", "host_spec", "debug", "bpf"]
        filters_list = ["bpf"]

        return make_tap(name, input_type, available_options, settings, filters_list)

    @classmethod
    def flow(cls, name, input_type="flow", settings=None):
        available_options = ["port", "bind", "flow_type"]

        return make_tap(name, input_type, available_options, settings)

    @classmethod
    def dnstap(cls, name, input_type="dnstap", settings=None):
        available_options = ["socket", "tcp", "only_hosts"]

        filters_list = ["only_hosts"]

        return make_tap(name, input_type, available_options, settings, filters_list)


def make_tap(name, input_type, available_options, settings, filters_list=None):
    if filters_list is None:
        filters_list = []
    kwargs_configs = list(settings.keys())

    assert_that(set(kwargs_configs).issubset(available_options), is_(True),
                f"Invalid configuration to tap {input_type}. \n "
                f"Options are: {available_options}. \n"
                f"Passed: {kwargs_configs}")

    filters = None
    for tap_filter in filters_list:
        if tap_filter in kwargs_configs:
            filters = {tap_filter: settings[tap_filter]}
            kwargs_configs.remove(tap_filter)

    if len(kwargs_configs) > 0:
        configs = dict()
    else:
        configs = None
    for configuration in kwargs_configs:
        configs.update({configuration: settings[configuration]})

    tap = {name: {"input_type": input_type}}

    if filters is not None:
        filters = {"filter": filters}
        tap[name].update(filters)

    if configs is not None:
        configs = {"config": configs}
        tap[name].update(configs)
    return tap
