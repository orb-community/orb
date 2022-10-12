import threading
from hamcrest import *
from selenium import webdriver
from selenium.webdriver.chrome.service import Service as ChromeService
from webdriver_manager.chrome import ChromeDriverManager
from selenium.webdriver.common.by import By
from selenium.webdriver.support.ui import WebDriverWait
from selenium.webdriver.support import expected_conditions as EC
from selenium.webdriver.common.action_chains import ActionChains
from page_objects import *
from selenium.common.exceptions import TimeoutException, StaleElementReferenceException, \
    ElementClickInterceptedException
from utils import threading_wait_until

time_webdriver_wait = 5


def go_to_page(page, headless=True, driver=None):
    """Open the page in Chrome browser
    Args:
        :param (str) page: site's URL
        :param (bool) headless: run webdriver in headless mode
        :param (obj) driver: selenium webdriver
    """
    if driver is None:
        options = webdriver.ChromeOptions()
        options.add_argument("--start-maximized")
        options.add_argument("--remote-debugging-port=9222")
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
    try:
        action = ActionChains(driver)
        WebDriverWait(driver, time_webdriver_wait).until(
            EC.element_to_be_clickable((By.ID, element_id)))
        data = driver.find_element(By.ID, (str(element_id)))
        action.send_keys_to_element(data, information).perform()
        WebDriverWait(driver, time_webdriver_wait).until(lambda check: data.get_attribute('value') == information)
    except Exception as error:
        raise ValueError(f"Fails to input text by ID using ID {element_id}. Exception: {error}")


@threading_wait_until
def input_text_by_xpath(element_xpath, information, driver, event=None):
    """Send information required on a page
    Args:
        element_xpath (string): xpath of the element to be located
        information (string): information that should be sent
        driver (webdriver): selenium webdriver
    """
    try:
        WebDriverWait(driver, time_webdriver_wait).until(
            EC.visibility_of_element_located((By.XPATH, element_xpath)))
        data = WebDriverWait(driver, time_webdriver_wait).until(
            EC.presence_of_element_located((By.XPATH, element_xpath)))
        data.click()
        data.send_keys(information)
        if data.get_attribute('value') == str(information):
            event.set()
    except ElementClickInterceptedException:
        print(ElementClickInterceptedException)
        event.wait(1)
    except Exception as error:
        raise ValueError(f"Fails to input text by XPATH using XPATH {element_xpath}. Exception: {error}")


def get_selector_options(driver, selector_options_xpath=UtilButton.selector_options()):
    """
    Args:
        driver (webdriver): selenium webdriver
        selector_options_xpath (str): xpath of the element to be located
    :return: dictionary in which the key is the name of the option and the value is the web element referring to the
    option in the selector

    """
    dict_options = dict()
    WebDriverWait(driver, time_webdriver_wait).until(
        EC.presence_of_all_elements_located((By.XPATH, selector_options_xpath)))
    options = driver.find_elements(By.XPATH, selector_options_xpath)
    for option in options:
        dict_options[option.text] = option
    return dict_options


def find_element_on_datatable(driver, xpath):
    """
    Find element present on agent datatable.

    :param (selenium obj) driver: webdriver running
    :param (str) xpath: xpath of the element to be found
    :return: web element, if found. None if not found.
    """
    WebDriverWait(driver, time_webdriver_wait).until(
        EC.presence_of_all_elements_located((By.XPATH, DataTable.page_count())), message="Unable to find page count")
    WebDriverWait(driver, time_webdriver_wait).until(
        EC.presence_of_all_elements_located((By.XPATH, DataTable.body())), message="Unable to find list body")
    pages = WebDriverWait(driver, time_webdriver_wait).until(
        EC.presence_of_all_elements_located((By.XPATH, DataTable.sub_pages())),
        message="Unable to find subpages")
    if len(pages) > 1:
        WebDriverWait(driver, time_webdriver_wait).until(
            EC.presence_of_element_located((By.XPATH, DataTable.last_page())), message="Unable to find 'go to last "
                                                                                       "page' button")
        button_was_clicked = button_click_by_xpath(DataTable.last_page(), driver,
                                                   "Unable to click on 'go to the last page' button")
        assert_that(button_was_clicked, equal_to(True), "Unable to click on 'go to the last page' button")

        last_pages = WebDriverWait(driver, time_webdriver_wait).until(EC.presence_of_all_elements_located((By.XPATH,
                                                                                                           DataTable.sub_pages())),
                                                                      message="Unable to find subpages")
        last_page = int(last_pages[-1].text)
        for page in range(last_page):
            try:
                element = WebDriverWait(driver, time_webdriver_wait).until(
                    EC.presence_of_element_located((By.XPATH, xpath)))
                return element
            except TimeoutException:
                button_was_clicked = button_click_by_xpath(DataTable.previous_page(), driver,
                                                           "Unable to click on 'go to the previous page' button")
                assert_that(button_was_clicked, equal_to(True), "Unable to click on 'go to the previous page' button")
            except StaleElementReferenceException:
                driver.refresh()
                threading.Event().wait(1)
            except OSError as err:
                print(err)
        return None
    else:
        try:
            element = WebDriverWait(driver, time_webdriver_wait).until(
                EC.presence_of_element_located((By.XPATH, xpath)))
            return element
        except TimeoutException:
            return None
        except StaleElementReferenceException:
            driver.refresh()
            threading.Event().wait(time_webdriver_wait)
            return StaleElementReferenceException
        except OSError as err:
            return err


@threading_wait_until
def find_element_on_datatable_by_condition(driver, element_xpath, xpath_referring_page, condition="is", event=None):
    try:
        assert_that(condition, any_of(equal_to("is"), equal_to("is not")), "Unexpected value for list condition")
        message = f"Unable to find {xpath_referring_page} icon on left menu"
        button_was_clicked = button_click_by_xpath(xpath_referring_page, driver, message)
        assert_that(button_was_clicked, equal_to(True), message)
        element_on_datatable = find_element_on_datatable(driver, element_xpath)
        if condition == "is" and element_on_datatable is not None:
            event.set()
        elif condition == "is not" and element_on_datatable is None:
            event.set()
        return element_on_datatable
    except TimeoutException:
        driver.refresh()
        event.wait(1)
        print(TimeoutException)
    except StaleElementReferenceException:
        driver.refresh()
        event.wait(1)
        print(StaleElementReferenceException)
    except Exception as exception:
        print(exception)
        event.wait(1)


@threading_wait_until
def button_click_by_xpath(button_path, driver, message, time_to_wait=5, event=None):
    try:
        WebDriverWait(driver, time_to_wait).until(
            EC.element_to_be_clickable((By.XPATH, button_path)), message=message).click()
        event.set()
        return event.is_set()
    except ElementClickInterceptedException:
        try:
            element_to_click = WebDriverWait(driver, time_to_wait).until(
                EC.element_to_be_clickable((By.XPATH, button_path)), message=message)
            driver.execute_script("arguments[0].click();", element_to_click)
            event.set()
            return event.is_set()
        except Exception as exception:
            print(exception)
            event.wait(time_to_wait)
            return event.is_set()
    except Exception as exception:
        print(exception)
        event.wait(time_to_wait)
        return event.is_set()


class OrbUrl:
    def __init__(self):
        pass

    @classmethod
    def agent_view(cls, orb_url, agent_id):
        return f"{orb_url}/pages/fleet/agents/view/{agent_id}"

    @classmethod
    def policy_view(cls, orb_url, policy_id):
        return f"{orb_url}/pages/datasets/policies/view/{policy_id}"
