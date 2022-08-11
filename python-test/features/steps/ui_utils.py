from selenium import webdriver
from selenium.webdriver.chrome.service import Service as ChromeService
from webdriver_manager.chrome import ChromeDriverManager
from selenium.webdriver.common.by import By
from selenium.webdriver.support.ui import WebDriverWait
from selenium.webdriver.support import expected_conditions as EC
from selenium.webdriver.common.action_chains import ActionChains
from page_objects import UtilButton
from utils import threading_wait_until


def go_to_page(page, headless=True):
    """Open the page in Chrome browser
    Args:
        :param (str) page: site's URL
        :param (bool) headless: run webdriver in headless mode
    """
    options = webdriver.ChromeOptions()
    options.add_argument("--start-maximized")
    if headless:
        options.add_argument('headless')
    driver = webdriver.Chrome(service=ChromeService(ChromeDriverManager().install()), options=options)
    driver.get(str(page))
    return driver, str(driver.current_url)


def input_text_by_id(element_id, information, driver):
    """Send information required on a page
    Args:
        element_id (string): id of the element to be located
        information (string): information that should be sent
        driver (webdriver): selenium webdriver
    """
    action = ActionChains(driver)
    WebDriverWait(driver, 3).until(
        EC.element_to_be_clickable((By.ID, element_id)))
    data = driver.find_element(By.ID, (str(element_id)))
    action.send_keys_to_element(data, information).perform()
    WebDriverWait(driver, 3).until(lambda check: data.get_attribute('value') == information)


@threading_wait_until
def input_text_by_xpath(element_xpath, information, driver, event=None):
    """Send information required on a page
    Args:
        element_xpath (string): xpath of the element to be located
        information (string): information that should be sent
        driver (webdriver): selenium webdriver
    """
    WebDriverWait(driver, 3).until(
        EC.visibility_of_element_located((By.XPATH, element_xpath)))
    data = WebDriverWait(driver, 3).until(EC.presence_of_element_located((By.XPATH, element_xpath)))
    data.clear()
    data.send_keys(information)
    if data.get_attribute('value') == str(information):
        event.set()


def get_selector_options(driver, selector_options_xpath=UtilButton.selector_options()):
    """
    Args:
        driver (webdriver): selenium webdriver
        selector_options_xpath (str): xpath of the element to be located
    :return: dictionary in which the key is the name of the option and the value is the web element referring to the
    option in the selector

    """
    dict_options = dict()
    WebDriverWait(driver, 3).until(
        EC.presence_of_all_elements_located((By.XPATH, selector_options_xpath)))
    options = driver.find_elements(By.XPATH, selector_options_xpath)
    for option in options:
        dict_options[option.text] = option
    return dict_options
