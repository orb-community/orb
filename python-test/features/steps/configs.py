import configparser
from hamcrest import *
import os

LOCAL_AGENT_CONTAINER_NAME = "orb-agent-int-test"


class TestConfig:
    _configs = None

    def __init__(self):
        raise RuntimeError('Call instance() instead')

    @classmethod
    def configs(cls):
        if cls._configs is None:
            cls._configs = _read_configs()
        return cls._configs


def _read_configs():
    parser = configparser.ConfigParser()
    parser.read('/home/amanda-silva/Documents/orb/python-test/test_config.ini')
    configs = parser['test_config']
    ################
    # todo remover

    # log level
    log_level = configs.get("log_level", "debug")
    configs["log_level"] = log_level.lower()
    assert_that(configs.get('log_level'), any_of("debug", "info", "warning", "error", "critical"),
                'Invalid log level')

    # stream logs
    stream_logs = configs.get("stream_logs", "true")
    assert_that(stream_logs.lower(), any_of("true", "false"), 'Invalid log stream option')
    configs["stream_logs"] = stream_logs.lower()

    # Make sure mandatory parameters are set!
    assert_that(configs.get('orb_address'), not_none(), 'No Orb URL was provided!')
    assert_that(configs.get('orb_address'), not_(""), 'No Orb URL was provided!')

    assert_that(configs.get('email'), not_none(), 'No Orb user email was provided!')
    assert_that(configs.get('email'), not_(""), 'No Orb user email was provided!')

    assert_that(configs.get('password'), not_none(), 'No Orb user password was provided!')
    assert_that(configs.get('password'), not_(""), 'No Orb user password was provided!')
    assert_that(configs.get('password'), has_length(greater_than_or_equal_to(8)),
                'Orb password must be at least 8 digits')

    configs['backend_type'] = configs.get("backend_type", "prometheus")
    if configs['backend_type'] == "otlphttp":
        assert_that(configs.get('otlp_publisher_username'), not_none(), 'No otlp_publisher username was provided!')
        assert_that(configs.get('otlp_publisher_username'), not_(""), 'No otlp_publisher username was provided!')

        assert_that(configs.get('otlp_publisher_key'), not_none(), 'No Orb otlp_publisher key was provided!')
        assert_that(configs.get('otlp_publisher_key'), not_(""), 'No Orb otlp_publisher key was provided!')

        assert_that(configs.get('otlp_publisher_endpoint'), not_none(), 'No Orb otlp_publisher endpoint was provided!')
        assert_that(configs.get('otlp_publisher_endpoint'), not_(""), 'No Orb otlp_publisher endpoint was provided!')
    else:
        assert_that(configs.get('prometheus_username'), not_none(), 'No prometheus username was provided!')
        assert_that(configs.get('prometheus_username'), not_(""), 'No prometheus username was provided!')

        assert_that(configs.get('prometheus_key'), not_none(), 'No Orb prometheus key was provided!')
        assert_that(configs.get('prometheus_key'), not_(""), 'No Orb prometheus key was provided!')

        assert_that(configs.get('remote_prometheus_endpoint'), not_none(), 'No Orb prometheus endpoint was provided!')
        assert_that(configs.get('remote_prometheus_endpoint'), not_(""), 'No Orb prometheus endpoint was provided!')

    local_orb_path = configs.get("orb_path",
                                 os.path.dirname(os.getcwd()))  # orb_path is required if user will use docker to test,
    # otherwise the function will map the local path.
    assert_that(os.path.exists(local_orb_path), equal_to(True), f"Invalid orb_path: {local_orb_path}.")
    configs['local_orb_path'] = local_orb_path

    use_orb_live_address_pattern = configs.get('use_orb_live_address_pattern', 'true').lower()
    assert_that(use_orb_live_address_pattern, any_of(equal_to('true'), equal_to('false')),
                'Invalid value to use_orb_live_address_pattern parameter. A boolean value is expected.')

    verify_ssl = configs.get('verify_ssl', 'true').lower()
    assert_that(verify_ssl, any_of(equal_to('true'), equal_to('false')),
                'Invalid value to verify_ssl parameter. A boolean value is expected.')
    configs['verify_ssl'] = verify_ssl
    # use agents. on the beginning of mqtt address if true
    if use_orb_live_address_pattern == "true":
        configs['mqtt_base_address'] = f"agents.{configs.get('orb_address')}"
        if verify_ssl == 'true':
            configs['orb_url'] = f"https://{configs.get('orb_address')}"
            configs['mqtt_url'] = f"tls://{configs['mqtt_base_address']}:8883"
        else:
            configs['orb_url'] = f"http://{configs.get('orb_address')}"
            configs['mqtt_url'] = f"{configs['mqtt_base_address']}:8883"
    else:
        configs['orb_url'] = configs.get('orb_cloud_api_address', 'None')
        assert_that(configs['orb_url'], is_not('None'), "If use_orb_live_address_pattern is not true, you need to "
                                                        "insert your orb_cloud_api_address")
        configs['mqtt_url'] = configs.get('orb_cloud_mqtt_address', 'None')
        assert_that(configs['mqtt_url'], is_not('None'), "If use_orb_live_address_pattern is not true, you need to "
                                                         "insert your orb_cloud_mqtt_address")

    is_credentials_registered = configs.get('is_credentials_registered', 'true').lower()
    assert_that(is_credentials_registered, any_of(equal_to('true'), equal_to('false')),
                'Invalid value to is_credentials_registered parameter. A boolean value is expected.')
    configs['is_credentials_registered'] = is_credentials_registered

    return configs
