from behave import when, then, step
from test_config import TestConfig
from ui_utils import *
from hamcrest import *
from utils import random_string, create_tags_set
from page_objects import *
from control_plane_agent_groups import agent_group_name_prefix, agent_group_description, list_agent_groups

configs = TestConfig.configs()
orb_url = configs.get('orb_url')
agent_group_to_be_deleted = "agent_group_to_delete"


@when('a new agent group is created through the UI with {orb_tags} orb tag')
def create_agent_group_through_the_agent_group_page_with_tags(context, orb_tags):
    context.orb_tags = create_tags_set(orb_tags)
    button_was_clicked = button_click_by_xpath(AgentGroupPage.new_agent_group_button(), context.driver,
                                               "Unable to click on new agent group button")
    assert_that(button_was_clicked, equal_to(True), "Unable to click on new agent group button")
    WebDriverWait(context.driver, 5).until(EC.url_to_be(f"{orb_url}/pages/fleet/groups/add"), message="Orb add"
                                                                                                      "agent group "
                                                                                                      "page not "
                                                                                                      "available")
    context.agent_group_name = agent_group_name_prefix + random_string(10)
    create_group_via_UI(context.agent_group_name, context.orb_tags, context.driver, context.token, context.agent_groups)


@step('a new agent group is created through the UI with same tags as the agent')
def create_agent_group_through_the_agent_group_page(context):
    context.agent_group_name = agent_group_name_prefix + random_string(10)
    button_was_clicked = button_click_by_xpath(AgentGroupPage.new_agent_group_button(), context.driver,
                                               "Unable to click on new agent group button")
    assert_that(button_was_clicked, equal_to(True), "Unable to click on new agent group button")
    tags_in_agent = context.agent["orb_tags"]
    if context.agent["agent_tags"] is not None:
        tags_in_agent.update(context.agent["agent_tags"])
    tags_keys = tags_in_agent.keys()
    context.orb_tags = tags_in_agent

    assert_that(len(tags_keys), greater_than(0), f"Unable to create group without tags. Tags:{tags_in_agent}. "
                                                 f"Agent:{context.agent}")
    create_group_via_UI(context.agent_group_name, context.orb_tags, context.driver, context.token, context.agent_groups)


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
    button_was_clicked = button_click_by_xpath(AgentGroupPage.new_agent_group_button(), context.driver,
                                               "Unable to click on new agent group button")
    assert_that(button_was_clicked, equal_to(True), "Unable to click on new agent group button")
    WebDriverWait(context.driver, 5).until(EC.url_to_be(f"{orb_url}/pages/fleet/groups/add"), message="Orb add"
                                                                                                      "agent group "
                                                                                                      "page not "
                                                                                                      "available")
    context.agent_group_name = agent_group_name_prefix + random_string(10)
    context.agent_group_description = agent_group_description + random_string(10)
    create_group_via_UI(context.agent_group_name, context.orb_tags, context.driver, context.token, context.agent_groups,
                        description=context.agent_group_description)
    context.initial_counter = check_total_counter(context.driver)


