from behave import given, then, step
from ui_utils import *
from control_plane_agents import agent_name_prefix, get_agent
from selenium.webdriver.support.ui import WebDriverWait
from selenium.webdriver.support import expected_conditions as EC
from selenium.webdriver.common.by import By
from selenium.common.exceptions import TimeoutException
from utils import random_string, create_tags_set, threading_wait_until
from test_config import TestConfig
from hamcrest import *
from page_objects import *

configs = TestConfig.configs()
orb_url = configs.get('orb_url')


@given("the user clicks on {element} on left menu")
def click_element_left_menu(context, element):
    dict_elements = {"Agents": LeftMenu.agents(), "Agent Groups": LeftMenu.agent_group(),
                     "Policy Management": LeftMenu.policies(), "Sink Management": LeftMenu.sinks()}
    WebDriverWait(context.driver, 3).until(
        EC.element_to_be_clickable((By.XPATH, dict_elements[element])), message=f"Unable to find {element} "
                                                                                f"icon on left menu")
    context.driver.find_element(By.XPATH, dict_elements[element]).click()


@step("that the user is on the orb {element} page")
def check_which_orb_page(context, element):
    dict_pages = {"Agents": OrbPagesUrl.Agents(orb_url), "Agent Groups": OrbPagesUrl.AgentGroups(orb_url),
                  "Policies": OrbPagesUrl.Policies(orb_url), "Sinks": OrbPagesUrl.Sinks(orb_url)}
    click_element_left_menu(context, element)
    WebDriverWait(context.driver, 5).until(EC.url_to_be(dict_pages[element]),
                                           message=f"Orb {element} page not available")
    current_url = context.driver.current_url
    assert_that(current_url, equal_to(dict_pages[element]), f"user not enabled to access orb {element} page")


@step("a new agent is created through the UI with {orb_tags} orb tag(s)")
def create_agent_through_the_agents_page(context, orb_tags):
    context.orb_tags = create_tags_set(orb_tags)
    WebDriverWait(context.driver, 3).until(
        EC.element_to_be_clickable((By.XPATH, AgentsPage.new_agent_button())), message="Unable to click on new agent"
                                                                                       " button").click()
    WebDriverWait(context.driver, 5).until(EC.url_to_be(f"{OrbPagesUrl.Agents(orb_url)}/add"),
                                           message="Orb add agents page not "
                                                   "available")
    context.agent_name = agent_name_prefix + random_string(10)
    input_text_by_xpath(AgentsPage.agent_name(), context.agent_name, context.driver)
    WebDriverWait(context.driver, 3).until(
        EC.element_to_be_clickable((By.XPATH, UtilButton.next_button())), message="Unable to click on next "
                                                                                  "button (page 1)").click()
    for tag_key, tag_value in context.orb_tags.items():
        input_text_by_xpath(AgentsPage.agent_tag_key(), tag_key, context.driver)
        input_text_by_xpath(AgentsPage.agent_tag_value(), tag_value, context.driver)
        WebDriverWait(context.driver, 3).until(
            EC.element_to_be_clickable((By.XPATH, AgentsPage.agent_add_tag_button())), message="Unable to click on add"
                                                                                               " tag button").click()
    WebDriverWait(context.driver, 3).until(
        EC.element_to_be_clickable((By.XPATH, UtilButton.next_button())), message="Unable to click on next "
                                                                                  "button (page 2)").click()
    WebDriverWait(context.driver, 3).until(
        EC.element_to_be_clickable((By.XPATH, UtilButton.save_button())), message="Unable to click on save "
                                                                                  "button").click()
    WebDriverWait(context.driver, 3).until(
        EC.text_to_be_present_in_element((By.CSS_SELECTOR, "span.title"), 'Agent successfully created'),
        message="Confirmation span of agent creation not displayed")
    context.agent_key = \
        WebDriverWait(context.driver, 3).until(EC.presence_of_element_located((By.XPATH, AgentsPage.agent_key())),
                                               message="Agent key not displayed").text
    agent_provisioning_command = \
        WebDriverWait(context.driver, 3).until(
            EC.presence_of_element_located((By.XPATH, AgentsPage.agent_provisioning_command())),
            message="Provisioning command not displayed").text

    context.agent_provisioning_command = agent_provisioning_command.replace("\n\n", " ")
    WebDriverWait(context.driver, 3).until(
        EC.element_to_be_clickable((By.XPATH, UtilButton.close_button())), message="Unable to click on close "
                                                                                   "button").click()
    agent = find_element_on_datatable_by_condition(context.driver, DataTable.agent(context.agent_name),
                                                   LeftMenu.agents())
    assert_that(agent, is_not(None), f"Unable to find the agent: {context.agent_name}")
    agent.click()
    context.agent = dict()
    context.agent['id'] = WebDriverWait(context.driver, 3).until(
        EC.presence_of_element_located((By.XPATH, AgentsPage.agent_view_id())), message="Agent id not displayed").text
    assert_that(context.agent['id'],
                matches_regexp(r'[a-zA-Z0-9]{8}\-[a-zA-Z0-9]{4}\-[a-zA-Z0-9]{4}\-[a-zA-Z0-9]{4}\-[a-zA-Z0-9]{12}'),
                f"Failed to get agent id {context.agent['id']}")
    context.agent['name'] = context.agent_name


