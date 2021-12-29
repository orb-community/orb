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

    ignore_ssl_and_certificate_errors = configs.get('ignore_ssl_and_certificate_errors', 'false')
    configs['ignore_ssl_and_certificate_errors'] = ignore_ssl_and_certificate_errors
    if ignore_ssl_and_certificate_errors.lower() == 'true':
        configs['base_orb_url'] = f"http://{configs.get('orb_address')}"
    else:
        configs['base_orb_url'] = f"https://{configs.get('orb_address')}"

    return configs
