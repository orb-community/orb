@agents @AUTORETRY
Feature: agent provider

@smoke
Scenario: Provision agent
    Given the Orb user has a registered account
        And the Orb user logs in
    When a new agent is created with 1 orb tag(s)
        And the agent container is started on an available port
        And the agent status is online
    Then the agent status in Orb should be online within 30 seconds
        And the container logs should contain the message "sending capabilities" within 30 seconds

@smoke
Scenario: Run two orb agents on the same port
    Given the Orb user has a registered account
        And the Orb user logs in
        And that an agent with 1 orb tag(s) already exists and is online
    When a new agent is created with 1 orb tag(s)
        And the agent container is started on an unavailable port
        And the agent status is new
    Then last container created is exited after 5 seconds
        And the container logs should contain the message "agent startup error" within 5 seconds
        And first container created is running after 5 seconds

@smoke
Scenario: Run two orb agents on different ports
    Given the Orb user has a registered account
        And the Orb user logs in
        And that an agent with 1 orb tag(s) already exists and is online
    When a new agent is created with 1 orb tag(s)
        And the agent container is started on an available port
        And the agent status is online
    Then last container created is running after 5 seconds
        And first container created is running after 5 seconds


@smoke
Scenario: Provision agent without tags
    Given the Orb user has a registered account
        And the Orb user logs in
    When a new agent is created with 0 orb tag(s)
        And the agent container is started on an available port
        And the agent status is online
    Then the agent status in Orb should be online within 30 seconds
        And the container logs should contain the message "sending capabilities" within 30 seconds


@smoke
Scenario: Provision agent with multiple tags
    Given the Orb user has a registered account
        And the Orb user logs in
    When a new agent is created with 5 orb tag(s)
        And the agent container is started on an available port
        And the agent status is online
    Then the agent status in Orb should be online within 30 seconds
        And the container logs should contain the message "sending capabilities" within 30 seconds


@smoke
Scenario: Edit agent tag
    Given the Orb user has a registered account
        And the Orb user logs in
        And a new agent is created with 5 orb tag(s)
        And the agent container is started on an available port
        And the agent status is online
    When edit the orb tags on agent and use 3 orb tag(s)
    Then the container logs should contain the message "sending capabilities" within 30 seconds
        And agent must have 3 tags
        And the agent status in Orb should be online within 30 seconds


@smoke
Scenario: Save agent without tag
    Given the Orb user has a registered account
        And the Orb user logs in
        And a new agent is created with 5 orb tag(s)
        And the agent container is started on an available port
        And the agent status is online
    When edit the orb tags on agent and use 0 orb tag(s)
    Then the container logs should contain the message "sending capabilities" within 30 seconds
        And agent must have 0 tags
        And the agent status in Orb should be online within 30 seconds


@smoke
Scenario: Insert tags in agents created without tags
    Given the Orb user has a registered account
        And the Orb user logs in
        And a new agent is created with 0 orb tag(s)
        And the agent container is started on an available port
        And the agent status is online
    When edit the orb tags on agent and use 2 orb tag(s)
    Then the container logs should contain the message "sending capabilities" within 30 seconds
        And agent must have 2 tags
        And the agent status in Orb should be online within 30 seconds


@smoke
Scenario: Edit agent name
    Given the Orb user has a registered account
        And the Orb user logs in
        And a new agent is created with 1 orb tag(s)
        And the agent container is started on an available port
        And the agent status is online
    When edit the agent name
    Then the container logs should contain the message "sending capabilities" within 30 seconds
        And agent must have 1 tags
        And the agent status in Orb should be online within 30 seconds


@smoke
Scenario: Edit agent name and tags
    Given the Orb user has a registered account
        And the Orb user logs in
        And a new agent is created with 1 orb tag(s)
        And the agent container is started on an available port
        And the agent status is online
    When edit the agent name and edit orb tags on agent using 3 orb tag(s)
    Then the container logs should contain the message "sending capabilities" within 30 seconds
        And agent must have 3 tags
        And the agent status in Orb should be online within 30 seconds


@smoke
Scenario: Stop agent container
    Given the Orb user has a registered account
        And the Orb user logs in
        And a new agent is created with 1 orb tag(s)
        And the agent container is started on an available port
        And the agent status is online
    When stop the orb-agent container
    Then the container logs should contain the message "stop signal received stopping agent" within 30 seconds
        And the container logs should contain the message "pktvisor process stopped" within 30 seconds
        And the agent status in Orb should be offline within 30 seconds
        And the container logs should not contain any error message
        And the container logs should not contain any panic message


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
        And the container logs should contain the message "error reconnecting with MQTT, stopping agent" within 120 seconds
        And last container created is exited within 70 seconds
        And the container logs should not contain any panic message
        And last container created is exited after 120 seconds


@smoke @sanity
Scenario: Create agent with name conflict
    Given the Orb user has a registered account
        And the Orb user logs in
        And that an agent with 1 orb tag(s) already exists and is online
    When a new agent is requested to be created with the same name as an existent one
    Then the error message on response is failed to create agent

@smoke @sanity
Scenario: Edit agent using an already existent name (conflict)
    Given the Orb user has a registered account
        And the Orb user logs in
        And that an agent with 1 orb tag(s) already exists and is online
        And that an agent with 1 orb tag(s) already exists and is online
    When edit the agent name using an already existent one
    Then the error message on response is entity already exists

@smoke
Scenario: Run agent with dnstap, pcap, sflow, netflow env vars
    Given the Orb user has a registered account
        And the Orb user logs in
        And a new agent is created with 1 orb tag(s)
    When the agent container is started on an available port and use ALL env vars
        And the agent status is online
    Then created agent has taps: dnstap, pcap, sflow, netflow


@smoke
Scenario: Run agent with dnstap env vars
    Given the Orb user has a registered account
        And the Orb user logs in
        And a new agent is created with 1 orb tag(s)
    When the agent container is started on an available port and use dnstap env vars
        And the agent status is online
    Then created agent has taps: dnstap


@smoke
Scenario: Run agent with pcap env vars
    Given the Orb user has a registered account
        And the Orb user logs in
        And a new agent is created with 1 orb tag(s)
    When the agent container is started on an available port and use pcap env vars
        And the agent status is online
    Then created agent has taps: pcap


@smoke
Scenario: Run agent with sflow env vars
    Given the Orb user has a registered account
        And the Orb user logs in
        And a new agent is created with 1 orb tag(s)
    When the agent container is started on an available port and use sflow env vars
        And the agent status is online
    Then created agent has taps: sflow


@smoke
Scenario: Run agent with netflow env vars
    Given the Orb user has a registered account
        And the Orb user logs in
        And a new agent is created with 1 orb tag(s)
    When the agent container is started on an available port and use netflow env vars
        And the agent status is online
    Then created agent has taps: netflow
