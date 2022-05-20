from users import authenticate
from behave import given, when, then, step
from test_config import TestConfig
from selenium.webdriver.support.ui import WebDriverWait
from selenium.webdriver.support import expected_conditions as EC
from ui_utils import go_to_page, input_text_by_id
from hamcrest import *


configs = TestConfig.configs()
user_email = configs.get('email')
user_password = configs.get('password')
orb_url = configs.get('orb_url')


@given("the Orb user logs in through the UI")
def logs_in_orb_ui(context):
    orb_page(context)
    use_credentials(context)
    check_home_page(context)
    context.token = authenticate(user_email, user_password)['token']


@step("that the user is on the orb page")
def orb_page(context):
    current_url = go_to_page(orb_url, context)
    assert_that(current_url, equal_to(f"{orb_url}/auth/login"), "user not enabled to access orb login page")


@when("the Orb user logs in Orb UI")
def use_credentials(context):
    input_text_by_id("input-email", user_email, context)
    input_text_by_id("input-password", user_password, context)
    context.driver.find_element_by_css_selector(str(".appearance-filled")).click()


@then("the user should have access to orb home page")
def check_home_page(context):
    WebDriverWait(context.driver, 10).until(EC.url_to_be(f"{orb_url}/pages/home"), message="user not enabled to "
                                                                                                "access orb home page")


