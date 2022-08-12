#XPATHs

class LeftMenu:
    def __init__(self):
        pass

    @classmethod
    def fleet_management(cls):
        return "//a[contains(@title, 'Fleet Management')]"

    @classmethod
    def agents(cls):
        return "//a[contains(@title, 'Agents')]"
    
    @classmethod
    def agent_group_menu(cls):
        return "//a[contains(@title, 'Agent Groups')]"


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
    def delete_agent_confirmation_button(cls):
        return "//button[contains(text( ), 'I Understand, Delete This Agent Group')]"
    
    @classmethod
    def delete_agent_confirmation_title(cls):
        return "//nb-card-header[contains(text( ), ' Delete Agent Group Confirmation ')]"
    
    @classmethod
    def delete_agent_confirmation_field(cls):
        return "//input[contains(@class, 'input-full-width')]"
    