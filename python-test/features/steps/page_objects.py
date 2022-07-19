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
        return "//div[contains(@class, 'datatable-body')]"

    @classmethod
    def page_count(cls):
        return "//div[contains(@class, 'page-count')]"

    @classmethod
    def sub_pages(cls):
        return "//ul[@class='pager']/child::li[contains(@class, 'pages')]"

    @classmethod
    def agent(cls, name):
        return f"//div[contains(@class, 'agent-name') and contains(text(), '{name}')]"

    @classmethod
    def next_page(cls):
        return "//a[@aria-label='go to next page']"
