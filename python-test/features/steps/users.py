from hamcrest import *
from behave import given
from test_config import TestConfig, base_orb_url
import requests

configs = TestConfig.configs()


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
    Asserts if the expected status code for an account registration with given credentials is correct

    :param (str) user_email: email of the user that is about to login
    :param (str) user_password: password of the user that is about to login
    :param (int) expected_status_code: expected request's status code. Default:201 (happy path).
    """

    headers = {'Content-type': 'application/json', 'Accept': '*/*'}
    response = requests.post(base_orb_url + '/api/v1/users',
                             json={'email': user_email, 'password': user_password},
                             headers=headers)
    assert_that(response.status_code, equal_to(expected_status_code),
                'Expected status code for registering an account failed with status=' + str(response.status_code))