@step("the agents list and the agents view should display agent's status as {status} within {time_to_wait} seconds")
def check_status_on_orb_ui(context, status, time_to_wait):
    agents_page = f"{orb_url}/pages/fleet/agents"
    context.driver, current_url = go_to_page(agents_page, driver=context.driver)
    agent_status_datatable = check_agent_status_on_orb_ui(context.driver, DataTable.agent_status(context.agent_name),
                                                          status, timeout=time_to_wait)
    assert_that(agent_status_datatable, is_not(None), f"Unable to find status of the agent: {context.agent_name}"
                                                      f" on datatable")
    assert_that(agent_status_datatable, equal_to(status), f"Agent {context.agent['id']} status failed on Agents list")
    WebDriverWait(context.driver, 3).until(
        EC.element_to_be_clickable((By.XPATH, DataTable.agent(context.agent_name))),
        message="Unable to click on agent name").click()
    agent_view_status = WebDriverWait(context.driver, 3).until(
        EC.presence_of_element_located(
            (By.XPATH, AgentsPage.agent_status())), message="Unable to find agent status on agent view page").text
    agent_view_status = agent_view_status.replace(".", "")
    assert_that(agent_view_status, equal_to(status), f"Agent {context.agent['id']} status failed")
    context.agent = get_agent(context.token, context.agent['id'])


@threading_wait_until
def check_agent_status_on_orb_ui(driver, agent_xpath, status, event=None):
    """

    :param driver: webdriver running
    :param (str) agent_xpath: xpath of the agent whose status need to be checked
    :param (str) status: agent expected status
    :param event: threading.event
    :return: web element refereed to the agent
    """
    agent_status_datatable = find_element_on_datatable_by_condition(driver, agent_xpath, LeftMenu.agents())
    if agent_status_datatable is not None and agent_status_datatable.text == status:
        event.set()
        return agent_status_datatable.text
    driver.refresh()
    return agent_status_datatable


@step("the policy must have status {status} on agent view page (Active Policies/Datasets)")
def get_policy_status_agent_view_page(context, status):
    agent_id = context.agent['id']
    agent_view_page = OrbUrl.agent_view(orb_url, agent_id)
    context.driver, current_url = go_to_page(agent_view_page, driver=context.driver)
    name = context.policy['name']
    WebDriverWait(context.driver, time_webdriver_wait).until(
        EC.url_to_be(agent_view_page))
    policy_status, policy_name = get_policy_ui_status(context.driver, name, status)
    assert_that(policy_name, equal_to(name), "Unexpected policy name")
    assert_that(policy_status, equal_to(status), "Unexpected policy status on agent view page")


