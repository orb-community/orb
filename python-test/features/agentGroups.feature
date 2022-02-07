@agentGroups
Feature: agent groups creation     
    
    Scenario: Create Agent Group
        Given the Orb user has a registered account
            And the Orb user logs in
            And that an agent with random orb tag(s): 1 tag(s) already exists and is online
        When an Agent Group is created with same tag as the agent
        Then one agent must be matching on response field matching_agents
            And the container logs should contain the message "completed RPC subscription to group" within 10 seconds
