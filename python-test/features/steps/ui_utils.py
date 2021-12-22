from hamcrest import *
from selenium import webdriver


def go_to_page(page, context, ignore_ssl_and_certificate_errors=True):
    """Open the page in Chrome browser
    Args:
        :param (str) page: site's URL
        :param (behave.runner.Context) context: object that hold contextual information during the running of tests
        :param (bool) ignore_ssl_and_certificate_errors: if true, ignore ssl and certificate errors
    """
    options = webdriver.ChromeOptions()
    if ignore_ssl_and_certificate_errors is True:
        options.add_argument('--ignore-ssl-errors=yes')
        options.add_argument('--ignore-certificate-errors')
    options.add_argument("--start-maximized")
    context.driver = webdriver.Chrome(options=options)
    context.driver.get(str(page))
    assert_that(context.driver.current_url, equal_to(f"{page}/auth/login"), "user not enabled to access orb "
                                                                            "login page")


def input_text_by_id(element_id, information, context):
    """Send information required on a page
    Args:
        element_id (string): id of the element to be located
        information (string): information that should be sent
        context (behave.runner.Context): object that hold contextual information during the running of tests
    """
    data = context.driver.find_element_by_id(str(element_id))
    data.send_keys(str(information))


def input_text_by_xpath(element_xpath, information, context):
    """Send information required on a page
    Args:
        element_xpath (string): xpath of the element to be located
        information (string): information that should be sent
        context (behave.runner.Context): object that hold contextual information during the running of tests
    """
    data = context.driver.find_element_by_xpath(str(element_xpath))
    data.send_keys(str(information))