@step("policy and dataset are clickable and redirect user to referred pages")
def click_policy_and_dataset(context):
    agent_id = context.agent['id']
    agent_view_page = OrbUrl.agent_view(orb_url, agent_id)
    policy_id = context.policy['id']
    policy_view_page = OrbUrl.policy_view(orb_url, policy_id)
    policy_name = context.policy['name']
    dataset_name = context.dataset['name']
    context.driver, current_url = go_to_page(agent_view_page, driver=context.driver)
    expand_policy_button = get_expand_policy_button(context.driver, policy_name)
    assert_that(expand_policy_button, is_not(None), "Unable to expand policy to check dataset on agent view page")
    WebDriverWait(context.driver, time_webdriver_wait).until(
        EC.element_to_be_clickable(expand_policy_button)).click()
    WebDriverWait(context.driver, time_webdriver_wait).until(
        EC.element_to_be_clickable((By.XPATH, AgentsPage.dataset_button(dataset_name)))).click()
    WebDriverWait(context.driver, time_webdriver_wait).until(
        EC.presence_of_element_located((By.XPATH, Dataset.DetailsModal())), "Dataset details modal was not opened "
                                                                            "after click above dataset name on agent "
                                                                            "views page")
    context.driver, current_url = go_to_page(agent_view_page, driver=context.driver)
    go_to_policy_button = get_policy_name_button(context.driver, policy_name)
    WebDriverWait(context.driver, time_webdriver_wait).until(
        EC.element_to_be_clickable(go_to_policy_button)).click()
    WebDriverWait(context.driver, time_webdriver_wait).until(
        EC.url_to_be(policy_view_page), "User not redirect to policy view page after click on policy name on agent view"
                                        " page")


@step("group must be listed on agent view page (Active Groups)")
def check_group_in_active_group(context):
    agent_id = context.agent['id']
    agent_view_page = OrbUrl.agent_view(orb_url, agent_id)
    context.driver, current_url = go_to_page(agent_view_page, driver=context.driver)
    group_name = list(context.agent_groups.values())[0]
    WebDriverWait(context.driver, time_webdriver_wait).until(
        EC.element_to_be_clickable((By.XPATH, AgentsPage.active_groups(group_name))),
        message="Unable to click on Active Group on agent view page").click()
    WebDriverWait(context.driver, time_webdriver_wait).until(
        EC.presence_of_element_located((By.XPATH, AgentGroupPage.MatchingGroupsModal())), "Matching groups modal was "
                                                                                          "not opened after click "
                                                                                          "above group name on agent "
                                                                                          "views page")


@step("policy must be removed from the agent")
def check_no_policies_agent(context):
    agent_id = context.agent['id']
    agent_view_page = OrbUrl.agent_view(orb_url, agent_id)
    context.driver, current_url = go_to_page(agent_view_page, driver=context.driver)
    found_span = wait_xpath(context.driver, AgentsPage.no_policies_span(), timeout=180)
    assert_that(found_span, equal_to(True), "Policies was note removed from the agent")


@threading_wait_until
def get_policy_ui_status(driver, name, status, event=None):
    elements = WebDriverWait(driver, time_webdriver_wait).until(
        EC.presence_of_all_elements_located((By.XPATH, AgentsPage.policies_and_datasets(name))))
    span_elements = list()
    policy_name = None
    for element in elements:
        if element.tag_name == 'span':
            span_elements.append(element)
        elif element.tag_name == 'button':
            policy_name = element
    policy_status = span_elements[2].text
    policy_name = policy_name.text
    if span_elements[0].text == "Policy:" and span_elements[1].text == "Status:" and policy_status == status:
        event.set()
    else:
        driver.refresh()
        event.wait(1)
    return policy_status, policy_name


@threading_wait_until
def get_expand_policy_button(driver, name, event=None):
    elements = WebDriverWait(driver, time_webdriver_wait).until(
        EC.presence_of_all_elements_located((By.XPATH, AgentsPage.policies_and_datasets(name))))
    expand_policy_icon = None
    for element in elements:
        if element.tag_name == 'nb-icon':
            expand_policy_icon = element
    if expand_policy_icon is not None:
        event.set()
    else:
        driver.refresh()
        event.wait(1)
    return expand_policy_icon


@threading_wait_until
def get_policy_name_button(driver, name, event=None):
    elements = WebDriverWait(driver, time_webdriver_wait).until(
        EC.presence_of_all_elements_located((By.XPATH, AgentsPage.policies_and_datasets(name))))
    policy_name = None
    for element in elements:
        if element.tag_name == 'button':
            policy_name = element
    if policy_name.text == name:
        event.set()
    else:
        driver.refresh()
        event.wait(1)
    return policy_name


@threading_wait_until
def wait_xpath(driver, xpath, event=None):
    try:
        WebDriverWait(driver, time_webdriver_wait).until(
            EC.presence_of_element_located((By.XPATH, xpath)))
        event.set()
        return event.is_set()
    except TimeoutException:
        driver.refresh()
        event.wait(1)
        return event.is_set()
