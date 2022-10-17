from behave import step
from control_plane_policies import parse_policy_params, get_policy, make_policy_json, list_policies
from test_config import TestConfig
from ui_utils import *
from utils import remove_empty_from_json
from page_objects import *
from hamcrest import *
from deepdiff import DeepDiff
from selenium.common.exceptions import ElementNotInteractableException

configs = TestConfig.configs()
orb_url = configs.get('orb_url')


@step("a new policy is created through the UI with: {kwargs}")
def create_new_policy_through_UI(context, kwargs):
    button_was_clicked = button_click_by_xpath(PolicyPage.new_policy_button(),
                                               context.driver, "Unable to click on new policy button")
    assert_that(button_was_clicked, equal_to(True), "Unable to click on new policy button")
    WebDriverWait(context.driver, time_webdriver_wait).until(
        EC.presence_of_element_located((By.XPATH, PolicyPage.policy_page_header())))
    WebDriverWait(context.driver, time_webdriver_wait).until(EC.url_to_be(OrbPagesUrl.PolicyAdd(orb_url)))
    params = parse_policy_params(kwargs)
    threading.Event().wait(time_webdriver_wait)
    input_text_by_xpath(PolicyPage.policy_name(), params["name"], context.driver)
    if params["description"] is not None:
        input_text_by_xpath(PolicyPage.policy_description(), params["description"], context.driver)
        WebDriverWait(context.driver, time_webdriver_wait).until(
            (EC.text_to_be_present_in_element_value((By.XPATH, PolicyPage.policy_description()),
                                                    params["description"])))
        threading.Event().wait(time_webdriver_wait)
    WebDriverWait(context.driver, time_webdriver_wait).until(
        (EC.text_to_be_present_in_element_value((By.XPATH, PolicyPage.policy_name()),
                                                params["name"])))
    button_was_clicked = button_click_by_xpath(UtilButton.next_button(),
                                               context.driver, "Unable to click on next button")
    assert_that(button_was_clicked, equal_to(True), "Unable to click on next button")
    button_was_clicked = button_click_by_xpath(PolicyPage.tap_selector_button(),
                                               context.driver, "Unable to click on tap selector button")
    assert_that(button_was_clicked, equal_to(True), "Unable to click on tap selector button")
    taps_options = get_selector_options(context.driver)
    chosen_tap = [val for key, val in taps_options.items() if params["tap"] in key]
    if len(chosen_tap) == 1:
        chosen_tap[0].click()
    else:  # todo improve logic for more than one
        raise "Invalid option for taps. More than one options was detected."
    if params["host_specification"] is not None:
        button_was_clicked = button_click_by_xpath(PolicyPage.advanced_options_expander(),
                                                   context.driver, "Unable to click on advanced options expander"
                                                                   " button")
        assert_that(button_was_clicked, equal_to(True), "Unable to click on advanced options expander button")
        input_text_by_xpath(PolicyPage.host_spec(), params["host_specification"], context.driver)
    if params["bpf_filter_expression"] is not None:
        input_text_by_xpath(PolicyPage.filter_expression(), params["bpf_filter_expression"], context.driver)
    button_was_clicked = button_click_by_xpath(UtilButton.next_button(),
                                               context.driver, "Unable to click on next button")
    assert_that(button_was_clicked, equal_to(True), "Unable to click on next button")
    button_was_clicked = button_click_by_xpath(PolicyPage.add_handler_button(),
                                               context.driver, "Unable to click on add handler button")
    assert_that(button_was_clicked, equal_to(True), "Unable to click on add handler button")
    button_was_clicked = button_click_by_xpath(PolicyPage.handler_selector_button(),
                                               context.driver, "Unable to click on select handler button")
    assert_that(button_was_clicked, equal_to(True), "Unable to click on select handler button")
    handlers_options = get_selector_options(context.driver)
    chosen_handler = [val for key, val in handlers_options.items() if params["handler"] in key]
    if len(chosen_handler) == 1:
        chosen_handler[0].click()
        params["handle_label"] = WebDriverWait(context.driver, time_webdriver_wait).until(
            EC.presence_of_element_located((By.XPATH, PolicyPage.handler_name()))).get_attribute('value')
    else:  # todo improve logic for more than one
        raise "Invalid option for handlers. More than one options was detected."
    if params["exclude_noerror"] is not None and params["exclude_noerror"].lower() == "true":
        button_was_clicked = button_click_by_xpath(PolicyPage.exclude_noerror_checkbox(),
                                                   context.driver, "Unable to click on noerror checkbox")
        assert_that(button_was_clicked, equal_to(True), "Unable to click on noerror checkbox")
    if params["only_qname_suffix"] is not None:
        params["only_qname_suffix"] = str(params["only_qname_suffix"]).replace("[", "").replace("]", "").replace("'",
                                                                                                                 "")
        input_text_by_xpath(PolicyPage.only_qname_suffix(), params["only_qname_suffix"], context.driver)
    if params["only_rcode"] is not None:
        button_was_clicked = button_click_by_xpath(PolicyPage.only_rcode_selector_button(),
                                                   context.driver, "Unable to click on only rcode button")
        assert_that(button_was_clicked, equal_to(True), "Unable to click on only rcode button")
        rcodes = get_selector_options(context.driver)
        chosen_rcode = [val for key, val in rcodes.items() if params["only_rcode"] in key]
        if len(chosen_rcode) == 1:
            chosen_rcode[0].click()
        else:  # todo improve logic for more than one
            raise "Invalid option for rcode. More than one options was detected."
    WebDriverWait(context.driver, time_webdriver_wait).until(
        EC.presence_of_element_located((By.XPATH, PolicyPage.save_handler_button())))
    button_was_clicked = button_click_by_xpath(PolicyPage.save_handler_button(),
                                               context.driver, "Unable to click on save handler button")
    assert_that(button_was_clicked, equal_to(True), "Unable to click on save handler button")
    button_was_clicked = button_click_by_xpath(UtilButton.save_button(),
                                               context.driver, "Unable to click on save policy button")
    assert_that(button_was_clicked, equal_to(True), "Unable to click on save policy button")
    WebDriverWait(context.driver, time_webdriver_wait).until(
        EC.text_to_be_present_in_element((By.CSS_SELECTOR, "span.title"), 'Agent Policy successfully created'),
        message="Confirmation span of policy creation is not correctly displayed")
    context.policy_name = params["name"]
    context.policy_json = make_policy_json(params["name"], params['handle_label'],
                                           params["handler"], params["description"], params["tap"],
                                           params["input_type"], params["host_specification"],
                                           params["bpf_filter_expression"], params["pcap_source"],
                                           params["only_qname_suffix"], params["only_rcode"],
                                           params["exclude_noerror"], params["backend_type"])
    all_policies = list_policies(context.token)
    for policy in all_policies:
        if policy['name'] == params["name"]:
            context.policy = policy
            break


