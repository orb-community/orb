from behave.model_core import Status


def before_scenario(context, scenario):
    context.containers_id = dict()
    context.agent_groups = dict()
    context.existent_sinks_id = list()


def after_scenario(context, scenario):
    if scenario.status != Status.failed:
        context.execute_steps('''
        Then stop the orb-agent container
        Then remove the orb-agent container
        ''')
    if "driver" in context:
        context.driver.close()
