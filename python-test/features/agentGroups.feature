@agentGroups
Feature: agent groups creation     
    
    Scenario: Create Agent Group
        Given that the user is logged in
            And that an agent already exists and is online
        When an Agent Group is created with same tag as the agent
        Then one agent must be matching on response field matching_agents
            And the container logs should contain the message "completed RPC subscription to group"
