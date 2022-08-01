from behave import given, then, step
from ui_utils import input_text_by_xpath
from utils import threading_wait_until
from control_plane_agents import agent_name_prefix
from selenium.webdriver.support.ui import WebDriverWait
from selenium.webdriver.support import expected_conditions as EC
from selenium.webdriver.common.by import By
from selenium.common.exceptions import TimeoutException
from utils import random_string, create_tags_set
from test_config import TestConfig
from hamcrest import *
from page_objects import *

configs = TestConfig.configs()
orb_url = configs.get('orb_url')
AGENTS_URL = f"{orb_url}/pages/fleet/agents"


@given("that fleet Management is clickable on ORB Menu")
def expand_fleet_management(context):
    context.driver.find_element(By.XPATH, (LeftMenu.agents())).click()


@given('that Agents is clickable on ORB Menu')
def agent_page(context):
    WebDriverWait(context.driver, 3).until(
        EC.element_to_be_clickable((By.XPATH, LeftMenu.agents())), message="Unable to find agent icon on left menu")
    context.driver.find_element(By.XPATH, (LeftMenu.agents())).click()
    WebDriverWait(context.driver, 5).until(EC.url_to_be(AGENTS_URL), message="Orb agents page not available")


@step("that the user is on the orb Agent page")
def orb_page(context):
    expand_fleet_management(context)
    agent_page(context)
    current_url = context.driver.current_url
    assert_that(current_url, equal_to(AGENTS_URL), "user not enabled to access orb login page")


@step("a new agent is created through the UI with {orb_tags} orb tag(s)")
def create_agent_through_the_agents_page(context, orb_tags):
    context.orb_tags = create_tags_set(orb_tags)
    WebDriverWait(context.driver, 3).until(
        EC.element_to_be_clickable((By.XPATH, AgentsPage.new_agent_button())), message="Unable to click on new agent"
                                                                                       " button").click()
    WebDriverWait(context.driver, 5).until(EC.url_to_be(f"{AGENTS_URL}/add"), message="Orb add agents page not "
                                                                                      "available")
    context.agent_name = agent_name_prefix + random_string(10)
    input_text_by_xpath(AgentsPage.agent_name(), context.agent_name, context)
    WebDriverWait(context.driver, 3).until(
        EC.element_to_be_clickable((By.XPATH, UtilButton.next_button())), message="Unable to click on next "
                                                                                  "button (page 1)").click()
    for tag_key, tag_value in context.orb_tags.items():
        input_text_by_xpath(AgentsPage.agent_tag_key(), tag_key, context)
        input_text_by_xpath(AgentsPage.agent_tag_value(), tag_value, context)
        WebDriverWait(context.driver, 3).until(
            EC.element_to_be_clickable((By.XPATH, AgentsPage.agent_add_tag_button())), message="Unable to click on add"
                                                                                               " tag button").click()
    WebDriverWait(context.driver, 3).until(
        EC.element_to_be_clickable((By.XPATH, UtilButton.next_button())), message="Unable to click on next "
                                                                                  "button (page 2)").click()
    WebDriverWait(context.driver, 3).until(
        EC.element_to_be_clickable((By.XPATH, UtilButton.save_button())), message="Unable to click on save "
                                                                                  "button").click()
    WebDriverWait(context.driver, 3).until(
        EC.text_to_be_present_in_element((By.CSS_SELECTOR, "span.title"), 'Agent successfully created'),
        message="Confirmation span of agent creation not displayed")
    context.agent_key = \
        WebDriverWait(context.driver, 3).until(EC.presence_of_element_located((By.XPATH, AgentsPage.agent_key())),
                                               message="Agent key not displayed").text
    agent_provisioning_command = \
        WebDriverWait(context.driver, 3).until(
            EC.presence_of_element_located((By.XPATH, AgentsPage.agent_provisioning_command())),
            message="Provisioning command not displayed").text

    context.agent_provisioning_command = agent_provisioning_command.replace("\n\n", " ")
    WebDriverWait(context.driver, 3).until(
        EC.element_to_be_clickable((By.XPATH, UtilButton.close_button())), message="Unable to click on close "
                                                                                   "button").click()
    agent = find_element_on_agent_datatable(context.driver, DataTable.agent(context.agent_name))
    assert_that(agent, is_not(None), f"Unable to find the agent: {context.agent_name}")
    agent.click()
    context.agent = dict()
    context.agent['id'] = WebDriverWait(context.driver, 3).until(
        EC.presence_of_element_located((By.XPATH, AgentsPage.agent_view_id())), message="Agent id not displayed").text
    assert_that(context.agent['id'],
                matches_regexp(r'[a-zA-Z0-9]{8}\-[a-zA-Z0-9]{4}\-[a-zA-Z0-9]{4}\-[a-zA-Z0-9]{4}\-[a-zA-Z0-9]{12}'),
                f"Failed to get agent id {context.agent['id']}")
    context.agent['name'] = context.agent_name


