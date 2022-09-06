from behave import when, then
from test_config import TestConfig
from ui_utils import *
from hamcrest import *
from utils import random_string, create_tags_set
from page_objects import *
from selenium.common.exceptions import TimeoutException, StaleElementReferenceException

configs = TestConfig.configs()
orb_url = configs.get('orb_url')
agent_group_name_prefix = "test_agent_group_name_"
agent_group_description_prefix = "test_agent_group_description_"
agent_group_to_be_deleted = "agent_group_to_delete"


@when('a new agent group is created through the UI with {orb_tags} orb tag')
def create_agent_group_through_the_agent_group_page(context, orb_tags):
    context.orb_tags = create_tags_set(orb_tags)
    WebDriverWait(context.driver, 5).until(
        EC.element_to_be_clickable((By.XPATH, AgentGroupPage.new_agent_group_button())),
        message="Unable to click on new agent group"
                " button").click()
    WebDriverWait(context.driver, 5).until(EC.url_to_be(f"{orb_url}/pages/fleet/groups/add"), message="Orb add"
                                                                                                      "agent group "
                                                                                                      "page not "
                                                                                                      "available")
    context.agent_group_name = agent_group_name_prefix + random_string(10)
    input_text_by_xpath(AgentGroupPage.agent_group_name(), context.agent_group_name, context.driver)
    WebDriverWait(context.driver, 3).until(
        EC.element_to_be_clickable((By.XPATH, UtilButton.next_button()))).click()
    for tag_key, tag_value in context.orb_tags.items():
        input_text_by_xpath(AgentGroupPage.agent_group_tag_key(), tag_key, context.driver)
        input_text_by_xpath(AgentGroupPage.agent_group_tag_value(), tag_value, context.driver)
        WebDriverWait(context.driver, 3).until(
            EC.element_to_be_clickable((By.XPATH, AgentGroupPage.agent_group_add_tag_button()))).click()
    WebDriverWait(context.driver, 3).until(
        EC.element_to_be_clickable((By.XPATH, UtilButton.next_button()))).click()
    WebDriverWait(context.driver, 3).until(
        EC.element_to_be_clickable((By.XPATH, UtilButton.save_button()))).click()
    WebDriverWait(context.driver, 3).until(
        EC.text_to_be_present_in_element((By.CSS_SELECTOR, "span.title"), 'Agent Group successfully created'))


@then("the new agent group {condition} shown on the datatable")
def check_presence_of_group_on_orb_ui(context, condition):
    group = find_element_on_datatable_by_condition(context.driver, DataTable.agent_group(context.agent_group_name),
                                                   LeftMenu.agent_group(), condition)
    if condition == "is":
        assert_that(group, is_not(None), "Unable to find group name on the datatable")
    else:
        assert_that(group, is_(None), f"{context.agent_group_name} found on group datable")


@when("a new agent group with description is created through the UI with {orb_tags} orb tag")
def create_agent_group_with_description_through_the_agent_group_page(context, orb_tags):
    context.orb_tags = create_tags_set(orb_tags)
    WebDriverWait(context.driver, 5).until(
        EC.element_to_be_clickable((By.XPATH, AgentGroupPage.new_agent_group_button())), message="Unable to click on "
                                                                                                 "new agent group "
                                                                                                 " button").click()
    WebDriverWait(context.driver, 5).until(EC.url_to_be(f"{orb_url}/pages/fleet/groups/add"), message="Orb add"
                                                                                                      "agent group "
                                                                                                      "page not "
                                                                                                      "available")
    context.agent_group_name = agent_group_name_prefix + random_string(10)
    input_text_by_xpath(AgentGroupPage.agent_group_name(), context.agent_group_name, context.driver)
    context.agent_group_description = agent_group_description_prefix + random_string(10)
    input_text_by_xpath(AgentGroupPage.agent_group_description(), context.agent_group_description, context.driver)
    WebDriverWait(context.driver, 3).until(
        EC.element_to_be_clickable((By.XPATH, UtilButton.next_button()))).click()
    for tag_key, tag_value in context.orb_tags.items():
        input_text_by_xpath(AgentGroupPage.agent_group_tag_key(), tag_key, context.driver)
        input_text_by_xpath(AgentGroupPage.agent_group_tag_value(), tag_value, context.driver)
        WebDriverWait(context.driver, 3).until(
            EC.element_to_be_clickable((By.XPATH, AgentGroupPage.agent_group_add_tag_button()))).click()
    WebDriverWait(context.driver, 3).until(
        EC.element_to_be_clickable((By.XPATH, UtilButton.next_button()))).click()
    WebDriverWait(context.driver, 3).until(
        EC.element_to_be_clickable((By.XPATH, UtilButton.save_button()))).click()
    WebDriverWait(context.driver, 3).until(
        EC.text_to_be_present_in_element((By.CSS_SELECTOR, "span.title"), 'Agent Group successfully created'))
    context.initial_counter = check_total_counter(context.driver)


