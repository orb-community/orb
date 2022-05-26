def before_scenario(context, scenario):
    context.containers_id = dict()
    context.agent_groups = dict()
    context.existent_sinks_id = list()


def after_scenario(context, feature):
    context.execute_steps('''
    Then remove the container
    ''')
