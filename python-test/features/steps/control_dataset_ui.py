from behave import step
from test_config import TestConfig
from control_plane_sink import get_sink
from ui_utils import *
from page_objects import *
from control_plane_datasets import dataset_name_prefix, list_datasets
from utils import random_string

configs = TestConfig.configs()
orb_url = configs.get('orb_url')


@step('a dataset is created through the UI')
def create_dataset_through_ui(context):
    policy_id = context.policy['id']
    policy_view_page = f"{orb_url}/pages/datasets/policies/view/{policy_id}"
    context.driver, current_url = go_to_page(policy_view_page, driver=context.driver)
    group_name = list(context.agent_groups.values())[0]
    sinks_names = list()
    for sink_id in context.existent_sinks_id:
        sink = get_sink(context.token, sink_id)
        sinks_names.append(sink['name'])
    button_was_clicked = button_click_by_xpath(PolicyPage.new_dataset_button(), context.driver,
                                               "Unable to click on new dataset button")
    assert_that(button_was_clicked, equal_to(True), "Unable to click on new dataset button")
    button_was_clicked = button_click_by_xpath(DatasetModal.agent_group(), context.driver,
                                               "Unable to click on button of agent groups on dataset modal")
    assert_that(button_was_clicked, equal_to(True), "Unable to click on button of agent groups on dataset modal")
    groups_options = get_selector_options(context.driver)
    ActionChains(context.driver).move_to_element(groups_options[group_name]).perform()
    WebDriverWait(context.driver, time_webdriver_wait).until(
        EC.element_to_be_clickable(groups_options[group_name])).click()
    button_was_clicked = button_click_by_xpath(DatasetModal.sinks_selector_button(), context.driver,
                                               "Unable to click on on sink selector button dataset modal")
    assert_that(button_was_clicked, equal_to(True), "Unable to click on sink selector button on dataset modal")
    sinks_options = get_selector_options(context.driver)
    button_was_clicked = button_click_by_xpath(DatasetModal.sinks_selector_button(), context.driver,
                                               "Unable to click on on sink selector button dataset modal (2)")
    assert_that(button_was_clicked, equal_to(True), "Unable to click on sink selector button on dataset modal (2)")
    for sink in sinks_names:
        button_was_clicked = button_click_by_xpath(DatasetModal.sinks_selector_button(), context.driver,
                                                   "Unable to click on on sink selector button dataset modal (3)")
        assert_that(button_was_clicked, equal_to(True), "Unable to click on sink selector button on dataset modal (3)")
        ActionChains(context.driver).move_to_element(sinks_options[sink]).perform()
        WebDriverWait(context.driver, time_webdriver_wait).until(
            EC.element_to_be_clickable(sinks_options[sink])).click()
    dataset_name = dataset_name_prefix + random_string()
    input_text_by_xpath(DatasetModal.dataset_name(), dataset_name, context.driver)
    button_was_clicked = button_click_by_xpath(UtilButton.save_button(), context.driver,
                                               "Unable to click on save dataset modal")
    assert_that(button_was_clicked, equal_to(True), "Unable to click on save dataset modal")
    WebDriverWait(context.driver, 3).until(
        EC.text_to_be_present_in_element((By.CSS_SELECTOR, "span.title"), 'Dataset successfully created'),
        message="Confirmation span of dataset creation not displayed")
    all_datasets = list_datasets(context.token)
    context.dataset = None
    for dataset in all_datasets:
        if dataset['name'] == dataset_name:
            context.dataset = dataset
            break
    assert_that(context.dataset, is_not(None), "Unable to find dataset on orb backend")


@step("the dataset is removed")
def remove_dataset(context):
    policy_id = context.policy['id']
    policy_view_page = f"{orb_url}/pages/datasets/policies/view/{policy_id}"
    context.driver, current_url = go_to_page(policy_view_page, driver=context.driver)
    button_was_clicked = button_click_by_xpath(PolicyPage.dataset_details_modal(context.dataset['name']),
                                               context.driver, "Unable to click on dataset details modal")
    assert_that(button_was_clicked, equal_to(True), "Unable to click on dataset details modal")
    button_was_clicked = button_click_by_xpath(PolicyPage.remove_dataset_button(),
                                               context.driver, "Unable to click on remove dataset button")
    assert_that(button_was_clicked, equal_to(True), "Unable to click on remove dataset button")
    input_text_by_xpath(Dataset.dataset_name(), context.dataset['name'], context.driver)
    button_was_clicked = button_click_by_xpath(Dataset.dataset_remove_confirmation(),
                                               context.driver, "Unable to click on remove dataset confirmation")
    assert_that(button_was_clicked, equal_to(True), "Unable to click on remove dataset confirmation")
    WebDriverWait(context.driver, 3).until(
        EC.text_to_be_present_in_element((By.CSS_SELECTOR, "span.title"), 'Dataset successfully deleted'),
        message="Confirmation span of dataset remove not displayed")


