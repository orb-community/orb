@integration
Feature: Integration tests

Scenario: Apply two policies to an agent
    Given the Orb user has a registered account
        And the Orb user logs in
        And that an agent already exists and is online
        And referred agent is subscribed to a group
        And that a sink already exists
    When 2 policies are applied to the agent
    Then this agent's heartbeat shows that 2 policies are successfully applied
        And the container logs contain the message "policy applied successfully" referred to each policy within 10 seconds
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
        And referred sink must have active state on response within 10 seconds


Scenario: Remove policy from agent
    Given the Orb user has a registered account
        And the Orb user logs in
        And that an agent already exists and is online
        And referred agent is subscribed to a group
        And that a sink already exists
    When 2 policies are applied to the agent
        And this agent's heartbeat shows that 2 policies are successfully applied
        And the container logs contain the message "policy applied successfully" referred to each policy within 10 seconds
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
        And referred sink must have active state on response within 10 seconds
    When one of applied policies is removed
    Then referred policy must not be listed on the orb policies list
        And this agent's heartbeat shows that 1 policies are successfully applied
        And container logs should inform that removed policy was stopped and removed within 10 seconds
        And the container logs that were output after the policy have been removed contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
        And the container logs that were output after the policy have been removed does not contain the message "scraped metrics for policy" referred to deleted policy anymore
