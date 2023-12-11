@integration_config_files_otel_backend @AUTORETRY
Feature: Integration tests using agent with otlp as backend provided via config file


@smoke_otel_backend
Scenario: provisioning agent with otel backend and applying 1 policy (config file - auto_provision=true)
    Given the Orb user has a registered account
        And the Orb user logs in
        And that a sink with default configuration type already exists
        And 1 Agent Group(s) is created with 1 orb tag(s) (lower case)
        And 1 policies with otel backend and yaml format are applied to the group
        And a new agent is created with 0 orb tag(s)
    When an agent(backend_type:otel, settings: {}) is provisioned via a configuration file on port available with matching 1 group agent tags and has status online. [Overwrite default: False. Paste only file: True]
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
        And a new agent is created with 0 orb tag(s)
    When an agent(backend_type:otel, settings: {}) is provisioned via a configuration file on port available with matching 1 group agent tags and has status online. [Overwrite default: False. Paste only file: True]
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
Scenario: removing policy from agent with otel backend (config file - auto_provision=true)
    Given the Orb user has a registered account
        And the Orb user logs in
        And that a sink with default configuration type already exists
        And 1 Agent Group(s) is created with 1 orb tag(s) (lower case)
        And 2 policies with otel backend and yaml format are applied to the group
        And a new agent is created with 0 orb tag(s)
    When an agent(backend_type:otel, settings: {}) is provisioned via a configuration file on port available with matching 1 group agent tags and has status online. [Overwrite default: False. Paste only file: True]
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
Scenario: Remove group to which agent with otel backend is linked
    Given the Orb user has a registered account
        And the Orb user logs in
        And a new agent is created with 0 orb tag(s)
        And an agent(backend_type:otel, settings: {}) is provisioned via a configuration file on port available with 3 agent tags and has status online. [Overwrite default: False. Paste only file: True]
        And otel state is waiting
        And referred agent is subscribed to 1 group
        And otel state is waiting
        And this agent's heartbeat shows that 1 groups are matching the agent
        And that a sink with default configuration type already exists
        And 2 policies with otel backend and yaml format are applied to the group
        And this agent's heartbeat shows that 2 policies are applied and all has status running
    When 1 group(s) to which the agent is linked is removed
    Then the container logs should contain the message "completed RPC unsubscription to group" within 30 seconds
        And the container logs contain the message "policy no longer used by any group, removing" referred to each policy within 30 seconds
        And this agent's heartbeat shows that 0 policies are applied to the agent
        And this agent's heartbeat shows that 0 groups are matching the agent
        And no dataset should be linked to the removed group anymore
        And 0 dataset(s) have validity valid and 2 have validity invalid in 30 seconds


@smoke_otel_backend
Scenario: Remotely restart agent with otel backend and policies applied
    Given the Orb user has a registered account
        And the Orb user logs in
        And 1 Agent Group(s) is created with 1 orb tag(s) (lower case)
        And that a sink with default configuration type already exists
        And 1 policies with otel backend and yaml format are applied to the group
        And a new agent is created with 0 orb tag(s)
        And an agent(backend_type:otel, settings: {}) is provisioned via a configuration file on port available with matching 1 group agent tags and has status online. [Overwrite default: False. Paste only file: True]
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


@smoke_otel_backend @config_file @auto_provision
Scenario: agent otel with only agent tags subscription to a group with policies created after provision the agent (config file - auto_provision=true)
    Given the Orb user has a registered account
        And the Orb user logs in
        And that a sink with default configuration type already exists
    When an agent(backend_type:otel, settings: {}) is self-provisioned via a configuration file on port available with 3 agent tags and has status online. [Overwrite default: False. Paste only file: True]
        And 1 Agent Group(s) is created with all tags contained in the agent
        And 2 policies with otel backend and yaml format are applied to the group
        And otel state is running
    Then 2 dataset(s) have validity valid and 0 have validity invalid in 30 seconds
        And the container logs should contain the message "completed RPC subscription to group" within 30 seconds
        And this agent's heartbeat shows that 1 groups are matching the agent
        And this agent's heartbeat shows that 2 policies are applied and all has status running
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
        And the container logs that were output after all policies have been applied contain the message "scraped and published telemetry" referred to each applied policy within 180 seconds
        And referred sink must have active state on response within 120 seconds
        And remove the agent .yaml generated on each scenario


@smoke_otel_backend @config_file @auto_provision
Scenario: agent otel with only agent tags subscription to a group with policies created before provision the agent (config file - auto_provision=true)
    Given the Orb user has a registered account
        And the Orb user logs in
        And that a sink with default configuration type already exists
        And 1 Agent Group(s) is created with 1 orb tag(s) (lower case)
        And 3 simple policies otel are applied to the group
        And a new agent is created with 0 orb tag(s)
    When an agent(backend_type:otel, settings: {}) is self-provisioned via a configuration file on port available with matching 1 group agent tags and has status online. [Overwrite default: False. Paste only file: False]
        And otel state is running
    Then 3 dataset(s) have validity valid and 0 have validity invalid in 30 seconds
        And this agent's heartbeat shows that 1 groups are matching the agent
        And the container logs should contain the message "completed RPC subscription to group" within 30 seconds
        And this agent's heartbeat shows that 3 policies are applied and all has status running
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
        And the container logs that were output after all policies have been applied contain the message "scraped and published telemetry" referred to each applied policy within 180 seconds
        And referred sink must have active state on response within 120 seconds
        And remove the agent .yaml generated on each scenario


