@integration
Feature: Integration tests

Scenario: Apply two policies to an agent
    Given the Orb user has a registered account
        And the Orb user logs in
        And that an agent with random orb tag(s): 1 tag(s) already exists and is online
        And referred agent is subscribed to a group
        And that a sink already exists
    When 2 policies are applied to the group
    Then this agent's heartbeat shows that 2 policies are successfully applied
        And the container logs contain the message "policy applied successfully" referred to each policy within 10 seconds
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
        And referred sink must have active state on response within 10 seconds
        And datasets related to all existing policies have validity valid


Scenario: apply one policy using multiple datasets to the same group
    Given the Orb user has a registered account
        And the Orb user logs in
        And that an agent with random orb tag(s): 2 tag(s) already exists and is online
        And referred agent is subscribed to a group
        And that a sink already exists
    When 2 policies are applied to the group by 3 datasets each
    Then this agent's heartbeat shows that 2 policies are successfully applied
        And 3 datasets are linked with each policy on agent's heartbeat
        And the container logs contain the message "policy applied successfully" referred to each policy within 10 seconds
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
        And referred sink must have active state on response within 10 seconds
        And datasets related to all existing policies have validity valid


Scenario: Remove group to which agent is linked
    Given the Orb user has a registered account
        And the Orb user logs in
        And that an agent with random orb tag(s): 1 tag(s) already exists and is online
        And referred agent is subscribed to a group
        And that a sink already exists
        And 1 policies are applied to the group
        And this agent's heartbeat shows that 1 policies are successfully applied
    When the group to which the agent is linked is removed
    Then the container logs should contain the message "completed RPC unsubscription to group" within 10 seconds
        And dataset related have validity invalid


    Scenario: Remove policy from agent
    Given the Orb user has a registered account
        And the Orb user logs in
        And that an agent with random orb tag(s): 3 tag(s) already exists and is online
        And referred agent is subscribed to a group
        And that a sink already exists
        And 2 policies are applied to the group
        And this agent's heartbeat shows that 2 policies are successfully applied
    When one of applied policies is removed
    Then referred policy must not be listed on the orb policies list
        And datasets related to removed policy has validity invalid
        And datasets related to all existing policies have validity valid
        And this agent's heartbeat shows that 1 policies are successfully applied
        And container logs should inform that removed policy was stopped and removed within 10 seconds
        And the container logs that were output after the policy have been removed contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
        And the container logs that were output after the policy have been removed does not contain the message "scraped metrics for policy" referred to deleted policy anymore


Scenario: Remove dataset from agent with just one dataset linked
    Given the Orb user has a registered account
        And the Orb user logs in
        And that an agent with random orb tag(s): 2 tag(s) already exists and is online
        And referred agent is subscribed to a group
        And that a sink already exists
        And 1 policies are applied to the group
        And this agent's heartbeat shows that 1 policies are successfully applied
    When a dataset linked to this agent is removed
    Then referred dataset must not be listed on the orb datasets list
        And this agent's heartbeat shows that 0 policies are successfully applied
        And container logs should inform that removed policy was stopped and removed within 10 seconds
        And the container logs that were output after removing dataset contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
        And the container logs that were output after removing dataset does not contain the message "scraped metrics for policy" referred to deleted policy anymore


Scenario: Remove dataset from agent with more than one dataset linked
    Given the Orb user has a registered account
        And the Orb user logs in
        And that an agent with random orb tag(s): 4 tag(s) already exists and is online
        And referred agent is subscribed to a group
        And that a sink already exists
        And 3 policies are applied to the group
        And this agent's heartbeat shows that 3 policies are successfully applied
    When a dataset linked to this agent is removed
    Then referred dataset must not be listed on the orb datasets list
        And this agent's heartbeat shows that 2 policies are successfully applied
        And container logs should inform that removed policy was stopped and removed within 10 seconds
        And the container logs that were output after removing dataset contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
        And the container logs that were output after removing dataset does not contain the message "scraped metrics for policy" referred to deleted policy anymore


