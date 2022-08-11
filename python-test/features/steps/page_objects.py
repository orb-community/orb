# XPATHs

class LeftMenu:
    def __init__(self):
        pass

    @classmethod
    def agents(cls):
        return "//a[contains(@title, 'Agents')]"

    @classmethod
    def agent_group(cls):
        return "//a[contains(@title, 'Agent Groups')]"

    @classmethod
    def policies(cls):
        return "//a[contains(@title, 'Policy Management')]"

    @classmethod
    def sinks(cls):
        return "//a[contains(@title, 'Sink Management')]"


class OrbPagesUrl:
    def __init__(self):
        pass

    @classmethod
    def Agents(cls, orb_url):
        return f"{orb_url}/pages/fleet/agents"

    @classmethod
    def AgentGroups(cls, orb_url):
        return f"{orb_url}/pages/fleet/groups"

    @classmethod
    def Policies(cls, orb_url):
        return f"{orb_url}/pages/datasets/policies"

    @classmethod
    def Sinks(cls, orb_url):
        return f"{orb_url}/pages/sinks"

    @classmethod
    def PolicyAdd(cls, orb_url):
        return f"{orb_url}/pages/datasets/policies/add"


class AgentsPage:
    def __init__(self):
        pass

    @classmethod
    def new_agent_button(cls):
        return "//button[contains(text(), 'New Agent')]"

    @classmethod
    def agent_name(cls):
        return "//input[contains(@data-orb-qa-id, 'input#name')]"

    @classmethod
    def agent_tag_key(cls):
        return "//input[contains(@data-orb-qa-id, 'input#orb_tag_key')]"

    @classmethod
    def agent_tag_value(cls):
        return "//input[contains(@data-orb-qa-id, 'input#orb_tag_value')]"

    @classmethod
    def agent_add_tag_button(cls):
        return "//button[contains(@data-orb-qa-id, 'button#addTag')]"

    @classmethod
    def agent_key(cls):
        return "//label[contains(text(), 'Agent Key')]/following::pre[1]"

    @classmethod
    def agent_provisioning_command(cls):
        return "//label[contains(text(), 'Provisioning Command')]/following::pre[1]"

    @classmethod
    def agent_view_id(cls):
        return "//label[contains(text(), 'Agent ID')]/following::p"

    @classmethod
    def agent_status(cls):
        return "//span[@class='float-right']//child::span"


class UtilButton:
    def __init__(self):
        pass

    @classmethod
    def next_button(cls):
        return "//button[contains(text(), 'Next')]"

    @classmethod
    def save_button(cls):
        return "//button[contains(text(), 'Save')]"

    @classmethod
    def close_button(cls):
        return "//span[contains(@class, 'nb-close')]"

    @classmethod
    def add_button(cls):
        return "//button[contains(@class, 'appearance-ghost size-medium status-primary')]"

    @classmethod
    def selector_options(cls):
        return "//ul[@class='option-list']//nb-option"


class DataTable:
    def __init__(self):
        pass

    @classmethod
    def body(cls):
        return "//div[contains(@class, 'datatable-body')]"

    @classmethod
    def page_count(cls):
        return "//div[contains(@class, 'page-count')]"

    @classmethod
    def sub_pages(cls):
        return "//ul[@class='pager']/child::li[contains(@class, 'pages')]"

    @classmethod
    def agent(cls, name):
        return f"//span[contains(@class, 'agent-name') and contains(text(), '{name}')]"

    @classmethod
    def agent_status(cls, name):
        return f"//span[contains(text(), '{name}')]/ancestor::div[contains(@class, " \
               f"'datatable-row-center')]/descendant::i[contains(@class, " \
               f"'fa fa-circle')]/ancestor::span[contains(@class, 'ng-star-inserted')]"

    @classmethod
    def next_page(cls):
        return "//a[@aria-label='go to next page']"

    @classmethod
    def previous_page(cls):
        return "//a[@aria-label='go to previous page']"

    @classmethod
    def last_page(cls):
        return "//a[@aria-label='go to last page']"

    @classmethod
    def first_page(cls):
        return "//a[@aria-label='go to first page']"

    @classmethod
    def destroyed_on_click_button(cls):
        return "//nb-toast[contains(@class, 'destroy-by-click has-icon custom-icon')]"

    @classmethod
    def agent_group(cls, name):
        return f"//span[@class='ng-star-inserted' and contains(text(), '{name}')]"


class AgentGroupPage:
    def __init__(self):
        pass

    @classmethod
    def new_agent_group_button(cls):
        return "//button[contains(text( ), 'New Agent Group')]"

    @classmethod
    def agent_group_name(cls):
        return "//input[contains(@formcontrolname, 'name')]"

    @classmethod
    def agent_group_description(cls):
        return "//input[contains(@formcontrolname, 'description')]"

    @classmethod
    def agent_group_tag_key(cls):
        return "//input[contains(@data-orb-qa-id, 'input#orb_tag_key')]"

    @classmethod
    def agent_group_tag_value(cls):
        return "//input[contains(@data-orb-qa-id, 'input#orb_tag_value')]"

    @classmethod
    def agent_group_add_tag_button(cls):
        return "//button[contains(@data-orb-qa-id, 'button#addTag')]"


class PolicyPage:
    def __init__(self):
        pass

    @classmethod
    def policy_page_header(cls):
        return "//h4[text()='Create Agent Policy']"

    @classmethod
    def new_policy_button(cls):
        return "//button[contains(text(), 'New Policy')]"

    @classmethod
    def policy_name(cls):
        return "//input[@data-orb-qa-id='name']"

    @classmethod
    def policy_description(cls):
        return "//input[@data-orb-qa-id='description']"

    @classmethod
    def tap_selector_button(cls):
        return "//nb-select[@data-orb-qa-id='taps']//button[@class='select-button placeholder']"

    @classmethod
    def advanced_options_expander(cls):
        return "//nb-accordion-item-header[contains(@class, 'policy-advanced-options')]"

    @classmethod
    def host_spec(cls):
        return "//input[contains(@placeholder, '10.0.1.0')]"

    @classmethod
    def filter_expression(cls):
        return "//input[contains(@placeholder, 'port 53')]"

    @classmethod
    def handler_selector_button(cls):
        return "//nb-select[@id='selected_handler']//button[@class='select-button placeholder']"

    @classmethod
    def only_qname_suffix(cls):
        return "//input[contains(@placeholder, '.com')]"

    @classmethod
    def only_rcode_selector_button(cls):
        return "//label[contains(text(), 'RCODE')]/following-sibling::nb-select"

    @classmethod
    def save_handler_button(cls):
        return "//button[@data-orb-qa-id='addHandler' and contains(text(), 'Save')]"

    @classmethod
    def add_handler_button(cls):
        return "//button[contains(@class, 'add-handler-button')]"

    @classmethod
    def exclude_noerror_checkbox(cls):
        return "//span[@class='custom-checkbox']"

    @classmethod
    def policy_configurations(cls):
        return "//div[contains(@class, 'monaco-scrollable-element editor-scrollable')]"