@then("the agents list and the agents view should display agent's status as {status} within {time_to_wait} seconds")
def check_status_on_orb_ui(context, status, time_to_wait):
    context.driver.get(f"{orb_url}/pages/fleet/agents")
    # agent_xpath = f"//span[contains(text(), '{context.agent_name}')]/ancestor::div[contains(@class, " \
    #               f"'datatable-row-center')]/descendant::i[contains(@class, " \
    #               f"'fa fa-circle')]/ancestor::span[contains(@class, 'ng-star-inserted')]"
    agent_status_datatable = check_agent_status_on_orb_ui(context.driver, DataTable.agent_status(context.agent_name),
                                                          status, timeout=time_to_wait)
    assert_that(agent_status_datatable, is_not(None), f"Unable to find status of the agent: {context.agent_name}"
                                                      f" on datatable")
    assert_that(agent_status_datatable, equal_to(status), f"Agent {context.agent['id']} status failed on Agents list")
    WebDriverWait(context.driver, 3).until(
        EC.element_to_be_clickable((By.XPATH, DataTable.agent(context.agent_name))),
        message="Unable to click on agent name").click()
    agent_view_status = WebDriverWait(context.driver, 3).until(
        EC.presence_of_element_located(
            (By.XPATH, AgentsPage.agent_status())), message="Unable to find agent status on agent view page").text
    agent_view_status = agent_view_status.replace(".", "")
    assert_that(agent_view_status, equal_to(status), f"Agent {context.agent['id']} status failed")


@threading_wait_until
def check_agent_status_on_orb_ui(driver, agent_xpath, status, event=None):
    """

    :param driver: webdriver running
    :param (str) agent_xpath: xpath of the agent whose status need to be checked
    :param (str) status: agent expected status
    :param event: threading.event
    :return: web element refereed to the agent
    """
    agent_status_datatable = find_element_on_agent_datatable(driver, agent_xpath)
    if agent_status_datatable is not None and agent_status_datatable.text == status:
        event.set()
        return agent_status_datatable.text
    driver.refresh()
    return agent_status_datatable


def find_element_on_agent_datatable(driver, xpath):
    """
    Find element present on agent datatable

    :param driver: webdriver running
    :param (str) xpath: xpath of the element to be found
    :return: web element, if found. None if not found.
    """
    WebDriverWait(driver, 3).until(
        EC.presence_of_all_elements_located((By.XPATH, DataTable.page_count())), message="Unable to find page count")
    WebDriverWait(driver, 3).until(
        EC.presence_of_all_elements_located((By.XPATH, DataTable.body())), message="Unable to find agent list body")
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
            except OSError as err:
                print(err)
        return None
    else:
        element = WebDriverWait(driver, 2).until(
            EC.presence_of_element_located((By.XPATH, xpath)))
    return element
