from selenium import webdriver
from behave import given, when, then
from test_config import TestConfig
from selenium.webdriver.support.ui import WebDriverWait
from selenium.webdriver.support import expected_conditions as EC
from hamcrest import *

configs = TestConfig.configs()
base_orb_url = "https://" + configs.get('orb_address')
user_email = configs.get('email')
user_password = configs.get('password')


@given("that the user is on the orb page")
def orb_page(context):
    """Open the required page on Chrome browser
    Args:
        context (behave.runner.Context): object that hold contextual information during the running of tests
    """
    go_to_page(base_orb_url, context)


@when("the Orb user logs in Orb UI")
def use_credentials(context):
    """
    Set the credentials for login in Orb
    :param context: Behave class that contains contextual information during the running of tests.
    """
    enter_information("input-email", user_email, context)
    enter_information("input-password", user_password, context)
    context.driver.find_element_by_css_selector(str(".appearance-filled")).click()


@then("the user should have access to orb home page")
def check_home_page(context):
    WebDriverWait(context.driver, 10).until(EC.url_to_be(f"{base_orb_url}/pages/home"), message="user not enabled to "
                                                                                                "access orb home page")


def go_to_page(page, context):
    """Open the page in Chrome browser
    Args:
        page (string): site's URL
        context (behave.runner.Context): object that hold contextual information during the running of tests
    """
    options = webdriver.ChromeOptions()
    options.add_argument('--ignore-ssl-errors=yes')
    options.add_argument('--ignore-certificate-errors')
    options.add_argument("--start-maximized")
    context.driver = webdriver.Chrome(options=options)
    context.driver.get(str(page))
    assert_that(context.driver.current_url, equal_to(f"{base_orb_url}/auth/login"), "user not enabled to access orb "
                                                                                    "login page")


def enter_information(html_information_name, information, context, source_by=0):
    """Send information required on a page
    Args:
        html_information_name (string): name of html object
        information (string): information that should be sent (for ex: container's id)
        context (behave.runner.Context): object that hold contextual information during the running of tests
        source_by (int, optional): type of source by. If 0: by id; if 1:css selector; if 2: xpath. Defaults to 0.
    """
    if source_by == 0:
        data = context.driver.find_element_by_id(str(html_information_name))
        data.send_keys(str(information))
    elif source_by == 1:
        data = context.driver.find_element_by_css_selector(str(html_information_name))
        data.send_keys(str(information))
    elif source_by == 2:
        data = context.driver.find_element_by_xpath(str(html_information_name))
        data.send_keys(str(information))
    else:
        raise Exception("source_by not found")
