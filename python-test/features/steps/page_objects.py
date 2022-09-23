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


class Dataset:

    def __init__(self):
        pass

    @classmethod
    def DetailsModal(cls):
        return f"//nb-card-header[contains(text(), 'Dataset Details')]/ancestor::nb-dialog-container"


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

    @classmethod
    def policies_and_datasets(cls, name):
        return f"//nb-card-header[contains(text(), 'Active " \
               f"Policies/Datasets')]/ancestor::nb-card/descendant::nb-card-body/descendant::nb-accordion/descendant" \
               f"::button[contains(text(),'{name}')]/ancestor::nb-accordion-item-header//*"

    @classmethod
    def active_groups(cls, name):
        return f"//nb-card-header[contains(text(), 'Active " \
               f"Groups')]/ancestor::nb-card/descendant::nb-card-body/descendant::span[contains(text()," \
               f"'{name}')]/ancestor::button"

    @classmethod
    def dataset_button(cls, name):
        return f"//span[contains(text(),'Dataset:')]/following-sibling::button[contains(text(), '{name}')]"


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

    @classmethod
    def clear_all_filters(cls):
        return "//button[contains(@class, 'clear-filter clear-all')]"


class DataTable:
    def __init__(self):
        pass

    @classmethod
    def body(cls):
        return "//*[contains(@class, 'datatable-body')]"

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

    @classmethod
    def trash_icon(cls):
        return "//*[@data-name='trash-2']"

    @classmethod
    def filter_by(cls):
        return "//nb-select[@placeholder='Filter by']"

    @classmethod
    def option_list(cls):
        return "//ul[@class='option-list']"

    @classmethod
    def all_filter_options(cls):
        return "//nb-option[@class='nb-transition ng-star-inserted']"

    @classmethod
    def filter_by_name_field(cls):
        return "//input[@placeholder='Name']"

    @classmethod
    def plus_button(cls):
        return "//button[contains(@class, 'appearance-ghost size-medium status-primary')]"

    @classmethod
    def sink_name_on_datatable(cls, name):
        return f"//span[@class='ng-star-inserted' and contains(text(), '{name}')]"

    @classmethod
    def edit_icon(cls):
        return "//*[@data-name='edit']"

    @classmethod
    def close_option_selected(cls):
        return "//*[@class='fa fa-window-close']"


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

    @classmethod
    def delete_agent_group_confirmation_button(cls):
        return "//button[contains(text(), 'I Understand, Delete This Agent Group')]"

    @classmethod
    def delete_agent_group_confirmation_title(cls):
        return "//nb-card-header[contains(text(), 'Delete Agent Group Confirmation')]"

    @classmethod
    def delete_agent_group_confirmation_field(cls):
        return "//input[contains(@class, 'input-full-width')]"

    @classmethod
    def MatchingGroupsModal(cls):
        return f"//nb-card-header[contains(text(), ' Matching Agents')]/ancestor::nb-dialog-container"


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
        return "//input[@id='name']"

    @classmethod
    def policy_description(cls):
        return "//input[@id='description']"

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
    def handler_name(cls):
        return "//input[@data-orb-qa-id='handler-label']"

    @classmethod
    def exclude_noerror_checkbox(cls):
        return "//span[@class='custom-checkbox']"

    @classmethod
    def policy_configurations(cls):
        return "//div[contains(@class, 'monaco-scrollable-element editor-scrollable')]"

    @classmethod
    def policy_configurations_lines(cls):
        return "//div[@class='view-line']"

    @classmethod
    def policy_view_header(cls):
        return "//h4[text()='Policy View']"

    @classmethod
    def policy_view_name(cls):
        return "//label[@class='summary-accent' and text()='Policy Name']//following-sibling::p"

    @classmethod
    def policy(cls, policy_name):
        return f"//button[contains(@class, 'view-policy-button') and contains(text(),'{policy_name}')]"

    @classmethod
    def remove_policy_button(cls, policy_name):
        return f"//button[contains(@class, 'view-policy-button') and contains(text(), '{policy_name}')]/ancestor::datatable-body-row//child::button[contains(@class, 'orb-action-hover del-button')]"

    @classmethod
    def remove_policy_confirmation_name(cls):
        return "//input[@data-orb-qa-id='input#name']"

    @classmethod
    def remove_policy_confirmation_button(cls):
        return "//button[@data-orb-qa-id='button#delete']"

    @classmethod
    def new_dataset_button(cls):
        return "//button[contains(text( ), 'New Dataset')]"


class DatasetModal:
    def __init__(self):
        pass

    @classmethod
    def agent_group(cls):
        return "//input[@formcontrolname='agent_group_name']"

    @classmethod
    def sinks_selector_button(cls):
        return "//button[contains(@class, 'select-button') and contains(@class, 'placeholder')]"

    @classmethod
    def dataset_name(cls):
        return "//input[@data-orb-qa-id='name']"


class SinkPage:
    def __init__(self):
        pass

    @classmethod
    def new_sink_button(cls):
        return "//button[contains(text( ), 'New Sink')]"

    @classmethod
    def name_label(cls):
        return "//input[(@data-orb-qa-id= 'name')]"

    @classmethod
    def sink_description(cls):
        return "//input[(@data-orb-qa-id= 'description')]"

    @classmethod
    def remote_url(cls):
        return "//input[(@data-orb-qa-id= 'remote_host')]"

    @classmethod
    def sink_tag_key(cls):
        return "//input[(@data-orb-qa-id= 'input#orb_tag_key')]"

    @classmethod
    def sink_tag_value(cls):
        return "//input[(@id= 'value')]"

    @classmethod
    def sink_add_tag_button(cls):
        return "//button[contains(@data-orb-qa-id, 'button#addTag')]"

    @classmethod
    def save_button(cls):
        return "//button[contains(@data-orb-qa-id, 'button#save')]"

    @classmethod
    def sink_username(cls):
        return "//input[(@data-orb-qa-id= 'username')]"

    @classmethod
    def sink_password(cls):
        return "//input[(@data-orb-qa-id= 'password')]"

    @classmethod
    def sink_password(cls):
        return "//input[(@data-orb-qa-id= 'password')]"

    @classmethod
    def delete_sink_confirmation_field(cls):
        return "//input[contains(@class, 'input-full-width')]"

    @classmethod
    def delete_sink_confirmation_title(cls):
        return "//nb-card-header[contains(text(), 'Delete Sink Confirmation')]"

    @classmethod
    def delete_sink_confirmation_button(cls):
        return "//button[contains(text(), 'I Understand, Delete This Sink')]"

    @classmethod
    def next_button(cls):
        return "//*[@data-orb-qa-id= 'button#next']"
