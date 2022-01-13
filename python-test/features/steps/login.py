from behave import given, when, then, step
from test_config import TestConfig
from hamcrest import *
from users import authenticate, register_account
from utils import random_string, insert_str
from random import randint

configs = TestConfig.configs()
base_orb_url = configs.get('base_orb_url')


@given('there is a registered account')
def check_registered_account(context):
    is_credentials_registered = configs.get('is_credentials_registered')
    assert_that(is_credentials_registered, equal_to('true'), "This test require a registered account")
    email = configs.get('email')
    password = configs.get('password')
    authenticate(email, password)


@when('request referred account registration using registered email, {password_status} password, {username} user name and {company} company name')
def request_account_registration(context, password_status, username, company):
    assert_that(password_status, any_of(equal_to('registered'), equal_to('unregistered')),
                "Not expected option to password")
    if username.lower() == 'none': username = None
    if company.lower() == 'none': company = None
    email = configs.get('email')
    password = configs.get('password')
    if password_status == "unregistered":
        password1 = password + random_string(1)
        response = register_account(email, password1, company, username, 409)
        assert_that(response.json()['error'], equal_to('email already taken'), 'Wrong message on API response')
        password2 = password[:8]
        response = register_account(email, password2, company, username, 409)
        assert_that(response.json()['error'], equal_to('email already taken'), 'Wrong message on API response')
        if len(password) > 8:
            password3 = password[:-1]
            response = register_account(email, password3, company, username, 409)
            assert_that(response.json()['error'], equal_to('email already taken'), 'Wrong message on API response')
    else:
        response = register_account(email, password, company, username, 409)
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

