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

    return configs


def bypass_ssl_certificate_check():
    ignore_ssl_and_certificate_errors = TestConfig.configs().get('ignore_ssl_and_certificate_errors', False)
    if ignore_ssl_and_certificate_errors:
        orb_url = f"http://{TestConfig.configs().get('orb_address')}"
    else:
        orb_url = f"https://{TestConfig.configs().get('orb_address')}"
    return ignore_ssl_and_certificate_errors, orb_url


bypass_ssl_certificate_check, base_orb_url = bypass_ssl_certificate_check()