@when("delete the agent group using filter by name with {orb_tags} orb tag")
def delete_agent_through_the_agent_group_page(context, orb_tags):
    create_agent_group_with_description_through_the_agent_group_page(context, orb_tags)
    button_was_clicked = button_click_by_xpath(DataTable.filter_by(), context.driver,
                                               "Unable to click on filter group button")
    assert_that(button_was_clicked, equal_to(True), "Unable to click on filter group button")
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
    button_was_clicked = button_click_by_xpath(DataTable.plus_button(), context.driver,
                                               "Unable to click on plus button on group page")
    assert_that(button_was_clicked, equal_to(True), "Unable to click on plus button on group page")
    button_was_clicked = button_click_by_xpath(DataTable.trash_icon(), context.driver,
                                               "Unable to click on trash icon on group page")
    assert_that(button_was_clicked, equal_to(True), "Unable to click on trash icon on group page")
    input_text_by_xpath(AgentGroupPage.delete_agent_group_confirmation_field(), context.agent_group_name,
                        context.driver)
    button_was_clicked = button_click_by_xpath(AgentGroupPage.delete_agent_group_confirmation_title(), context.driver,
                                               "Unable to click on delete group on group page")
    assert_that(button_was_clicked, equal_to(True), "Unable to click on delete group on group page")
    button_was_clicked = button_click_by_xpath(AgentGroupPage.delete_agent_group_confirmation_button(), context.driver,
                                               "Unable to click on delete group confirmation on group page")
    assert_that(button_was_clicked, equal_to(True), "Unable to click on delete group confirmation on group page")
    WebDriverWait(context.driver, 3).until(
        EC.text_to_be_present_in_element((By.CSS_SELECTOR, "span.title"), 'Agent Group successfully deleted'))
    button_was_clicked = button_click_by_xpath(UtilButton.clear_all_filters(), context.driver,
                                               "Unable to clear all filters")
    assert_that(button_was_clicked, equal_to(True), "Unable to clear all filters")


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
    create_agent_group_through_the_agent_group_page_with_tags(context, orb_tags)
    context.initial_counter_datatable = check_total_counter(context.driver)
    button_was_clicked = button_click_by_xpath(DataTable.filter_by(), context.driver,
                                               "Unable to click on filter on group page")
    assert_that(button_was_clicked, equal_to(True), "Unable to click on filter on group page")
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
    button_was_clicked = button_click_by_xpath(DataTable.plus_button(), context.driver,
                                               "Unable to click on plus button on group page")
    assert_that(button_was_clicked, equal_to(True), "Unable to click on plus button on group page")
    button_was_clicked = button_click_by_xpath(DataTable.edit_icon(), context.driver,
                                               "Unable to click on edit button on group page")
    assert_that(button_was_clicked, equal_to(True), "Unable to click on edit button on group page")
    WebDriverWait(context.driver, 3).until(
        EC.element_to_be_clickable((By.XPATH, AgentGroupPage.agent_group_name()))).clear()
    context.agent_group_name = agent_group_name_prefix + "upd" + random_string(5)
    input_text_by_xpath(AgentGroupPage.agent_group_name(), context.agent_group_name, context.driver)
    button_was_clicked = button_click_by_xpath(UtilButton.next_button(), context.driver,
                                               "Unable to click on next button on group page")
    assert_that(button_was_clicked, equal_to(True), "Unable to click on next button on group page")
    button_was_clicked = button_click_by_xpath(UtilButton.next_button(), context.driver,
                                               "Unable to click on next button on group page (2)")
    assert_that(button_was_clicked, equal_to(True), "Unable to click on next button on group page (2)")
    button_was_clicked = button_click_by_xpath(UtilButton.save_button(), context.driver,
                                               "Unable to click on save button on group page")
    assert_that(button_was_clicked, equal_to(True), "Unable to click on save button on group page")
    WebDriverWait(context.driver, 3).until(
        EC.text_to_be_present_in_element((By.CSS_SELECTOR, "span.title"), 'Agent Group successfully updated'))
    context.initial_counter = check_total_counter(context.driver)
    button_was_clicked = button_click_by_xpath(DataTable.close_option_selected(), context.driver,
                                               "Unable to click on close option button on group page")
    assert_that(button_was_clicked, equal_to(True), "Unable to click on close option button on group page")


def create_group_via_UI(name, orb_tags, driver, token, existent_groups, description=None, time_to_wait_until=5):
    assert_that(str(driver.current_url), equal_to(f"{orb_url}/pages/fleet/groups/add"), "Not possible to create a "
                                                                                        "group because the driver is "
                                                                                        "not on the group add page")
    input_text_by_xpath(AgentGroupPage.agent_group_name(), name, driver)
    if description is not None:
        input_text_by_xpath(AgentGroupPage.agent_group_description(), agent_group_description, driver)
    button_was_clicked = button_click_by_xpath(UtilButton.next_button(), driver,
                                               "Unable to click on next button on group page")
    assert_that(button_was_clicked, equal_to(True), "Unable to click on next button on group page")
    for tag_key, tag_value in orb_tags.items():
        input_text_by_xpath(AgentGroupPage.agent_group_tag_key(), tag_key, driver)
        input_text_by_xpath(AgentGroupPage.agent_group_tag_value(), tag_value, driver)
        button_was_clicked = button_click_by_xpath(AgentGroupPage.agent_group_add_tag_button(), driver,
                                                   "Unable to click on add tags button on group page")
        assert_that(button_was_clicked, equal_to(True), "Unable to click on add tags button on group page")
    button_was_clicked = button_click_by_xpath(UtilButton.next_button(), driver,
                                               "Unable to click on next button on group page")
    assert_that(button_was_clicked, equal_to(True), "Unable to click on next tags button on group page")
    button_was_clicked = button_click_by_xpath(UtilButton.save_button(), driver,
                                               "Unable to click on save button on group page")
    assert_that(button_was_clicked, equal_to(True), "Unable to click on save tags button on group page")
    WebDriverWait(driver, time_to_wait_until).until(
        EC.text_to_be_present_in_element((By.CSS_SELECTOR, "span.title"), 'Agent Group successfully created'))
    all_groups = list_agent_groups(token)
    for group in all_groups:
        if group['name'] == name:
            existent_groups[group['id']] = group['name']
            break
