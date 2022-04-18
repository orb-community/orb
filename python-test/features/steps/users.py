from hamcrest import *
from behave import given, step, then
from test_config import TestConfig
import requests
from utils import random_string

configs = TestConfig.configs()
orb_url = configs.get('orb_url')


@step("that there is an unregistered {email} email with {password} password")
def check_non_registered_account(context, email, password):
    assert_that(email, any_of(equal_to("valid"), equal_to("invalid")), "Unexpected value for email")
    if email == "valid":
        context.email = [f"tester_{random_string(4)}@email.com", email]
        status_code = 403
    else:
        context.email = [f"tester.com", email]
        status_code = 400
    context.password = password
    authenticate(context.email[0], context.password, status_code)


@step("the Orb user request this account registration with {company} as company and {fullname} as fullname")
def request_account_registration(context, company, fullname):
    if company.lower() == "none": company = None
    if fullname.lower() == "none": fullname = None
    context.company = company
    context.full_name = fullname

    if context.email[1] == "invalid" or not (8 <= len(context.password) <= 50):
        status_code = 400
    else:
        status_code = 201
    context.registration_response, context.status_code = register_account(context.email[0], context.password,
                                                                          context.company, context.full_name,
                                                                          status_code)

@then("the status code must be {status_code}")
def register_orb_account(context, status_code):
    status_code = int(status_code)
    assert_that(context.status_code, equal_to(status_code), "Unexpected status code for account registration")



@step("account is registered with email, with password, {company_status} company and {fullname_status} full name")
def check_account_information(context, company_status, fullname_status):
    if company_status.lower() == "none": company_status = None
    if fullname_status.lower() == "none": fullname_status = None
    context.token = authenticate(context.email[0], context.password)["token"]
    account_details = get_account_information(context.token)
    expected_keys = ["id", "email", "metadata"]
    for key in expected_keys:
        assert_that(account_details, has_key(key), f"key {key} not present in account data")
    expected_metadata = dict()
    if company_status is not None:
        expected_metadata["company"] = context.company
    if fullname_status is not None:
        expected_metadata["fullName"] = context.full_name
    for metadata in expected_metadata.keys():
        assert_that(account_details["metadata"], has_key(metadata), f"Missing {metadata} in account metadata")
        assert_that(account_details["metadata"][metadata], equal_to(expected_metadata[metadata]),
                f"Unexpected value for {metadata}")


@given("the Orb user has a registered account")
def check_account_registration(context):
    email = configs.get('email')
    password = configs.get('password')
    if configs.get('is_credentials_registered') == 'false':
        register_account(email, password)
        configs['is_credentials_registered'] = 'true'
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
    json_request = {'email': user_email, 'password': user_password}
    json_request = {parameter: value for parameter, value in json_request.items() if value}
    response = requests.post(orb_url + '/api/v1/tokens',
                             json=json_request,
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
    response = requests.post(orb_url + '/api/v1/users',
                             json=json_request,
                             headers=headers)
    assert_that(response.status_code, equal_to(expected_status_code),
                f"Current value of is_credentials_registered parameter = {configs.get('is_credentials_registered')}."
                f"\nExpected status code for registering an account failed with status= {str(response.status_code)}.")
    return response, response.status_code,


def get_account_information(token, expected_status_code=200):
    """

    :param (str) token: used for API authentication
    :param (int) expected_status_code: expected request's status code. Default:200 (happy path).
    :return: (dict) account_response
    """
    response = requests.get(orb_url + '/api/v1/users/profile',
                            headers={'Authorization': token})
    assert_that(response.status_code, equal_to(expected_status_code), "Unexpected status code for get account data")
    return response.json()
