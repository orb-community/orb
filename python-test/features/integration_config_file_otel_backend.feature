@integration_config_files_otel_backend
Feature: Integration tests using agent provided via config file and otel backend

@smoke_otel_backend
Scenario: provisioning agent with otel backend and applying 1 policy (config file - auto_provision=true)
    Given the Orb user has a registered account
        And the Orb user logs in
        And that a sink with default configuration type already exists
        And 1 Agent Group(s) is created with 1 orb tag(s) (lower case)
        And 1 policies with otel backend and yaml format are applied to the group
    When an agent with otel backend is self-provisioned via a configuration file on port available with matching 1 group agent tags and has status online. [Overwrite default: True. Paste only file: True]
        And otel state is running
    Then 1 dataset(s) have validity valid and 0 have validity invalid in 30 seconds
        And this agent's heartbeat shows that 1 groups are matching the agent
        And the container logs should contain the message "completed RPC subscription to group" within 30 seconds
        And this agent's heartbeat shows that 1 policies are applied and all has status running
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And the container logs that were output after all policies have been applied contain the message "scraped and published telemetry" referred to each applied policy within 180 seconds
        And referred sink must have active state on response within 180 seconds
        And remove the agent .yaml generated on each scenario


@smoke_otel_backend
Scenario: provisioning agent with otel backend and applying 2 policies (config file - auto_provision=true)
    Given the Orb user has a registered account
        And the Orb user logs in
        And that a sink with default configuration type already exists
        And 1 Agent Group(s) is created with 1 orb tag(s) (lower case)
        And 2 policies with otel backend and yaml format are applied to the group
    When an agent with otel backend is self-provisioned via a configuration file on port available with matching 1 group agent tags and has status online. [Overwrite default: True. Paste only file: True]
        And otel state is running
    Then 2 dataset(s) have validity valid and 0 have validity invalid in 30 seconds
        And this agent's heartbeat shows that 1 groups are matching the agent
        And the container logs should contain the message "completed RPC subscription to group" within 30 seconds
        And this agent's heartbeat shows that 2 policies are applied and all has status running
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And the container logs that were output after all policies have been applied contain the message "scraped and published telemetry" referred to each applied policy within 180 seconds
        And referred sink must have active state on response within 180 seconds
        And remove the agent .yaml generated on each scenario


@smoke_otel_backend
Scenario: provisioning agent with otel backend and removing 1 policy (config file - auto_provision=true)
    Given the Orb user has a registered account
        And the Orb user logs in
        And that a sink with default configuration type already exists
        And 1 Agent Group(s) is created with 1 orb tag(s) (lower case)
        And 2 policies with otel backend and yaml format are applied to the group
    When an agent with otel backend is self-provisioned via a configuration file on port available with matching 1 group agent tags and has status online. [Overwrite default: True. Paste only file: True]
        And otel state is running
    Then 2 dataset(s) have validity valid and 0 have validity invalid in 30 seconds
    And this agent's heartbeat shows that 2 policies are applied and all has status running
    When one of applied policies is removed
    Then referred policy must not be listed on the orb policies list
        And no dataset should be linked to the removed policy anymore
        And 1 dataset(s) have validity valid and 0 have validity invalid in 30 seconds
        And this agent's heartbeat shows that 1 policies are applied and all has status running
        And container logs should inform that removed policy was stopped and removed within 30 seconds
        And the container logs that were output after the policy have been removed contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
        And the container logs that were output after the policy have been removed does not contain the message "scraped and published telemetry" referred to deleted policy anymore


@smoke_otel_backend
Scenario: Remotely restart agent with otel backend and policies applied
    Given the Orb user has a registered account
        And the Orb user logs in
        And 1 Agent Group(s) is created with 1 orb tag(s) (lower case)
        And that a sink with default configuration type already exists
        And 1 policies with otel backend and yaml format are applied to the group
        And an agent with otel backend is self-provisioned via a configuration file on port available with matching 1 group agent tags and has status online. [Overwrite default: True. Paste only file: True]
        And this agent's heartbeat shows that 1 groups are matching the agent
        And otel state is running
        And this agent's heartbeat shows that 1 policies are applied and all has status running
    When remotely restart the agent
    Then the container logs that were output after reset the agent contain the message "otel process stopped" within 30 seconds
        And the container logs should contain the message "all backends and comms were restarted" within 30 seconds
        And the container logs that were output after reset the agent contain the message "removing policies" within 30 seconds
        And the container logs that were output after reset the agent contain the message "resetting backend" within 30 seconds
        And the container logs that were output after reset the agent contain the message "all backends and comms were restarted" within 30 seconds
        And the container logs that were output after reset the agent contain the message "policy applied successfully" referred to each applied policy within 30 seconds
        And the container logs that were output after reset the agent contain the message "scraped and published telemetry" referred to each applied policy within 180 seconds
