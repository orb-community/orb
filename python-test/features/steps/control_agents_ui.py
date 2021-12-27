from behave import given, when, then
from ui_utils import input_text_by_xpath
from control_plane_agents import agent_name_prefix, tag_key_prefix, tag_value_prefix
from selenium.webdriver.support.ui import WebDriverWait
from selenium.webdriver.support import expected_conditions as EC
from selenium.webdriver.common.by import By
from utils import random_string
from test_config import TestConfig
import time
from hamcrest import *

configs = TestConfig.configs()
base_orb_url = configs.get('base_orb_url')


@given("that fleet Management is clickable on ORB Menu")
def expand_fleet_management(context):
    context.driver.find_elements_by_xpath("//a[contains(@title, 'Fleet Management')]")[0].click()


@given('that Agents is clickable on ORB Menu')
def agent_page(context):
    WebDriverWait(context.driver, 3).until(
        EC.element_to_be_clickable((By.XPATH, "//a[contains(@title, 'Agents')]")))
    context.driver.find_element_by_xpath("//a[contains(@title, 'Agents')]").click()
    WebDriverWait(context.driver, 5).until(EC.url_to_be(f"{base_orb_url}/pages/fleet/agents"), message="Orb agents "
                                                                                                       "page not "
                                                                                                       "available")


@when("a new agent is created through the UI")
def create_agent_through_the_agents_page(context):
    WebDriverWait(context.driver, 3).until(
        EC.element_to_be_clickable((By.XPATH, "//button[contains(text(), 'New Agent')]"))).click()
    WebDriverWait(context.driver, 5).until(EC.url_to_be(f"{base_orb_url}/pages/fleet/agents/add"), message="Orb add"
                                                                                                           "agents "
                                                                                                           "page not "
                                                                                                           "available")
    context.agent_name = agent_name_prefix + random_string(10)
    context.agent_tag_key = tag_key_prefix + random_string(4)
    context.agent_tag_value = tag_value_prefix + random_string(4)
    input_text_by_xpath("//input[contains(@data-orb-qa-id, 'input#name')]", context.agent_name, context)
    WebDriverWait(context.driver, 3).until(
        EC.element_to_be_clickable((By.XPATH, "//button[contains(text(), 'Next')]"))).click()
    input_text_by_xpath("//input[contains(@data-orb-qa-id, 'input#orb_tag_key')]", context.agent_tag_key, context)
    input_text_by_xpath("//input[contains(@data-orb-qa-id, 'input#orb_tag_value')]", context.agent_tag_value, context)
    WebDriverWait(context.driver, 3).until(
        EC.element_to_be_clickable((By.XPATH, "//button[contains(@data-orb-qa-id, 'button#addTag')]"))).click()
    WebDriverWait(context.driver, 3).until(
        EC.element_to_be_clickable((By.XPATH, "//button[contains(text(), 'Next')]"))).click()
    WebDriverWait(context.driver, 3).until(
        EC.element_to_be_clickable((By.XPATH, "//button[contains(text(), 'Save')]"))).click()
    WebDriverWait(context.driver, 3).until(
        EC.text_to_be_present_in_element((By.CSS_SELECTOR, "span.title"), 'Agent successfully created'))
    agent_key_xpath = "//label[contains(text(), 'Agent Key')]/following::pre[1]"
    context.agent_key = \
        WebDriverWait(context.driver, 3).until(EC.presence_of_all_elements_located((By.XPATH, agent_key_xpath)))[0].text
    agent_provisioning_command_xpath = "//label[contains(text(), 'Provisioning Command')]/following::pre[1]"
    agent_provisioning_command = \
        WebDriverWait(context.driver, 3).until(
            EC.presence_of_all_elements_located((By.XPATH, agent_provisioning_command_xpath)))[0].text

    context.agent_provisioning_command = agent_provisioning_command.replace("\n\n", " ")
    WebDriverWait(context.driver, 3).until(
        EC.presence_of_all_elements_located((By.XPATH, "//span[contains(@class, 'nb-close')]")))[0].click()
    WebDriverWait(context.driver, 3).until(
        EC.presence_of_all_elements_located((By.XPATH, f"//div[contains(@class, 'agent-name') and contains(text(),"
                                                       f"{context.agent_name})]")))[0].click()
    context.agent = dict()
    context.agent['id'] = WebDriverWait(context.driver, 3).until(
        EC.presence_of_all_elements_located((By.XPATH, "//label[contains(text(), 'Agent ID')]/following::p")))[0].text


@then("the agents list and the agents view should display agent's status as {status} within {time_to_wait} seconds")
def check_status_on_orb_ui(context, status, time_to_wait):
    context.driver.get(f"{base_orb_url}/pages/fleet/agents")

    time_waiting = 0
    sleep_time = 0.5
    timeout = int(time_to_wait)
    list_of_datatable_body_cell = list()

    while time_waiting < timeout:
        list_of_datatable_body_cell = WebDriverWait(context.driver, 3).until(
            EC.presence_of_all_elements_located([By.XPATH, f"//div[contains(text(), {context.agent_name})]/ancestor"
                                                           f"::datatable-body-row/descendant::i[contains(@class, "
                                                           f"'fa fa-circle')]/ancestor::div[contains(@class, "
                                                           f"'ng-star-inserted')]"]))
        if list_of_datatable_body_cell[0].text.splitlines()[1] == status:
            break
        time.sleep(sleep_time)
        time_waiting += sleep_time
        context.driver.refresh()
    assert_that(list_of_datatable_body_cell[0].text.splitlines()[1], equal_to(status),
                f"Agent {context.agent['id']} status failed")
    assert_that(list_of_datatable_body_cell[1].text, equal_to(status))
    WebDriverWait(context.driver, 3).until(
        EC.presence_of_all_elements_located((By.XPATH, f"//div[contains(@class, 'agent-name') and contains(text(),"
                                                       f"{context.agent_name})]")))[0].click()
    agent_view_status = WebDriverWait(context.driver, 3).until(
        EC.presence_of_all_elements_located(
            (By.XPATH, "//label[contains(text(), 'Health Status')]/following::p")))[0].text
    assert_that(agent_view_status, equal_to(status), f"Agent {context.agent['id']} status failed")
