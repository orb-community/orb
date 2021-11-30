from hamcrest import *
from behave import given
from test_config import TestConfig
import requests

configs = TestConfig.configs()
base_orb_url = "https://" + configs.get('orb_address')


@given('that the user is logged in')
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