@step('user must be directed to the policy view page')
def check_policy_view_page(context):
    WebDriverWait(context.driver, time_webdriver_wait).until(
        EC.presence_of_element_located((By.XPATH, PolicyPage.policy_view_header())))
    policy_name = \
        WebDriverWait(context.driver, time_webdriver_wait).until(
            EC.presence_of_element_located((By.XPATH, PolicyPage.policy_view_name()))).text
    assert_that(policy_name, equal_to(context.policy_name), "Policy name is not correct in policy view page")
    current_url = str(context.driver.current_url)
    context.policy_id = current_url.split("/")[-1]


@step('created policy must have the chosen parameters')
def check_json_policies_ui(context):
    policy_back = get_policy(context.token, context.policy_id)
    policy_back = remove_empty_from_json(policy_back)
    policy_ui = remove_empty_from_json(context.policy_json)
    ddiff = DeepDiff(policy_back, policy_ui, view='tree', exclude_paths=["root['ts_created']", "root['id']",
                                                                         "root['schema_version']",
                                                                         "root['format']",
                                                                         "root['ts_last_modified']",
                                                                         "root['version']"])
    assert_that(ddiff, equal_to({}), f"{ddiff}")
    # todo validate editor text


@step('created policy {condition} displayed on policy pages')
def find_policy_in_policies_list(context, condition):
    policy_to_be_sought = PolicyPage.policy(context.policy_name)
    policy_on_datatable = find_element_on_datatable_by_condition(context.driver, policy_to_be_sought,
                                                                 LeftMenu.policies(), condition)
    if condition == "is":
        assert_that(policy_on_datatable, is_not(None), "Unable to find policy in policy datatable")
        policy_on_datatable.click()
        check_policy_view_page(context)
    else:
        get_policy(context.token, context.policy_id, 404)
        assert_that(policy_on_datatable, equal_to(None))


@step('remove policy from Orb UI')
def remove_policy_from_orb_ui(context):
    policy_removal_icon = PolicyPage.remove_policy_button(context.policy_name)
    removal_policy_button = \
        find_element_on_datatable_by_condition(context.driver, policy_removal_icon, LeftMenu.policies())
    try:
        removal_policy_button.click()
    except ElementNotInteractableException:
        context.driver.execute_script("arguments[0].click();", removal_policy_button)
    except Exception as err:
        raise ValueError(err)
    input_text_by_xpath(PolicyPage.remove_policy_confirmation_name(), context.policy_name, context.driver)
    button_click_by_xpath("//html", context.driver, "Unable to click on save policy button")  # blank space
    button_was_clicked = button_click_by_xpath(PolicyPage.remove_policy_confirmation_button(),
                                               context.driver, "Unable to click on remove policy confirmation button")
    assert_that(button_was_clicked, equal_to(True), "Unable to click on remove policy confirmation button")
    WebDriverWait(context.driver, time_webdriver_wait).until(
        EC.text_to_be_present_in_element((By.CSS_SELECTOR, "span.title"), 'Agent Policy successfully deleted'),
        message="Confirmation span of policy removal is not correctly displayed")
