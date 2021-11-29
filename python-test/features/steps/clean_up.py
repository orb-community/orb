from behave import then
from control_plane_agents import delete_agents, list_agents
from control_plane_agent_groups import delete_agent_groups, list_agent_groups


@then("all agents should be deleted from orb")
def clean_agents(context):
    context.execute_steps('''
    Given that the user is logged in
    ''')
    token = context.token
    agents_list = list_agents(token)
    delete_agents(token, agents_list)


@then("all agent groups should be deleted from orb")
def clean_agent_groups(context):
    context.execute_steps('''
    Given that the user is logged in
    ''')
    token = context.token
    agent_groups_list = list_agent_groups(token)
    delete_agent_groups(token, agent_groups_list)
