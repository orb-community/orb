from behave.model_core import Status
import os


def before_scenario(context, scenario):
    context.containers_id = dict()
    context.agent_groups = dict()
    context.existent_sinks_id = list()
    context.tap_tags = dict()
    # user = os.getuid()
    # if 'root' in scenario.tags and user != 0:
    #     scenario.skip('Root privileges are required')


def after_scenario(context, scenario):
    if 'access_denied' in context and context.access_denied is True:
        scenario.set_status(Status.skipped)
    # if scenario.status != Status.failed:
    #     context.execute_steps('''
    #     Then stop the orb-agent container
    #     Then remove the orb-agent container
    #     ''')
    # if "driver" in context:
    #     context.driver.close()
    #     context.driver.quit()
    # if "mocked_interface" in scenario.tags and scenario.status == Status.passed:
    #     context.execute_steps('''
    #     Then remove virtual switch
    #     Then remove dummy interface
    #     ''')
