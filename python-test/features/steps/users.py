from hamcrest import *
from behave import given
from test_config import TestConfig
import requests

configs = TestConfig.configs()

if "localhost" in configs.get('orb_address'):
    base_orb_url = f"http://{configs.get('orb_address')}"
else:
    base_orb_url = f"https://{configs.get('orb_address')}"


@given("the Orb user has a registered account")
def register_orb_account(context):
    email = configs.get('email')
    password = configs.get('password')
    register_account(email, password)


@given('the Orb user logs in')
def get_auth_token(context):
    email = configs.get('email')
    password = configs.get('password')
    context.token = authenticate(email, password)


def authenticate(user_email, user_password):
    """
    Logs in to orb with given credentials

    :param (str) user_email: email of the user that is about to login
    :param (str) user_password: password of the user that is about to login
    :returns: (str) authentication token used to perform API calls to orb
    """

    headers = {'Content-type': 'application/json', 'Accept': '*/*'}
    response = requests.post(base_orb_url + '/api/v1/tokens',
                             json={'email': user_email, 'password': user_password},
                             headers=headers)
    assert_that(response.status_code, equal_to(201),
                'Get token request failed with status=' + str(response.status_code))
    assert_that(response.json(), has_key('token'), 'Authentication token not found in login response!')

    return response.json()['token']


def register_account(user_email, user_password, expected_status_code=201):
    """
    Logs in to orb with given credentials

    :param (str) user_email: email of the user that is about to login
    :param (str) user_password: password of the user that is about to login
    :param (int) expected_status_code: expected request's status
    """

    headers = {'Content-type': 'application/json', 'Accept': '*/*'}
    response = requests.post(base_orb_url + '/api/v1/users',
                             json={'email': user_email, 'password': user_password},
                             headers=headers)
    assert_that(response.status_code, equal_to(expected_status_code),
                'Register an account request failed with status=' + str(response.status_code))
