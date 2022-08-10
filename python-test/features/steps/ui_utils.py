from selenium import webdriver
from selenium.webdriver.chrome.service import Service as ChromeService
from webdriver_manager.chrome import ChromeDriverManager
from selenium.webdriver.common.by import By


def go_to_page(page, context, headless=True):
    """Open the page in Chrome browser
    Args:
        :param (str) page: site's URL
        :param (behave.runner.Context) context: object that hold contextual information during the running of tests
        :param (bool) headless: run webdriver in headless mode
    """
    options = webdriver.ChromeOptions()
    options.add_argument("--start-maximized")
    if headless:
        options.add_argument('headless')
    context.driver = webdriver.Chrome(service=ChromeService(ChromeDriverManager().install()), options=options)
    context.driver.get(str(page))
    return str(context.driver.current_url)


def input_text_by_id(element_id, information, context):
    """Send information required on a page
    Args:
        element_id (string): id of the element to be located
        information (string): information that should be sent
        context (behave.runner.Context): object that hold contextual information during the running of tests
    """
    data = context.driver.find_element(By.ID, (str(element_id)))
    data.send_keys(str(information))


def input_text_by_xpath(element_xpath, information, context):
    """Send information required on a page
    Args:
        element_xpath (string): xpath of the element to be located
        information (string): information that should be sent
        context (behave.runner.Context): object that hold contextual information during the running of tests
    """
    data = context.driver.find_element(By.XPATH, (str(element_xpath)))
    data.send_keys(str(information))
