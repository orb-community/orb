@agents_ui @AUTORETRY
Feature: Create agents using orb ui

    @smoke_ui
    Scenario: Create agent
        Given that the Orb user logs in Orb UI
            And the user clicks on Agents on left menu
        When a new agent is created through the UI with 3 orb tag(s)
        Then the agents list and the agents view should display agent's status as New within 30 seconds

    @smoke_ui
    Scenario: Provision agent
        Given that the Orb user logs in Orb UI
            And the user clicks on Agents on left menu
        When a new agent is created through the UI with 2 orb tag(s)
            And the agent container is started using the command provided by the UI on an available port
        Then the agents list and the agents view should display agent's status as Online within 30 seconds
            And the agent status in Orb should be online within 30 seconds
            And the container logs should contain the message "sending capabilities" within 30 seconds

    @smoke_ui
    Scenario: Run two orb agents on the same port
        Given that the Orb user logs in Orb UI
            And that the user is on the orb Agents page
            And a new agent is created through the UI with 1 orb tag(s)
            And the agent container is started using the command provided by the UI on an available port
            And that the user is on the orb Agents page
        When a new agent is created through the UI with 1 orb tag(s)
            And the agent container is started using the command provided by the UI on an unavailable port
        Then last container created is running after 5 seconds
            And the container logs should contain the message "agent startup error" within 5 seconds
            And the container logs should contain "[error] unable to bind to localhost:port" as log within 5 seconds
            And first container created is running after 5 seconds

    @smoke_ui
    Scenario: Run two orb agents on the same port without restart always parameter
        Given that the Orb user logs in Orb UI
            And that the user is on the orb Agents page
            And a new agent is created through the UI with 1 orb tag(s)
            And the agent container is started using the command provided by the UI on an available port
            And that the user is on the orb Agents page
        When a new agent is created through the UI with 1 orb tag(s)
            And the agent container is started using the command provided by the UI without --restart=always on an unavailable port
        Then last container created is exited after 5 seconds
            And the container logs should contain the message "agent startup error" within 5 seconds
            And the container logs should contain "[error] unable to bind to localhost:port" as log within 5 seconds
            And first container created is running after 5 seconds

    @smoke_ui
    Scenario: Run two orb agents on different ports
        Given that the Orb user logs in Orb UI
            And that the user is on the orb Agents page
            And a new agent is created through the UI with 1 orb tag(s)
            And the agent container is started using the command provided by the UI on an available port
            And that the user is on the orb Agents page
        When a new agent is created through the UI with 1 orb tag(s)
            And the agent container is started using the command provided by the UI on an available port
        Then last container created is running after 5 seconds
            And first container created is running after 5 seconds
