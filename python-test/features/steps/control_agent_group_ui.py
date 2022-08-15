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
agent_group_description_prefix = "test_agent_group_description_"
agent_group_to_be_deleted = "agent_group_to_delete"


@given("the user clicks on new agent group on left menu")
def agent_page(context):
    WebDriverWait(context.driver, 3).until(
        EC.element_to_be_clickable((By.XPATH, LeftMenu.agent_group_menu())))
    context.driver.find_element(By.XPATH, (LeftMenu.agent_group_menu())).click()
    WebDriverWait(context.driver, 10).until(EC.url_to_be(f"{orb_url}/pages/fleet/groups"), message="Orb agent group "
                                                                                                       "page not "
                                                                                                       "available")


@when("a new agent group is created through the UI with {orb_tags} orb tag")
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
        EC.presence_of_all_elements_located((By.XPATH, DataTable.body())), message="Unable to find agent group list body")
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
        try:   
            element = WebDriverWait(driver, 2).until(
            EC.presence_of_all_elements_located((By.XPATH, xpath)))
            return element[0]
        except TimeoutException:
            return None    
    

@when("a new agent group with description is created through the UI with {orb_tags} orb tag")
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
    context.agent_group_description = agent_group_description_prefix + random_string(10)
    input_text_by_xpath(AgentGroupPage.agent_group_description(), context.agent_group_description, context)
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
  
    
@when("delete the agent group using filter by name with {orb_tags} orb tag")
def create_agent_through_the_agent_group_page(context, orb_tags):
    context.orb_tags = create_tags_set(orb_tags)
    WebDriverWait(context.driver, 5).until(
        EC.element_to_be_clickable((By.XPATH, AgentGroupPage.new_agent_group_button())), message="Unable to click on new agent group"
                                                                                       " button").click()
    WebDriverWait(context.driver, 5).until(EC.url_to_be(f"{orb_url}/pages/fleet/groups/add"), message="Orb add"
                                                                                                           "agent group "
                                                                                                           "page not "
                                                                                                           "available")
    context.agent_group_name = agent_group_to_be_deleted + random_string(5)
    input_text_by_xpath(AgentGroupPage.agent_group_name(), context.agent_group_name, context)
    context.agent_group_description = agent_group_description_prefix + random_string(10)
    input_text_by_xpath(AgentGroupPage.agent_group_description(), context.agent_group_description, context)
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
    WebDriverWait(context.driver, 3).until(
        EC.element_to_be_clickable((By.XPATH, DataTable.filter_by()))).click()
    WebDriverWait(context.driver, 3).until(
        EC.presence_of_all_elements_located((By.XPATH, DataTable.filter_by())))
    WebDriverWait(context.driver, 3).until(
        EC.presence_of_all_elements_located((By.XPATH, DataTable.option_list())))
    select_list = WebDriverWait(context.driver, 3).until(
        EC.presence_of_all_elements_located((By.XPATH, DataTable.all_filter_options())))
    select_list[0].click()
    WebDriverWait(context.driver, 3).until(
        EC.presence_of_element_located((By.XPATH, DataTable.filter_by_name_field())))
    input_text_by_xpath(DataTable.filter_by_name_field(), context.agent_group_name, context)
    WebDriverWait(context.driver, 3).until(
        EC.element_to_be_clickable((By.XPATH, DataTable.plus_button()))).click() 
    context.initial_counter_group = check_total_counter(context.driver)
    WebDriverWait(context.driver, 5).until( 
        EC.element_to_be_clickable((By.XPATH, DataTable.trash_icon()))).click()
    input_text_by_xpath(AgentGroupPage.delete_agent_confirmation_field(), context.agent_group_name, context)
    WebDriverWait(context.driver, 5).until(
        EC.element_to_be_clickable((By.XPATH, AgentGroupPage.delete_agent_confirmation_title()))).click()
    WebDriverWait(context.driver, 3).until(
        EC.element_to_be_clickable((By.XPATH, AgentGroupPage.delete_agent_confirmation_button()))).click()
    WebDriverWait(context.driver, 3).until(
        EC.text_to_be_present_in_element((By.CSS_SELECTOR, "span.title"), 'Agent Group successfully deleted'))  

    
@then("the agent group is not shown on the datatable")
def check_status_on_orb_ui(context):
    context.driver.get(f"{orb_url}/pages/fleet/groups")
    group = find_element_on_agent_group_datatable(context.driver, DataTable.agent_group(context.agent_group_name))
    assert_that(group, is_(None), "Unable to find group name on the datatable")
    
    
@then("total number was decreased in one unit")
def check_total_counter_final(context):
    final_counter_group = check_total_counter(context.driver)
    assert_that(final_counter_group, equal_to(context.initial_counter_group - 1), 'The counter was not decreased wirh successfully')


def check_total_counter(driver):
    WebDriverWait(driver, 3).until(
        EC.presence_of_element_located((By.XPATH,DataTable.page_count())))
    return int(driver.find_element(By.XPATH, DataTable.page_count()).text.split()[0])
        