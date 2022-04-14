from behave import given, when, then, step
from test_config import TestConfig
from hamcrest import *
from users import authenticate, register_account
from utils import random_string, insert_str
from random import randint

configs = TestConfig.configs()
orb_url = configs.get('orb_url')


@given('there is a registered account')
def check_registered_account(context):
    is_credentials_registered = configs.get('is_credentials_registered')
    assert_that(is_credentials_registered, equal_to('true'), "This test require a registered account")
    email = configs.get('email')
    password = configs.get('password')
    authenticate(email, password)


@when("user request account registration {email} email, {password} password, {username} user name and {company} "
      "company name")
def check_account_input(context, email, password, username, company):
    inputs = {'email': email, 'password': password, 'username': username, 'company': company}
    for key, value in inputs.items():
        assert_that(value, any_of(equal_to('with'), equal_to('without')),
                    f"Not expected option to {key}")
    account_input = {'email': None, 'password': None, 'company': None, 'username': None, 'reg_status': 201,
                     'auth_status': 201}
    if email == "without" or password == "without":
        account_input['reg_status'] = 400
        account_input['auth_status'] = 400
    elif email == "with":
        account_input['email'] = f"test_email_{random_string(3)}@email.com"
        if password == "without":
            account_input['auth_status'] = 403
    elif password == "with":
        account_input['password'] = configs.get('password')
    if username == "with":
        account_input['username'] = f"test_user {random_string(3)}"
    if company == "with":
        account_input['company'] = f"test_company {random_string(3)}"

    register_account(account_input['email'], account_input['password'], account_input['company'],
                     account_input['username'], account_input['reg_status'])
    context.auth_response = authenticate(account_input['email'], account_input['password'], account_input['auth_status'])


@when('request referred account registration using registered email, {password_status} password, {username} user name '
      'and {company} company name')
def request_account_registration(context, password_status, username, company):
    assert_that(password_status, any_of(equal_to('registered'), equal_to('unregistered')),
                "Not expected option to password")
    if username.lower() == 'none': username = None
    if company.lower() == 'none': company = None
    email = configs.get('email')
    password = configs.get('password')
    if password_status == "unregistered":
        passwords_to_test = list()
        passwords_to_test.append(password + random_string(1))
        passwords_to_test.append(password[:8])
        if len(password) > 8: passwords_to_test.append(password[:-1])
        for current_password in passwords_to_test:
            response, status_code = register_account(email, current_password, company, username, 409)
            assert_that(response.json()['error'], equal_to('email already taken'), 'Wrong message on API response')
    else:
        response, status_code = register_account(email, password, company, username, 409)
        assert_that(response.json()['error'], equal_to('email already taken'), 'Wrong message on API response')


@then('account register should not be changed')
def check_users_account(context):
    email = configs.get('email')
    password = configs.get('password')
    auth_response = authenticate(email, password)
    assert_that(auth_response, has_key('token'), 'Authentication token not found in login response!')


@when('the Orb user request an authentication token using {email_status} email and {password_status} password')
def request_orb_authentication(context, email_status, password_status):
    assert_that(email_status, any_of(equal_to('correct'), equal_to('incorrect')), 'unexpected email status')
    assert_that(password_status, any_of(equal_to('correct'), equal_to('incorrect')), 'unexpected password status')
    email = configs.get('email')
    password = configs.get('password')
    if email_status == 'incorrect':
        email = insert_str(email, random_string(1), randint(0, len(email)))
    if password_status == 'incorrect':
        password = insert_str(password, random_string(1), randint(0, len(password)))

    context.auth_response = authenticate(email, password, 403)
    assert_that(context.auth_response['error'], equal_to("missing or invalid credentials provided"))


@then('user should not be able to authenticate')
def check_access_denied(context):
    assert_that(context.auth_response, not_(has_key("token")))
    assert_that(context.auth_response.keys(), only_contains("error"))