@smoke_otel_backend @config_file @auto_provision
Scenario: agent otel with mixed tags subscription to a group with policies created after provision the agent (config file - auto_provision=true)
    Given the Orb user has a registered account
        And the Orb user logs in
        And that a sink with default configuration type already exists
    When an agent(backend_type:otel, settings: {}) is self-provisioned via a configuration file on port available with 3 agent tags and has status online. [Overwrite default: False. Paste only file: True]
        And otel state is waiting
        And edit the orb tags on agent and use 2 orb tag(s)
        And 1 Agent Group(s) is created with all tags contained in the agent
        And otel state is waiting
        And 3 policies with otel backend and yaml format are applied to the group
        And otel state is running
    Then 3 dataset(s) have validity valid and 0 have validity invalid in 30 seconds
        And this agent's heartbeat shows that 1 groups are matching the agent
        And the container logs should contain the message "completed RPC subscription to group" within 30 seconds
        And this agent's heartbeat shows that 3 policies are applied and all has status running
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
        And the container logs that were output after all policies have been applied contain the message "scraped and published telemetry" referred to each applied policy within 180 seconds
        And referred sink must have active state on response within 120 seconds
        And remove the agent .yaml generated on each scenario


@smoke_otel_backend @config_file @auto_provision
Scenario: agent otel with mixed tags subscription to a group with policies created before provision the agent (config file - auto_provision=true)
    Given the Orb user has a registered account
        And the Orb user logs in
        And that a sink with default configuration type already exists
        And 1 Agent Group(s) is created with 2 orb tag(s) (lower case)
        And 3 simple policies otel are applied to the group
        And a new agent is created with 2 orb tag(s)
    When an agent(backend_type:otel, settings: {}) is self-provisioned via a configuration file on port available with matching 1 group agent tags and has status online. [Overwrite default: False. Paste only file: False]
        And otel state is running
    Then 3 dataset(s) have validity valid and 0 have validity invalid in 30 seconds
        And this agent's heartbeat shows that 1 groups are matching the agent
        And the container logs should contain the message "completed RPC subscription to group" within 30 seconds
        And this agent's heartbeat shows that 3 policies are applied and all has status running
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
        And the container logs that were output after all policies have been applied contain the message "scraped and published telemetry" referred to each applied policy within 180 seconds
        And referred sink must have active state on response within 120 seconds
        And remove the agent .yaml generated on each scenario

@smoke_otel_backend @config_file
Scenario: agent otel with only agent tags subscription to a group with policies created after provision the agent (config file - auto_provision=false)
    Given the Orb user has a registered account
        And the Orb user logs in
        And that a sink with default configuration type already exists
        And a new agent is created with 0 orb tag(s)
    When an agent(backend_type:otel, settings: {}) is provisioned via a configuration file on port available with 3 agent tags and has status online. [Overwrite default: False. Paste only file: True]
        And otel state is waiting
        And 1 Agent Group(s) is created with all tags contained in the agent
        And 3 policies with otel backend and yaml format are applied to the group
        And otel state is running
    Then 3 dataset(s) have validity valid and 0 have validity invalid in 30 seconds
        And the container logs should contain the message "completed RPC subscription to group" within 30 seconds
        And this agent's heartbeat shows that 1 groups are matching the agent
        And this agent's heartbeat shows that 3 policies are applied and all has status running
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
        And the container logs that were output after all policies have been applied contain the message "scraped and published telemetry" referred to each applied policy within 180 seconds
        And referred sink must have active state on response within 120 seconds
        And remove the agent .yaml generated on each scenario

#@smoke_otel_backend @config_file
@MUTE
Scenario: agent otel with only agent tags subscription to a group with policies created before provision the agent (config file - auto_provision=false)
    Given the Orb user has a registered account
        And the Orb user logs in
        And that a sink with default configuration type already exists
        And 1 Agent Group(s) is created with 1 orb tag(s) (lower case)
        And 3 simple policies otel are applied to the group
        And a new agent is created with 0 orb tag(s)
    When an agent(backend_type:otel, settings: {}) is provisioned via a configuration file on port available with matching 1 group agent tags and has status online. [Overwrite default: False. Paste only file: False]
        And otel state is waiting
    Then 3 dataset(s) have validity valid and 0 have validity invalid in 30 seconds
        And this agent's heartbeat shows that 1 groups are matching the agent
        And the container logs should contain the message "completed RPC subscription to group" within 30 seconds
        And this agent's heartbeat shows that 3 policies are applied and all has status running
        And otel state is running
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
        And the container logs that were output after all policies have been applied contain the message "scraped and published telemetry" referred to each applied policy within 180 seconds
        And referred sink must have active state on response within 120 seconds
        And remove the agent .yaml generated on each scenario

