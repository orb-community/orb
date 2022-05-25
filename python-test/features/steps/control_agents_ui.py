from behave import given, then, step
from ui_utils import input_text_by_xpath
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


@given("that fleet Management is clickable on ORB Menu")
def expand_fleet_management(context):
    context.driver.find_elements_by_xpath(LeftMenu.agents())[0].click()


@given('that Agents is clickable on ORB Menu')
def agent_page(context):
    WebDriverWait(context.driver, 3).until(
        EC.element_to_be_clickable((By.XPATH, LeftMenu.agents())))
    context.driver.find_element_by_xpath(LeftMenu.agents()).click()
    WebDriverWait(context.driver, 5).until(EC.url_to_be(f"{orb_url}/pages/fleet/agents"), message="Orb agents "
                                                                                                       "page not "
                                                                                                       "available")


@step("that the user is on the orb Agent page")
def orb_page(context):
    expand_fleet_management(context)
    agent_page(context)
    current_url = context.driver.current_url
    assert_that(current_url, equal_to(f"{orb_url}/pages/fleet/agents"),
                "user not enabled to access orb login page")


@step("a new agent is created through the UI with {orb_tags} orb tag(s)")
def create_agent_through_the_agents_page(context, orb_tags):
    context.orb_tags = create_tags_set(orb_tags)
    WebDriverWait(context.driver, 3).until(
        EC.element_to_be_clickable((By.XPATH, AgentsPage.new_agent_button()))).click()
    WebDriverWait(context.driver, 5).until(EC.url_to_be(f"{orb_url}/pages/fleet/agents/add"), message="Orb add"
                                                                                                           "agents "
                                                                                                           "page not "
                                                                                                           "available")
    context.agent_name = agent_name_prefix + random_string(10)
    input_text_by_xpath(AgentsPage.agent_name(), context.agent_name, context)
    WebDriverWait(context.driver, 3).until(
        EC.element_to_be_clickable((By.XPATH, UtilButton.next_button()))).click()
    for tag_key, tag_value in context.orb_tags.items():
        input_text_by_xpath(AgentsPage.agent_tag_key(), tag_key, context)
        input_text_by_xpath(AgentsPage.agent_tag_value(), tag_value, context)
        WebDriverWait(context.driver, 3).until(
            EC.element_to_be_clickable((By.XPATH, AgentsPage.agent_add_tag_button()))).click()
    WebDriverWait(context.driver, 3).until(
        EC.element_to_be_clickable((By.XPATH, UtilButton.next_button()))).click()
    WebDriverWait(context.driver, 3).until(
        EC.element_to_be_clickable((By.XPATH, UtilButton.save_button()))).click()
    WebDriverWait(context.driver, 3).until(
        EC.text_to_be_present_in_element((By.CSS_SELECTOR, "span.title"), 'Agent successfully created'))
    context.agent_key = \
        WebDriverWait(context.driver, 3).until(EC.presence_of_all_elements_located((By.XPATH,
                                                                                    AgentsPage.agent_key())))[0].text
    agent_provisioning_command = \
        WebDriverWait(context.driver, 3).until(
            EC.presence_of_all_elements_located((By.XPATH, AgentsPage.agent_provisioning_command())))[0].text

    context.agent_provisioning_command = agent_provisioning_command.replace("\n\n", " ")
    WebDriverWait(context.driver, 3).until(
        EC.presence_of_all_elements_located((By.XPATH, UtilButton.close_button())))[0].click()
    WebDriverWait(context.driver, 3).until(
        EC.presence_of_all_elements_located((By.XPATH, f"//div[contains(@class, 'agent-name') and contains(text(),"
                                                       f"'{context.agent_name}')]")))[0].click()
    context.agent = dict()
    context.agent['id'] = WebDriverWait(context.driver, 3).until(
        EC.presence_of_all_elements_located((By.XPATH, AgentsPage.agent_view_id())))[0].text


@threading_wait_until
def check_agent_status_on_orb_ui(driver, agent_xpath, status, event=None):
    list_of_datatable_body_cell = WebDriverWait(driver, 3).until(EC.presence_of_all_elements_located([By.XPATH,
                                                                                                      agent_xpath]))
    if list_of_datatable_body_cell[0].text.splitlines()[1] == status:
        event.set()
        return list_of_datatable_body_cell
    driver.refresh()
    return list_of_datatable_body_cell


@then("the agents list and the agents view should display agent's status as {status} within {time_to_wait} seconds")
def check_status_on_orb_ui(context, status, time_to_wait):
    context.driver.get(f"{orb_url}/pages/fleet/agents")
    agent_xpath = f"//div[contains(text(), '{context.agent_name}')]/ancestor::datatable-body-row/descendant::" \
                  f"i[contains(@class, 'fa fa-circle')]/ancestor::div[contains(@class, 'ng-star-inserted')]"
    list_of_datatable_body_cell = check_agent_status_on_orb_ui(context.driver, agent_xpath, status,
                                                               timeout=time_to_wait)

    assert_that(list_of_datatable_body_cell[0].text.splitlines()[1], equal_to(status),
                f"Agent {context.agent['id']} status failed")
    assert_that(list_of_datatable_body_cell[1].text, equal_to(status))
    WebDriverWait(context.driver, 3).until(
        EC.presence_of_all_elements_located((By.XPATH, f"//div[contains(@class, 'agent-name') and contains(text(),"
                                                       f"'{context.agent_name}')]")))[0].click()
    agent_view_status = WebDriverWait(context.driver, 3).until(
        EC.presence_of_all_elements_located(
            (By.XPATH, AgentsPage.agent_status())))[0].text
    agent_view_status = agent_view_status.replace(".", "")
    assert_that(agent_view_status, equal_to(status), f"Agent {context.agent['id']} status failed")
