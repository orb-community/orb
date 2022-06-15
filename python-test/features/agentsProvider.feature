@agents
Feature: agent provider

@smoke
Scenario: Provision agent
    Given the Orb user has a registered account
        And the Orb user logs in
    When a new agent is created with 1 orb tag(s)
        And the agent container is started on an available port
        And the agent status is online
    Then the agent status in Orb should be online within 10 seconds
        And the container logs should contain the message "sending capabilities" within 10 seconds

@smoke
Scenario: Run two orb agents on the same port
    Given the Orb user has a registered account
        And the Orb user logs in
        And that an agent with 1 orb tag(s) already exists and is online
    When a new agent is created with 1 orb tag(s)
        And the agent container is started on an unavailable port
        And the agent status is new
    Then last container created is exited after 5 seconds
        And the container logs should contain the message "agent startup error" within 2 seconds
        And first container created is running after 5 seconds

@smoke
Scenario: Run two orb agents on different ports
    Given the Orb user has a registered account
        And the Orb user logs in
        And that an agent with 1 orb tag(s) already exists and is online
    When a new agent is created with 1 orb tag(s)
        And the agent container is started on an available port
        And the agent status is online
    Then last container created is running after 2 seconds
        And first container created is running after 5 seconds


@smoke
Scenario: Provision agent without tags
    Given the Orb user has a registered account
        And the Orb user logs in
    When a new agent is created with 0 orb tag(s)
        And the agent container is started on an available port
        And the agent status is online
    Then the agent status in Orb should be online within 10 seconds
        And the container logs should contain the message "sending capabilities" within 10 seconds


@smoke
Scenario: Provision agent with multiple tags
    Given the Orb user has a registered account
        And the Orb user logs in
    When a new agent is created with 5 orb tag(s)
        And the agent container is started on an available port
        And the agent status is online
    Then the agent status in Orb should be online within 10 seconds
        And the container logs should contain the message "sending capabilities" within 10 seconds


@smoke
Scenario: Edit agent tag
    Given the Orb user has a registered account
        And the Orb user logs in
        And a new agent is created with 5 orb tag(s)
        And the agent container is started on an available port
        And the agent status is online
    When edit the orb tags on agent and use 3 orb tag(s)
    Then the container logs should contain the message "sending capabilities" within 10 seconds
        And agent must have 3 tags
        And the agent status in Orb should be online within 10 seconds


@smoke
Scenario: Save agent without tag
    Given the Orb user has a registered account
        And the Orb user logs in
        And a new agent is created with 5 orb tag(s)
        And the agent container is started on an available port
        And the agent status is online
    When edit the orb tags on agent and use 0 orb tag(s)
    Then the container logs should contain the message "sending capabilities" within 10 seconds
        And agent must have 0 tags
        And the agent status in Orb should be online within 10 seconds


@smoke
Scenario: Insert tags in agents created without tags
    Given the Orb user has a registered account
        And the Orb user logs in
        And a new agent is created with 0 orb tag(s)
        And the agent container is started on an available port
        And the agent status is online
    When edit the orb tags on agent and use 2 orb tag(s)
    Then the container logs should contain the message "sending capabilities" within 10 seconds
        And agent must have 2 tags
        And the agent status in Orb should be online within 10 seconds


@smoke
Scenario: Edit agent name
    Given the Orb user has a registered account
        And the Orb user logs in
        And a new agent is created with 1 orb tag(s)
        And the agent container is started on an available port
        And the agent status is online
    When edit the agent name
    Then the container logs should contain the message "sending capabilities" within 10 seconds
        And agent must have 1 tags
        And the agent status in Orb should be online within 10 seconds


@smoke
Scenario: Edit agent name and tags
    Given the Orb user has a registered account
        And the Orb user logs in
        And a new agent is created with 1 orb tag(s)
        And the agent container is started on an available port
        And the agent status is online
    When edit the agent name and edit orb tags on agent using 3 orb tag(s)
    Then the container logs should contain the message "sending capabilities" within 10 seconds
        And agent must have 3 tags
        And the agent status in Orb should be online within 10 seconds


@smoke
Scenario: Stop agent container
    Given the Orb user has a registered account
        And the Orb user logs in
        And a new agent is created with 1 orb tag(s)
        And the agent container is started on an available port
        And the agent status is online
    When stop the orb-agent container
    Then the container logs should contain the message "pktvisor stopping" within 10 seconds
        And the agent status in Orb should be offline within 10 seconds
        And the container logs should not contain any error message


@sanity
Scenario: Forced remove agent container
    Given the Orb user has a registered account
        And the Orb user logs in
        And a new agent is created with 1 orb tag(s)
        And the agent container is started on an available port
        And the agent status is online
    When forced remove the orb-agent container
    Then the agent status in Orb should be stale within 360 seconds


@smoke
Scenario: Remove agent
    Given the Orb user has a registered account
        And the Orb user logs in
        And a new agent is created with 1 orb tag(s)
        And the agent container is started on an available port
        And the agent status is online
    When this agent is removed
    Then the container logs should contain the message "ERROR mqtt log" within 120 seconds
        And last container created is running after 120 seconds
