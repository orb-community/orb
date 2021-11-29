@agentGroups
Feature: agent groups creation     
    
    Scenario: Create Agent Group
        Given that the user is logged in
            And that an agent already exists and be online
        When an Agent Group is created with same tag as agent
        Then one matching agent must be seen
            And the container logs should contain the message "completed RPC subscription to group"
