@agentGroups
Feature: agent groups creation     
    
    Scenario: Create Agent Group  with one tag
        Given the Orb user has a registered account
            And the Orb user logs in
            And that an agent with 1 orb tag(s) already exists and is online
        When an Agent Group is created with same tag as the agent
        Then one agent must be matching on response field matching_agents
            And the container logs should contain the message "completed RPC subscription to group" within 10 seconds


    Scenario: Create Agent Group with multiple tags
        Given the Orb user has a registered account
        And the Orb user logs in
        And that an agent with random orb tag(s): 5 tag(s) already exists and is online
        When an Agent Group is created with same tag as the agent
        Then one agent must be matching on response field matching_agents
        And the container logs should contain the message "completed RPC subscription to group" within 10 seconds


    Scenario: Create Agent Group without tags
        Given the Orb user has a registered account
        And the Orb user logs in
        When an Agent Group is created with random orb tag(s): 0 tag(s)
        Then Agent Group creation response must be an error with message 'malformed entity specification'