@when("delete the agent group using filter by name with {orb_tags} orb tag")
def delete_agent_through_the_agent_group_page(context, orb_tags):
    create_agent_group_with_description_through_the_agent_group_page(context, orb_tags)
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
    input_text_by_xpath(DataTable.filter_by_name_field(), context.agent_group_name, context.driver)
    WebDriverWait(context.driver, 3).until(
        EC.element_to_be_clickable((By.XPATH, DataTable.plus_button()))).click()
    WebDriverWait(context.driver, 5).until(
        EC.element_to_be_clickable((By.XPATH, DataTable.trash_icon()))).click()
    input_text_by_xpath(AgentGroupPage.delete_agent_group_confirmation_field(), context.agent_group_name, context.driver)
    WebDriverWait(context.driver, 5).until(
        EC.element_to_be_clickable((By.XPATH, AgentGroupPage.delete_agent_group_confirmation_title()))).click()
    WebDriverWait(context.driver, 3).until(
        EC.element_to_be_clickable((By.XPATH, AgentGroupPage.delete_agent_group_confirmation_button()))).click()
    WebDriverWait(context.driver, 3).until(
        EC.text_to_be_present_in_element((By.CSS_SELECTOR, "span.title"), 'Agent Group successfully deleted'))
    WebDriverWait(context.driver, 3).until(EC.element_to_be_clickable((By.XPATH, UtilButton.clear_all_filters())),
                                           "Unable to clear all filters").click()


@then("total number was decreased in one unit")
def check_total_counter_final(context):
    final_counter_group = check_total_counter(context.driver)
    assert_that(final_counter_group, equal_to(context.initial_counter - 1),
                'The counter was not decrease with successful')


def check_total_counter(driver):
    threading.Event().wait(3)
    WebDriverWait(driver, 3).until(
        EC.presence_of_element_located((By.XPATH, DataTable.page_count())))
    return int(driver.find_element(By.XPATH, DataTable.page_count()).text.split()[0])



@when("update the agent group using filter by name with {orb_tags} orb tag")
def update_an_agent_group_by_name_through_the_agent_group_page(context, orb_tags):
    create_agent_group_through_the_agent_group_page(context, orb_tags)
    context.initial_counter_datatable = check_total_counter(context.driver)
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
    input_text_by_xpath(DataTable.filter_by_name_field(), context.agent_group_name, context.driver)
    WebDriverWait(context.driver, 3).until(
        EC.element_to_be_clickable((By.XPATH, DataTable.plus_button()))).click()
    WebDriverWait(context.driver, 5).until(
        EC.element_to_be_clickable((By.XPATH, DataTable.edit_icon()))).click()
    WebDriverWait(context.driver, 3).until(
        EC.element_to_be_clickable((By.XPATH, AgentGroupPage.agent_group_name()))).clear()
    context.agent_group_name = agent_group_name_prefix + random_string(5)
    input_text_by_xpath(AgentGroupPage.agent_group_name(), context.agent_group_name, context.driver)
    WebDriverWait(context.driver, 3).until(
        EC.element_to_be_clickable((By.XPATH, UtilButton.next_button()))).click()
    WebDriverWait(context.driver, 3).until(
        EC.element_to_be_clickable((By.XPATH, UtilButton.next_button()))).click()
    WebDriverWait(context.driver, 3).until(
        EC.element_to_be_clickable((By.XPATH, UtilButton.save_button()))).click()
    WebDriverWait(context.driver, 3).until(
        EC.text_to_be_present_in_element((By.CSS_SELECTOR, "span.title"), 'Agent Group successfully updated'))
    context.initial_counter = check_total_counter(context.driver)
    WebDriverWait(context.driver, 3).until(
        EC.element_to_be_clickable((By.XPATH, DataTable.close_option_selected()))).click()
    