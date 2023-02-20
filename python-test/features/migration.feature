@migration @AUTORETRY
Feature: Migration tests

@pre-migration
Scenario: Agent legacy + sink legacy -> sink OTEL
    Given the Orb user has a registered account
        And the Orb user logs in
        And that an agent with 1 orb tag(s) and OTEL disabled already exists and is online within 30 seconds
        And pktvisor state is running
        And referred agent is subscribed to 1 group
        And this agent's heartbeat shows that 1 groups are matching the agent
        And that a/an legacy sink already exists (migration)
        And 2 simple policies are applied to the group
        And this agent's heartbeat shows that 2 policies are applied and all has status running
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
        And referred sink must have active state on response within 120 seconds
        And 2 dataset(s) have validity valid and 0 have validity invalid in 30 seconds
    When the sink is updated and OTEL is enabled
        And referred sink must have active state on response after 10 seconds


@pre-migration
Scenario: Agent legacy + sink OTEL
    Given the Orb user has a registered account
        And the Orb user logs in
        And that an agent with 1 orb tag(s) and OTEL disabled already exists and is online within 30 seconds
        And pktvisor state is running
        And referred agent is subscribed to 1 group
        And this agent's heartbeat shows that 1 groups are matching the agent
        And that a/an OTEL sink already exists (migration)
        And 2 simple policies are applied to the group
        And this agent's heartbeat shows that 2 policies are applied and all has status running
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
        And referred sink must have active state on response within 120 seconds
        And 2 dataset(s) have validity valid and 0 have validity invalid in 30 seconds


@pre-migration
Scenario: Adding policies to an Agent legacy after migrate sink legacy to sink OTEL
    Given the Orb user has a registered account
        And the Orb user logs in
        And that an agent with 1 orb tag(s) and OTEL disabled already exists and is online within 30 seconds
        And pktvisor state is running
        And that a/an legacy sink already exists (migration)
        And the sink is updated and OTEL is enabled
        And referred agent is subscribed to 1 group
        And this agent's heartbeat shows that 1 groups are matching the agent
    When 2 simple policies are applied to the group
    Then this agent's heartbeat shows that 2 policies are applied and all has status running
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
        And referred sink must have active state on response within 120 seconds
        And 2 dataset(s) have validity valid and 0 have validity invalid in 30 seconds
        And referred sink must have active state on response within 120 seconds


@pre-migration
Scenario: Adding policies to an Agent legacy after migrate sink legacy to sink OTEL and agent legacy to otel
    Given the Orb user has a registered account
        And the Orb user logs in
        And that an agent with 1 orb tag(s) and OTEL disabled already exists and is online within 30 seconds
        And pktvisor state is running
        And that a/an legacy sink already exists (migration)
        And referred agent is subscribed to 1 group
        And this agent's heartbeat shows that 1 groups are matching the agent
        And 2 simple policies are applied to the group
        And this agent's heartbeat shows that 2 policies are applied and all has status running
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
        And referred sink must have active state on response within 120 seconds
        And 2 dataset(s) have validity valid and 0 have validity invalid in 30 seconds
        And referred sink must have active state on response within 120 seconds
    When the sink is updated and OTEL is enabled
        And stop the orb-agent container
        And the agent container is started on an available port and use otel:enabled env vars
        And the agent status is online
        And this agent's heartbeat shows that 1 groups are matching the agent
        And this agent's heartbeat shows that 2 policies are applied and all has status running
        And 2 dataset(s) have validity valid and 0 have validity invalid in 30 seconds
        And referred sink must have active state on response after 120 seconds
        And 2 simple policies are applied to the group
    Then this agent's heartbeat shows that 4 policies are applied and all has status running
        And 4 dataset(s) have validity valid and 0 have validity invalid in 30 seconds
        And referred sink must have active state on response within 120 seconds


@pos-migration
Scenario: Check if all sinks are OTEL after migration
    Given the Orb user has a registered account
        And the Orb user logs in
    Then all existing sinks must have OTEL enabled