Scenario: Provision agent with tags matching an existent group
    Given the Orb user has a registered account
        And the Orb user logs in
        And an Agent Group is created with random orb tag(s): 2 tag(s)
    When a new agent is created with tags matching an existing group
        And the agent container is started on port default
    Then the agent status in Orb should be online
        And the container logs should contain the message "completed RPC subscription to group" within 10 seconds


Scenario: Provision agent with tag matching existing group linked to a valid dataset
    Given the Orb user has a registered account
        And the Orb user logs in
        And an Agent Group is created with random orb tag(s): 3 tag(s)
        And that a sink already exists
        And 2 policies are applied to the group
    When a new agent is created with tags matching an existing group
        And the agent container is started on port default
    Then this agent's heartbeat shows that 2 policies are successfully applied
        And the container logs contain the message "policy applied successfully" referred to each policy within 10 seconds
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
        And referred sink must have active state on response within 10 seconds
        And datasets related to all existing policies have validity valid


Scenario: Sink with invalid endpoint
    Given the Orb user has a registered account
        And the Orb user logs in
        And that an agent with random orb tag(s): 1 tag(s) already exists and is online
        And referred agent is subscribed to a group
        And that a sink with invalid endpoint already exists
        And that a policy already exists
    When a new dataset is created using referred group, sink and policy ID
    Then the container logs should contain the message "managing agent policy from core" within 10 seconds
        And the container logs should contain the message "policy applied successfully" within 10 seconds
        And the container logs should contain the message "scraped metrics for policy" within 180 seconds
        And referred sink must have error state on response within 10 seconds
        And dataset related have validity valid


Scenario: Sink with invalid username
    Given the Orb user has a registered account
        And the Orb user logs in
        And that an agent with random orb tag(s): 1 tag(s) already exists and is online
        And referred agent is subscribed to a group
        And that a sink with invalid username already exists
        And that a policy already exists
    When a new dataset is created using referred group, sink and policy ID
    Then the container logs should contain the message "managing agent policy from core" within 10 seconds
        And the container logs should contain the message "policy applied successfully" within 10 seconds
        And the container logs should contain the message "scraped metrics for policy" within 180 seconds
        And referred sink must have error state on response within 10 seconds
        And dataset related have validity valid


Scenario: Sink with invalid password
    Given the Orb user has a registered account
        And the Orb user logs in
        And that an agent with random orb tag(s): 1 tag(s) already exists and is online
        And referred agent is subscribed to a group
        And that a sink with invalid password already exists
        And that a policy already exists
    When a new dataset is created using referred group, sink and policy ID
    Then the container logs should contain the message "managing agent policy from core" within 10 seconds
        And the container logs should contain the message "policy applied successfully" within 10 seconds
        And the container logs should contain the message "scraped metrics for policy" within 180 seconds
        And referred sink must have error state on response within 10 seconds
        And dataset related have validity valid


    Scenario: Agent subscription to multiple groups created after provisioning agent
    Given the Orb user has a registered account
        And the Orb user logs in
        And a new agent is created with defined orb tag(s): region:br, demo:true, ns1:true
        And the agent container is started on port default
    When an Agent Group is created with defined orb tag(s): demo:true, ns1:true
        And an Agent Group is created with defined orb tag(s): region:br
        And an Agent Group is created with defined orb tag(s): demo:true
    Then the container logs contain the message "completed RPC subscription to group" referred to each group within 10 seconds


    Scenario: Agent subscription to multiple groups created before provisioning agent
        Given the Orb user has a registered account
            And the Orb user logs in
            And an Agent Group is created with defined orb tag(s): demo:true, ns1:true
            And an Agent Group is created with defined orb tag(s): region:br
            And an Agent Group is created with defined orb tag(s): demo:true
        When a new agent is created with defined orb tag(s): region:br, demo:true, ns1:true
            And the agent container is started on port default
        Then the container logs contain the message "completed RPC subscription to group" referred to each group within 10 seconds