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
from selenium.common.exceptions import TimeoutException, StaleElementReferenceException
from utils import threading_wait_until


def go_to_page(page, headless=True):
    """Open the page in Chrome browser
    Args:
        :param (str) page: site's URL
        :param (bool) headless: run webdriver in headless mode
    """
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
    data.click()
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


def find_element_on_datatable(driver, xpath):
    """
    Find element present on agent datatable.

    :param (selenium obj) driver: webdriver running
    :param (str) xpath: xpath of the element to be found
    :return: web element, if found. None if not found.
    """
    WebDriverWait(driver, 3).until(
        EC.presence_of_all_elements_located((By.XPATH, DataTable.page_count())), message="Unable to find page count")
    WebDriverWait(driver, 3).until(
        EC.presence_of_all_elements_located((By.XPATH, DataTable.body())), message="Unable to find list body")
    pages = WebDriverWait(driver, 3).until(EC.presence_of_all_elements_located((By.XPATH, DataTable.sub_pages())),
                                           message="Unable to find subpages")
    if len(pages) > 1:
        WebDriverWait(driver, 3).until(
            EC.presence_of_element_located((By.XPATH, DataTable.last_page())), message="Unable to find 'go to last "
                                                                                       "page' button")
        try:  # avoid failure because of ghost button
            WebDriverWait(driver, 3).until(
                EC.element_to_be_clickable((By.XPATH, DataTable.destroyed_on_click_button())),
                message="ghost button").click()
        except TimeoutException:
            pass
        WebDriverWait(driver, 3).until(
            EC.element_to_be_clickable((By.XPATH, DataTable.last_page())), message="Unable to click on 'go to the last "
                                                                                   "page' button").click()
        last_pages = WebDriverWait(driver, 3).until(EC.presence_of_all_elements_located((By.XPATH,
                                                                                         DataTable.sub_pages())),
                                                    message="Unable to find subpages")
        last_page = int(last_pages[-1].text)
        for page in range(last_page):
            try:
                element = WebDriverWait(driver, 2).until(
                    EC.presence_of_element_located((By.XPATH, xpath)))
                return element
            except TimeoutException:
                WebDriverWait(driver, 2).until(
                    EC.element_to_be_clickable((By.XPATH, DataTable.previous_page())),
                    message="Unable to click on 'go to the previous page' button").click()
            except StaleElementReferenceException:
                driver.refresh()
                threading.Event().wait(1)
            except OSError as err:
                print(err)
        return None
    else:
        try:
            element = WebDriverWait(driver, 2).until(
                EC.presence_of_element_located((By.XPATH, xpath)))
            return element
        except TimeoutException:
            return None
        except StaleElementReferenceException:
            driver.refresh()
            threading.Event().wait(1)
            return StaleElementReferenceException
        except OSError as err:
            return err


@threading_wait_until
def find_element_on_datatable_by_condition(driver, element_xpath, xpath_referring_page, condition="is", event=None):
    try:
        assert_that(condition, any_of(equal_to("is"), equal_to("is not")), "Unexpected value for list condition")
        WebDriverWait(driver, 3).until(
            EC.element_to_be_clickable((By.XPATH, xpath_referring_page)),
            message=f"Unable to find {xpath_referring_page} icon on left menu")
        driver.find_element(By.XPATH, xpath_referring_page).click()
        element_on_datatable = find_element_on_datatable(driver, element_xpath)
        if condition == "is" and element_on_datatable is not None:
            event.set()
        elif condition == "is not" and element_on_datatable is None:
            event.set()
        return element_on_datatable
    except TimeoutException:
        print(TimeoutException)
        raise TimeoutException
    except StaleElementReferenceException:
        driver.refresh()
        event.wait(1)
        print(StaleElementReferenceException)
    except OSError as err:
        raise err
