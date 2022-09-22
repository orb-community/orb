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
    WebDriverWait(context.driver, time_webdriver_wait).until(
        EC.element_to_be_clickable((By.XPATH, PolicyPage.new_dataset_button()))).click()
    WebDriverWait(context.driver, time_webdriver_wait).until(
        EC.element_to_be_clickable((By.XPATH, DatasetModal.agent_group()))).click()
    groups_options = get_selector_options(context.driver)
    ActionChains(context.driver).move_to_element(groups_options[group_name]).perform()
    WebDriverWait(context.driver, time_webdriver_wait).until(
        EC.element_to_be_clickable(groups_options[group_name])).click()
    WebDriverWait(context.driver, time_webdriver_wait).until(
        EC.element_to_be_clickable((By.XPATH, DatasetModal.sinks_selector_button()))).click()
    sinks_options = get_selector_options(context.driver)
    WebDriverWait(context.driver, time_webdriver_wait).until(
        EC.element_to_be_clickable((By.XPATH, DatasetModal.sinks_selector_button()))).click()
    for sink in sinks_names:
        WebDriverWait(context.driver, time_webdriver_wait).until(
            EC.element_to_be_clickable((By.XPATH, DatasetModal.sinks_selector_button()))).click()
        ActionChains(context.driver).move_to_element(sinks_options[sink]).perform()
        WebDriverWait(context.driver, time_webdriver_wait).until(
            EC.element_to_be_clickable(sinks_options[sink])).click()
    dataset_name = dataset_name_prefix + random_string()
    input_text_by_xpath(DatasetModal.dataset_name(), dataset_name, context.driver)
    WebDriverWait(context.driver, time_webdriver_wait).until(
        EC.element_to_be_clickable((By.XPATH, UtilButton.save_button()))).click()
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
