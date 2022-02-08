@agents
Feature: agent provider
   
    Scenario: Provision agent
        Given the Orb user has a registered account
            And the Orb user logs in
        When a new agent is created with 1 orb tag(s)
            And the agent container is started on port default
        Then the agent status in Orb should be online
            And the container logs should contain the message "sending capabilities" within 10 seconds

    Scenario: Run two orb agents on the same port
        Given the Orb user has a registered account
            And the Orb user logs in
            And that an agent with 1 orb tag(s) already exists and is online
        When a new agent is created with 1 orb tag(s)
            And the agent container is started on port default
        Then last container created is exited after 2 seconds
            And the container logs should contain the message "agent startup error" within 2 seconds
            And container on port default is running after 2 seconds

    Scenario: Run two orb agents on different ports
        Given the Orb user has a registered account
            And the Orb user logs in
            And that an agent with 1 orb tag(s) already exists and is online
        When a new agent is created with 1 orb tag(s)
            And the agent container is started on port 10854
        Then last container created is running after 2 seconds
            And container on port default is running after 2 seconds