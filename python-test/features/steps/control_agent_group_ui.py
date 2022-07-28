from behave import given, when, then, step
from test_config import TestConfig
from selenium.webdriver.support.ui import WebDriverWait
from selenium.webdriver.support import expected_conditions as EC
from selenium.webdriver.common.by import By
from selenium.webdriver.support import *
from selenium.common.exceptions import TimeoutException
from ui_utils import *
from hamcrest import *
from utils import random_string, create_tags_set
from utils import threading_wait_until
from page_objects import *

configs = TestConfig.configs()
orb_url = configs.get('orb_url')
agent_group_name_prefix = "test_agent_group_name_"


@given('the user clicks on new agent group on left menu')
def agent_page(context):
    WebDriverWait(context.driver, 3).until(
        EC.element_to_be_clickable((By.XPATH, LeftMenu.agent_group_menu())))
    context.driver.find_element(By.XPATH, (LeftMenu.agent_group_menu())).click()
    WebDriverWait(context.driver, 10).until(EC.url_to_be(f"{orb_url}/pages/fleet/groups"), message="Orb agent group "
                                                                                                       "page not "
                                                                                                       "available")


@when('a new agent group is created with {orb_tags} orb tag')
def create_agent_through_the_agent_group_page(context, orb_tags):
    context.orb_tags = create_tags_set(orb_tags)
    WebDriverWait(context.driver, 5).until(
        EC.element_to_be_clickable((By.XPATH, AgentGroupPage.new_agent_group_button())), message="Unable to click on new agent group"
                                                                                       " button").click()
    WebDriverWait(context.driver, 5).until(EC.url_to_be(f"{orb_url}/pages/fleet/groups/add"), message="Orb add"
                                                                                                           "agent group "
                                                                                                           "page not "
                                                                                                           "available")
    context.agent_group_name = agent_group_name_prefix + random_string(10)
    input_text_by_xpath(AgentGroupPage.agent_group_name(), context.agent_group_name, context)
    WebDriverWait(context.driver, 3).until(
        EC.element_to_be_clickable((By.XPATH, UtilButton.next_button()))).click()
    for tag_key, tag_value in context.orb_tags.items():
        input_text_by_xpath(AgentGroupPage.agent_group_tag_key(), tag_key, context)
        input_text_by_xpath(AgentGroupPage.agent_group_tag_value(), tag_value, context)
        WebDriverWait(context.driver, 3).until(
            EC.element_to_be_clickable((By.XPATH, AgentGroupPage.agent_group_add_tag_button()))).click()
    WebDriverWait(context.driver, 3).until(
        EC.element_to_be_clickable((By.XPATH, UtilButton.next_button()))).click()
    WebDriverWait(context.driver, 3).until(
        EC.element_to_be_clickable((By.XPATH, UtilButton.save_button()))).click()
    WebDriverWait(context.driver, 3).until(
        EC.text_to_be_present_in_element((By.CSS_SELECTOR, "span.title"), 'Agent Group successfully created'))


#@then("the agent group list and the agents view should display agent's status as {status} within {time_to_wait} seconds")
#def check_status_on_orb_ui(context, status, time_to_wait):
#    context.driver.get(f"{orb_url}/pages/fleet/agents")
    # agent_xpath = f"//span[contains(text(), '{context.agent_name}')]/ancestor::div[contains(@class, " \
    #               f"'datatable-row-center')]/descendant::i[contains(@class, " \
    #               f"'fa fa-circle')]/ancestor::span[contains(@class, 'ng-star-inserted')]"
#    agent_group_status_datatable = check_agent_group_status_on_orb_ui(context.driver, DataTable.agent_group_status(context.agent_group_name),
#                                                          status, timeout=time_to_wait)
#    assert_that(agent_group_status_datatable, is_not(None), f"Unable to find status of the agent group: {context.agent_group_name}"
#                                                      f" on datatable")
#    assert_that(agent_group_status_datatable, equal_to(status), f"Agent {context.agent['id']} status failed on Agents list")
#    WebDriverWait(context.driver, 3).until(
#        EC.element_to_be_clickable((By.XPATH, DataTable.agent(context.agent_group_name))),
#        message="Unable to click on agent name").click()
#   agent_view_status = WebDriverWait(context.driver, 3).until(
#        EC.presence_of_element_located(
#           (By.XPATH, AgentsPage.agent_group_status())), message="Unable to find agent status on agent view page").text
#    #agent_view_status = agent_view_status.replace(".", "")
#    #assert_that(agent_view_status, equal_to(status), f"Agent {context.agent['id']} status failed")

#@threading_wait_until
#def check_agent_group_status_on_orb_ui(driver, agent_xpath, status, event=None):
    """

    :param driver: webdriver running
    :param (str) agent_xpath: xpath of the agent whose status need to be checked
    :param (str) status: agent expected status
    :param event: threading.event
    :return: web element refereed to the agent
    """
#    agent_group_status_datatable = find_element_on_agent_group_datatable(driver, agent_xpath)
#    if agent_group_status_datatable is not None and agent_group_status_datatable.text == status:
#        event.set()
#        return agent_group_status_datatable.text
#    driver.refresh()
#    return agent_group_status_datatable


@then("the new agent group is shown on the datatable")
def check_status_on_orb_ui(context):
    context.driver.get(f"{orb_url}/pages/fleet/groups")
    group = find_element_on_agent_group_datatable(context.driver, DataTable.agent_group(context.agent_group_name))
    assert_that(group, is_not(None), "Unable to find group name on the datatable")
       
    
def find_element_on_agent_group_datatable(driver, xpath):
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
                    EC.presence_of_all_elements_located((By.XPATH, xpath)))
                return element[0]
            except TimeoutException:
                WebDriverWait(driver, 2).until(
                    EC.element_to_be_clickable((By.XPATH, DataTable.previous_page())),
                    message="Unable to click on 'go to the previous page' button").click()
            except OSError as err:
                print(err)
        return None
    else:
        element = WebDriverWait(driver, 2).until(
            EC.presence_of_all_elements_located((By.XPATH, xpath)))
    return element[0]
