@agentGroups
Feature: agent groups creation

    @smoke
    Scenario: Create Agent Group  with one tag
        Given the Orb user has a registered account
            And the Orb user logs in
            And that an agent with 1 orb tag(s) already exists and is online
        When an Agent Group is created with all tags contained in the agent
        Then 1 agent must be matching on response field matching_agents
            And the container logs should contain the message "completed RPC subscription to group" within 10 seconds

    @sanity
    Scenario: Create Agent Group with multiple tags
        Given the Orb user has a registered account
            And the Orb user logs in
            And that an agent with 5 orb tag(s) already exists and is online
        When an Agent Group is created with all tags contained in the agent
        Then 1 agent must be matching on response field matching_agents
            And the container logs should contain the message "completed RPC subscription to group" within 10 seconds

    @sanity
    Scenario: Create Agent Group without tags
        Given the Orb user has a registered account
            And the Orb user logs in
        When an Agent Group is created with 0 orb tag(s)
        Then Agent Group creation response must be an error with message 'malformed entity specification'

    @smoke
    Scenario: Create Agent Group without description
        Given the Orb user has a registered account
            And the Orb user logs in
            And that an agent with 1 orb tag(s) already exists and is online
        When an Agent Group is created with same tag as the agent and without description
        Then 1 agent must be matching on response field matching_agents
            And the container logs should contain the message "completed RPC subscription to group" within 10 seconds

    @smoke
    Scenario: Edit Agent Group name
        Given the Orb user has a registered account
            And the Orb user logs in
            And that an agent with 1 orb tag(s) already exists and is online
            And an Agent Group is created with all tags contained in the agent
        When the name of Agent Group is edited using: name=group_name
        Then 1 agent must be matching on response field matching_agents
            And the container logs should contain the message "completed RPC subscription to group" within 10 seconds

    @smoke
    Scenario: Agent Group name editing without name
        Given the Orb user has a registered account
            And the Orb user logs in
            And that an agent with 1 orb tag(s) already exists and is online
        And an Agent Group is created with all tags contained in the agent
        When the name of Agent Group is edited using: name=None
        Then agent group editing must fail
            And 1 agent must be matching on response field matching_agents
            And the container logs should contain the message "completed RPC subscription to group" within 10 seconds

    @smoke
    Scenario: Edit Agent Group description (without description)
        Given the Orb user has a registered account
            And the Orb user logs in
            And that an agent with 1 orb tag(s) already exists and is online
            And an Agent Group is created with all tags contained in the agent
        When the description of Agent Group is edited using: description=None
        Then 1 agent must be matching on response field matching_agents
            And the container logs should contain the message "completed RPC subscription to group" within 10 seconds

    @smoke
    Scenario: Edit Agent Group description (with description)
        Given the Orb user has a registered account
            And the Orb user logs in
            And that an agent with 1 orb tag(s) already exists and is online
            And an Agent Group is created with all tags contained in the agent
        When the description of Agent Group is edited using: description="Agent group test description"
        Then 1 agent must be matching on response field matching_agents
            And the container logs should contain the message "completed RPC subscription to group" within 10 seconds

    @smoke
    Scenario: Edit Agent Group tags (with tags - unsubscription)
        Given the Orb user has a registered account
            And the Orb user logs in
            And that an agent with 1 orb tag(s) already exists and is online
            And an Agent Group is created with all tags contained in the agent
        When the tags of Agent Group is edited using: tags=2 orb tag(s)
        Then 0 agent must be matching on response field matching_agents
            And the container logs should contain the message "completed RPC unsubscription to group" within 10 seconds

    @smoke
    Scenario: Edit Agent Group tags (with tags - subscription)
        Given the Orb user has a registered account
            And the Orb user logs in
            And that an agent with 1 orb tag(s) already exists and is online
            And an Agent Group is created with 1 orb tag(s)
        When the tags of Agent Group is edited using: tags=matching the agent
        Then 1 agent must be matching on response field matching_agents
            And the container logs contain the message "completed RPC subscription to group" referred to each matching group within 10 seconds

    @smoke
    Scenario: Edit Agent Group tags (without tags)
        Given the Orb user has a registered account
            And the Orb user logs in
            And that an agent with 1 orb tag(s) already exists and is online
            And an Agent Group is created with all tags contained in the agent
        When the tags of Agent Group is edited using: tags=None
        Then agent group editing must fail
            And 1 agent must be matching on response field matching_agents
            And the agent status in Orb should be online

    @smoke
    Scenario: Edit Agent Group name, description and tags
        Given the Orb user has a registered account
            And the Orb user logs in
            And that an agent with 1 orb tag(s) already exists and is online
            And an Agent Group is created with all tags contained in the agent
        When the name, tags, description of Agent Group is edited using: name=new_name/ tags=2 orb tag(s)/ description=None
        Then the container logs should contain the message "completed RPC unsubscription to group" within 10 seconds
            And 0 agent must be matching on response field matching_agents
            And the agent status in Orb should be online