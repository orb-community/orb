@agents
Feature: agent provider
   
    Scenario: Provision agent
        Given the Orb user has a registered account
            And the Orb user logs in
        When a new agent is created
            And the agent container is started
        Then the agent status in Orb should be online
            And the container logs should contain the message "sending capabilities" within 10 seconds
