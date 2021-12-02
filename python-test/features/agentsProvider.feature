@agents
Feature: agent provider
   
    Scenario: Provision agent
        Given that the user is logged in on orb account
        When a new agent is created
            And the agent container is started
        Then the agent status in Orb should be online
            And the container logs should contain the message "sending capabilities"
