#@agents_ui
#Feature: Create agents using orb ui
#
#    Scenario: Create agent
#        Given the Orb user logs in through the UI
#            And that fleet Management is clickable on ORB Menu
#            And that Agents is clickable on ORB Menu
#        When a new agent is created through the UI
#        Then the agents list and the agents view should display agent's status as New within 10 seconds
#
#    Scenario: Provision agent
#        Given the Orb user logs in through the UI
#            And that fleet Management is clickable on ORB Menu
#            And that Agents is clickable on ORB Menu
#        When a new agent is created through the UI
#            And the agent container is started using the command provided by the UI
#        Then the agents list and the agents view should display agent's status as Online within 10 seconds
#            And the agent status in Orb should be online
#            And the container logs should contain the message "sending capabilities" within 10 seconds