@step("the dataset is edited and one more sink is inserted and name is changed")
def edit_dataset_UI_using_multiple_sinks(context):
    policy_id = context.policy['id']
    policy_view_page = f"{orb_url}/pages/datasets/policies/view/{policy_id}"
    context.driver, current_url = go_to_page(policy_view_page, driver=context.driver)
    sinks_names = list()
    for sink_id in context.existent_sinks_id:
        sink = get_sink(context.token, sink_id)
        sinks_names.append(sink['name'])
    button_was_clicked = button_click_by_xpath(PolicyPage.dataset_details_modal(context.dataset['name']),
                                               context.driver, "Unable to open dataset modal")
    assert_that(button_was_clicked, equal_to(True), "Unable to open dataset modal")
    button_was_clicked = button_click_by_xpath(DatasetModal.sinks_selector_button(),
                                               context.driver, "Unable to click on sink selector button")
    assert_that(button_was_clicked, equal_to(True), "Unable to click on sink selector button")
    sinks_options = get_selector_options(context.driver)
    button_was_clicked = button_click_by_xpath(DatasetModal.sinks_selector_button(),
                                               context.driver, "Unable to click on sink selector button (2)")
    assert_that(button_was_clicked, equal_to(True), "Unable to click on sink selector button (2)")
    selected_sinks = \
        WebDriverWait(context.driver,
                      time_webdriver_wait).until(EC.presence_of_all_elements_located((By.XPATH,
                                                                                      Dataset.selected_sinks())))
    selected_sinks_names = list()
    for sink in selected_sinks:
        selected_sinks_names.append(sink.text)
    sinks_names = list(set(sinks_names) - set(selected_sinks_names))
    for sink in sinks_names:
        button_was_clicked = button_click_by_xpath(DatasetModal.sinks_selector_button(),
                                                   context.driver, "Unable to click on sink selector button (3)")
        assert_that(button_was_clicked, equal_to(True), "Unable to click on sink selector button (3)")
        ActionChains(context.driver).move_to_element(sinks_options[sink]).perform()
        WebDriverWait(context.driver, time_webdriver_wait).until(
            EC.element_to_be_clickable(sinks_options[sink])).click()
    WebDriverWait(context.driver, time_webdriver_wait).until(
        EC.presence_of_element_located((By.XPATH, DatasetModal.dataset_name()))).clear()
    dataset_name = dataset_name_prefix + random_string(10)
    input_text_by_xpath(DatasetModal.dataset_name(), dataset_name, context.driver)
    button_was_clicked = button_click_by_xpath(UtilButton.save_button(),
                                               context.driver, "Unable to click on save dataset button")
    assert_that(button_was_clicked, equal_to(True), "Unable to click on save dataset button")
    WebDriverWait(context.driver, 3).until(
        EC.text_to_be_present_in_element((By.CSS_SELECTOR, "span.title"), 'Dataset successfully updated'),
        message="Confirmation span of dataset creation not displayed")
    all_datasets = list_datasets(context.token)
    context.dataset = None
    for dataset in all_datasets:
        if dataset['name'] == dataset_name:
            context.dataset = dataset
            break
    assert_that(context.dataset, is_not(None), "Unable to find dataset on orb backend")


@step("{amount_of_sinks} sinks are linked to the dataset and the new name is displayed")
def check_amount_sinks_ui(context, amount_of_sinks):
    button_was_clicked = button_click_by_xpath(PolicyPage.dataset_details_modal(context.dataset['name']),
                                               context.driver, "Unable to click on dataset modal")
    assert_that(button_was_clicked, equal_to(True), "Unable to click on dataset modal")
    selected_sinks = \
        WebDriverWait(context.driver,
                      time_webdriver_wait).until(EC.presence_of_all_elements_located((By.XPATH,
                                                                                      Dataset.selected_sinks())))
    assert_that(len(selected_sinks), equal_to(int(amount_of_sinks)), "Unexpected amount of sinks linked to dataset")
