import configparser
from hamcrest import *

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
    parser.read('test_config.ini')
    configs = parser['test_config']

    # Make sure mandatory parameters are set!
    assert_that(configs.get('orb_address'), not_none(), 'No Orb URL was provided!')
    assert_that(configs.get('orb_address'), not_(""), 'No Orb URL was provided!')

    assert_that(configs.get('email'), not_none(), 'No Orb user email was provided!')
    assert_that(configs.get('email'), not_(""), 'No Orb user email was provided!')

    assert_that(configs.get('password'), not_none(), 'No Orb user password was provided!')
    assert_that(configs.get('password'), not_(""), 'No Orb user password was provided!')
    assert_that(configs.get('password'), has_length(greater_than_or_equal_to(8)), 'Orb password must be at least 8 digits')

    assert_that(configs.get('prometheus_username'), not_none(), 'No Orb user password was provided!')
    assert_that(configs.get('prometheus_username'), not_(""), 'No Orb user password was provided!')

    assert_that(configs.get('prometheus_key'), not_none(), 'No Orb user password was provided!')
    assert_that(configs.get('prometheus_key'), not_(""), 'No Orb user password was provided!')

    assert_that(configs.get('remote_prometheus_endpoint'), not_none(), 'No Orb user password was provided!')
    assert_that(configs.get('remote_prometheus_endpoint'), not_(""), 'No Orb user password was provided!')

    ignore_ssl_and_certificate_errors = configs.get('ignore_ssl_and_certificate_errors', 'false').lower()
    assert_that(ignore_ssl_and_certificate_errors, any_of(equal_to('true'), equal_to('false')),
                'Invalid value to ignore_ssl_and_certificate_errors parameter. A boolean value is expected.')
    configs['ignore_ssl_and_certificate_errors'] = ignore_ssl_and_certificate_errors
    if ignore_ssl_and_certificate_errors.lower() == 'true':
        configs['orb_url'] = f"http://{configs.get('orb_address')}"
    else:
        configs['orb_url'] = f"https://{configs.get('orb_address')}"

    is_credentials_registered = configs.get('is_credentials_registered').lower()
    assert_that(is_credentials_registered, any_of(equal_to('true'), equal_to('false')),
                'Invalid value to is_credentials_registered parameter. A boolean value is expected.')
    configs['is_credentials_registered'] = is_credentials_registered
    return configs
