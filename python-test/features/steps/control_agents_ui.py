from behave import given, then, step
from ui_utils import input_text_by_xpath, find_element_on_datatable
from utils import threading_wait_until
from control_plane_agents import agent_name_prefix
from selenium.webdriver.support.ui import WebDriverWait
from selenium.webdriver.support import expected_conditions as EC
from selenium.webdriver.common.by import By
from utils import random_string, create_tags_set
from test_config import TestConfig
from hamcrest import *
from page_objects import *

configs = TestConfig.configs()
orb_url = configs.get('orb_url')


@given("the user clicks on {element} on left menu")
def click_element_left_menu(context, element):
    dict_elements = {"Agents": LeftMenu.agents(), "Agent Groups": LeftMenu.agent_group(),
                     "Policy Management": LeftMenu.policies(), "Sink Management": LeftMenu.sinks()}
    WebDriverWait(context.driver, 3).until(
        EC.element_to_be_clickable((By.XPATH, dict_elements[element])), message=f"Unable to find {element} "
                                                                                f"icon on left menu")
    context.driver.find_element(By.XPATH, dict_elements[element]).click()


@step("that the user is on the orb {element} page")
def check_which_orb_page(context, element):
    dict_pages = {"Agents": OrbPagesUrl.Agents(orb_url), "Agent Groups": OrbPagesUrl.AgentGroups(orb_url),
                  "Policies": OrbPagesUrl.Policies(orb_url), "Sinks": OrbPagesUrl.Sinks(orb_url)}
    dict_elements = {"Agents": LeftMenu.agents(), "Agent Groups": LeftMenu.agent_group(),
                     "Policy Management": LeftMenu.policies(), "Sink Management": LeftMenu.sinks()}
    click_element_left_menu(context, dict_elements[element])
    WebDriverWait(context.driver, 5).until(EC.url_to_be(dict_pages[element]),
                                           message=f"Orb {element} page not available")
    current_url = context.driver.current_url
    assert_that(current_url, equal_to(dict_pages[element]), f"user not enabled to access orb {element} page")


@step("a new agent is created through the UI with {orb_tags} orb tag(s)")
def create_agent_through_the_agents_page(context, orb_tags):
    context.orb_tags = create_tags_set(orb_tags)
    WebDriverWait(context.driver, 3).until(
        EC.element_to_be_clickable((By.XPATH, AgentsPage.new_agent_button())), message="Unable to click on new agent"
                                                                                       " button").click()
    WebDriverWait(context.driver, 5).until(EC.url_to_be(f"{OrbPagesUrl.Agents(orb_url)}/add"),
                                           message="Orb add agents page not "
                                                   "available")
    context.agent_name = agent_name_prefix + random_string(10)
    input_text_by_xpath(AgentsPage.agent_name(), context.agent_name, context.driver)
    WebDriverWait(context.driver, 3).until(
        EC.element_to_be_clickable((By.XPATH, UtilButton.next_button())), message="Unable to click on next "
                                                                                  "button (page 1)").click()
    for tag_key, tag_value in context.orb_tags.items():
        input_text_by_xpath(AgentsPage.agent_tag_key(), tag_key, context.driver)
        input_text_by_xpath(AgentsPage.agent_tag_value(), tag_value, context.driver)
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
    agent = find_element_on_datatable(context.driver, DataTable.agent(context.agent_name))
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
    agent_status_datatable = find_element_on_datatable(driver, agent_xpath)
    if agent_status_datatable is not None and agent_status_datatable.text == status:
        event.set()
        return agent_status_datatable.text
    driver.refresh()
    return agent_status_datatable