@smoke_otel_backend @config_file
Scenario: agent otel with mixed tags subscription to a group with policies created after provision the agent (config file - auto_provision=false)
    Given the Orb user has a registered account
        And the Orb user logs in
        And that a sink with default configuration type already exists
        And a new agent is created with 2 orb tag(s)
    When an agent(backend_type:otel, settings: {}) is provisioned via a configuration file on port available with 3 agent tags and has status online. [Overwrite default: False. Paste only file: False]
        And edit the orb tags on agent and use 2 orb tag(s)
        And 1 Agent Group(s) is created with all tags contained in the agent
        And 3 policies with otel backend and yaml format are applied to the group
        And otel state is running
    Then 3 dataset(s) have validity valid and 0 have validity invalid in 30 seconds
        And this agent's heartbeat shows that 1 groups are matching the agent
        And the container logs should contain the message "completed RPC subscription to group" within 30 seconds
        And this agent's heartbeat shows that 3 policies are applied and all has status running
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
        And the container logs that were output after all policies have been applied contain the message "scraped and published telemetry" referred to each applied policy within 180 seconds
        And referred sink must have active state on response within 120 seconds
        And remove the agent .yaml generated on each scenario

#@smoke_otel_backend @config_file
@MUTE
Scenario: agent otel with mixed tags subscription to a group with policies created before provision the agent (config file - auto_provision=false)
    Given the Orb user has a registered account
        And the Orb user logs in
        And that a sink with default configuration type already exists
        And 1 Agent Group(s) is created with 2 orb tag(s) (lower case)
        And 3 simple policies otel are applied to the group
        And a new agent is created with 2 orb tag(s)
    When an agent(backend_type:otel, settings: {}) is provisioned via a configuration file on port available with matching 1 group agent tags and has status online. [Overwrite default: False. Paste only file: False]
        And otel state is running
    Then 3 dataset(s) have validity valid and 0 have validity invalid in 30 seconds
        And this agent's heartbeat shows that 1 groups are matching the agent
        And the container logs should contain the message "completed RPC subscription to group" within 30 seconds
        And this agent's heartbeat shows that 3 policies are applied and all has status running
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
        And the container logs that were output after all policies have been applied contain the message "scraped and published telemetry" referred to each applied policy within 180 seconds
        And referred sink must have active state on response within 120 seconds
        And remove the agent .yaml generated on each scenario


########### provisioning agents without specify otel configs on backend


@smoke_otel_backend @config_file @otel_configs @auto_provision
Scenario: provisioning agent without specify path to otel config file (config file - auto_provision=true)
    Given the Orb user has a registered account
        And the Orb user logs in
        And that a sink with default configuration type already exists
        And 1 Agent Group(s) is created with 1 orb tag(s) (lower case)
        And 3 policies with otel backend and yaml format are applied to the group
        And a new agent is created with 0 orb tag(s)
    When an agent(backend_type:otel, settings: {}) is self-provisioned via a configuration file on port available with matching 1 group agent tags and has status online. [Overwrite default: True. Paste only file: True. Use specif backend config {"config_file":"None"}]
        And otel state is running
    Then 3 dataset(s) have validity valid and 0 have validity invalid in 30 seconds
        And this agent's heartbeat shows that 1 groups are matching the agent
        And the container logs should contain the message "completed RPC subscription to group" within 30 seconds
        And this agent's heartbeat shows that 3 policies are applied and all has status running
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
        And the container logs that were output after all policies have been applied contain the message "scraped and published telemetry" referred to each applied policy within 180 seconds
        And remove the agent .yaml generated on each scenario


@smoke_otel_backend @config_file @otel_configs
Scenario: provisioning agent without specify path to otel config file (config file - auto_provision=false)
    Given the Orb user has a registered account
        And the Orb user logs in
        And that a sink with default configuration type already exists
        And a new agent is created with 2 orb tag(s)
    When an agent(backend_type:otel, settings: {}) is provisioned via a configuration file on port available with 3 agent tags and has status online. [Overwrite default: True. Paste only file: True. Use specif backend config {"config_file":"None"}]
        And edit the orb tags on agent and use 2 orb tag(s)
        And 1 Agent Group(s) is created with all tags contained in the agent
        And 3 policies with otel backend and yaml format are applied to the group
        And otel state is running
    Then 3 dataset(s) have validity valid and 0 have validity invalid in 30 seconds
        And this agent's heartbeat shows that 1 groups are matching the agent
        And the container logs should contain the message "completed RPC subscription to group" within 30 seconds
        And this agent's heartbeat shows that 3 policies are applied and all has status running
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
        And the container logs that were output after all policies have been applied contain the message "scraped and published telemetry" referred to each applied policy within 180 seconds
        And remove the agent .yaml generated on each scenario
