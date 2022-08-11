from behave import step
from control_plane_policies import parse_policy_params
from test_config import TestConfig
from ui_utils import *
from page_objects import *

configs = TestConfig.configs()
orb_url = configs.get('orb_url')


@step("a new policy is created through the UI with: {kwargs}")
def create_new_policy_through_UI(context, kwargs):
    params = parse_policy_params(kwargs)
    WebDriverWait(context.driver, 3).until(
        EC.element_to_be_clickable((By.XPATH, PolicyPage.new_policy_button()))).click()
    WebDriverWait(context.driver, 3).until(
        EC.presence_of_element_located((By.XPATH, PolicyPage.policy_page_header())))
    WebDriverWait(context.driver, 3).until(EC.url_to_be(OrbPagesUrl.PolicyAdd(orb_url)))
    input_text_by_xpath(PolicyPage.policy_name(), params["name"], context.driver)
    if params["description"] is not None:
        input_text_by_xpath(PolicyPage.policy_description(), params["description"], context.driver)
    WebDriverWait(context.driver, 3).until(EC.element_to_be_clickable((By.XPATH, UtilButton.next_button()))).click()
    WebDriverWait(context.driver, 3).until(
        EC.element_to_be_clickable((By.XPATH, PolicyPage.tap_selector_button()))).click()
    taps_options = get_selector_options(context.driver)
    chosen_tap = [val for key, val in taps_options.items() if params["tap"] in key]
    if len(chosen_tap) == 1:
        chosen_tap[0].click()
    else:  # todo improve logic for more than one
        raise "Invalid option for taps. More than one options was detected."
    if params["host_specification"] is not None:
        WebDriverWait(context.driver, 3).until(
            EC.element_to_be_clickable((By.XPATH, PolicyPage.advanced_options_expander()))).click()
        input_text_by_xpath(PolicyPage.host_spec(), params["host_specification"], context.driver)
    if params["bpf_filter_expression"] is not None:
        input_text_by_xpath(PolicyPage.filter_expression(), params["bpf_filter_expression"], context.driver)
    WebDriverWait(context.driver, 3).until(EC.element_to_be_clickable((By.XPATH, UtilButton.next_button()))).click()
    WebDriverWait(context.driver, 3).until(
        EC.element_to_be_clickable((By.XPATH, PolicyPage.add_handler_button()))).click()
    WebDriverWait(context.driver, 3).until(
        EC.element_to_be_clickable((By.XPATH, PolicyPage.handler_selector_button()))).click()
    handlers_options = get_selector_options(context.driver)
    chosen_handler = [val for key, val in handlers_options.items() if params["handler"] in key]
    if len(chosen_handler) == 1:
        chosen_handler[0].click()
    else:  # todo improve logic for more than one
        raise "Invalid option for handlers. More than one options was detected."
    if params["exclude_noerror"] is not None and params["exclude_noerror"].lower() == "true":
        WebDriverWait(context.driver, 3).until(EC.element_to_be_clickable((By.XPATH, PolicyPage.exclude_noerror_checkbox()))).click()
    if params["only_qname_suffix"] is not None:
        params["only_qname_suffix"] = str(params["only_qname_suffix"]).replace("[", "").replace("]", "").replace("'",
                                                                                                                 "")
        input_text_by_xpath(PolicyPage.only_qname_suffix(), params["only_qname_suffix"], context.driver)
    if params["only_rcode"] is not None:
        WebDriverWait(context.driver, 3).until(
            EC.element_to_be_clickable((By.XPATH, PolicyPage.only_rcode_selector_button()))).click()
        rcodes = get_selector_options(context.driver)
        chosen_rcode = [val for key, val in rcodes.items() if params["only_rcode"] in key]
        if len(chosen_rcode) == 1:
            chosen_rcode[0].click()
        else:  # todo improve logic for more than one
            raise "Invalid option for rcode. More than one options was detected."
    WebDriverWait(context.driver, 3).until(EC.presence_of_element_located((By.XPATH, PolicyPage.save_handler_button())))
    WebDriverWait(context.driver, 3).until(
        EC.element_to_be_clickable((By.XPATH, PolicyPage.save_handler_button()))).click()
    WebDriverWait(context.driver, 3).until(EC.element_to_be_clickable((By.XPATH, UtilButton.save_button()))).click()
    context.policy_name = params["name"]


# @step('created policy must have the chosen parameters')
# def check_json_policies_ui(context):
#     WebDriverWait(context.driver, 3).until(EC.presence_of_element_located((By.XPATH, PolicyPage.policy_configurations())))
