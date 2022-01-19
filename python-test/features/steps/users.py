from hamcrest import *
from behave import given
from test_config import TestConfig
import requests

configs = TestConfig.configs()
base_orb_url = configs.get('base_orb_url')


@given("the Orb user has a registered account")
def register_orb_account(context):
    email = configs.get('email')
    password = configs.get('password')
    if configs.get('is_credentials_registered') == 'false':
        register_account(email, password)
    authenticate(email, password)


@given('the Orb user logs in')
def get_auth_token(context):
    email = configs.get('email')
    password = configs.get('password')
    response = authenticate(email, password)
    assert_that(response, has_key('token'), 'Authentication token not found in login response!')
    context.token = response['token']


def authenticate(user_email, user_password, expected_status_code=201):
    """
    Logs in to orb with given credentials

    :param (str) user_email: email of the user that is about to login
    :param (str) user_password: password of the user that is about to login
    :param (int) expected_status_code: expected request's status code. Default:201.
    :returns: (dict) response of auth request
    """

    headers = {'Content-type': 'application/json', 'Accept': '*/*'}
    response = requests.post(base_orb_url + '/api/v1/tokens',
                             json={'email': user_email, 'password': user_password},
                             headers=headers)
    assert_that(response.status_code, equal_to(expected_status_code),
                'Authentication failed with status= ' + str(response.status_code))

    return response.json()


def register_account(user_email, user_password, company_name=None, user_full_name=None, expected_status_code=201):
    """
    Attempt to register an account and asserts if the expected status code for an account registration with given
    credentials is correct

    :param (str) user_email: email of the user that is about to login
    :param (str) user_password: password of the user that is about to login
    :param (str) company_name: name of the company the user belongs to. Default: None
    :param (str) user_full_name: user full name. Default: None
    :param (int) expected_status_code: expected request's status code. Default:201 (happy path).
    """

    json_request = {"email": user_email, "password": user_password, "metadata": {"company": company_name,
                                                                                 "fullName": user_full_name}}
    json_request['metadata'] = {parameter: value for parameter, value in json_request['metadata'].items() if value}
    json_request = {parameter: value for parameter, value in json_request.items() if value}
    headers = {'Content-type': 'application/json', 'Accept': '*/*'}
    response = requests.post(base_orb_url + '/api/v1/users',
                             json=json_request,
                             headers=headers)
    assert_that(response.status_code, equal_to(expected_status_code),
                f"Current value of is_credentials_registered parameter = {configs.get('is_credentials_registered')}."
                f"\nExpected status code for registering an account failed with status= {str(response.status_code)}.")
    return response
