@integration
Feature: Integration tests

Scenario: Apply two policies to an agent
    Given the Orb user has a registered account
        And the Orb user logs in
        And that an agent with 1 orb tag(s) already exists and is online
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
        And that an agent with 2 orb tag(s) already exists and is online
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
        And that an agent with 1 orb tag(s) already exists and is online
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
        And that an agent with 3 orb tag(s) already exists and is online
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
        And that an agent with 3 orb tag(s) already exists and is online
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
        And that an agent with 4 orb tag(s) already exists and is online
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
        And an Agent Group is created with 2 orb tag(s)
    When a new agent is created with tags matching an existing group
        And the agent container is started on port default
    Then the agent status in Orb should be online
        And the container logs should contain the message "completed RPC subscription to group" within 10 seconds


Scenario: Provision agent with tag matching existing group linked to a valid dataset
    Given the Orb user has a registered account
        And the Orb user logs in
        And an Agent Group is created with 3 orb tag(s)
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
        And that an agent with 1 orb tag(s) already exists and is online
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
        And that an agent with 1 orb tag(s) already exists and is online
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
        And that an agent with 1 orb tag(s) already exists and is online
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
        And a new agent is created with region:br, demo:true, ns1:true orb tag(s)
        And the agent container is started on port default
    When an Agent Group is created with demo:true, ns1:true orb tag(s)
        And an Agent Group is created with region:br orb tag(s)
        And an Agent Group is created with demo:true orb tag(s)
    Then the container logs contain the message "completed RPC subscription to group" referred to each matching group within 10 seconds


Scenario: Agent subscription to multiple groups created before provisioning agent
    Given the Orb user has a registered account
        And the Orb user logs in
        And an Agent Group is created with demo:true, ns1:true orb tag(s)
        And an Agent Group is created with region:br orb tag(s)
        And an Agent Group is created with demo:true orb tag(s)
    When a new agent is created with demo:true, ns1:true orb tag(s)
        And the agent container is started on port default
    Then the container logs contain the message "completed RPC subscription to group" referred to each matching group within 10 seconds


Scenario: Agent subscription to group after editing agent's tags (agent provisioned before editing and group created after)
    Given the Orb user has a registered account
        And the Orb user logs in
        And a new agent is created with demo:true, ns1:true orb tag(s)
        And the agent container is started on port default
        And an Agent Group is created with demo:true, ns1:true orb tag(s)
    When edit the agent tags and use region:br orb tag(s)
        And an Agent Group is created with region:br orb tag(s)
    Then the container logs contain the message "completed RPC subscription to group" referred to each matching group within 10 seconds
        And this agent's heartbeat shows that 1 groups are matching the agent


Scenario: Agent subscription to group after editing agent's tags (editing tags after agent provision)
    Given the Orb user has a registered account
        And the Orb user logs in
        And an Agent Group is created with demo:true, ns1:true orb tag(s)
        And an Agent Group is created with region:br orb tag(s)
        And a new agent is created with demo:true, ns1:true orb tag(s)
        And the agent container is started on port default
    When edit the agent tags and use region:br orb tag(s)
    Then the container logs contain the message "completed RPC subscription to group" referred to each matching group within 10 seconds
        And this agent's heartbeat shows that 1 groups are matching the agent


Scenario: Agent subscription to group after editing agent's tags (editing tags before agent provision)
    Given the Orb user has a registered account
        And the Orb user logs in
        And a new agent is created with demo:true, ns1:true orb tag(s)
        And edit the agent tags and use region:br orb tag(s)
        And the agent container is started on port default
    When an Agent Group is created with region:br orb tag(s)
    Then the container logs contain the message "completed RPC subscription to group" referred to each matching group within 10 seconds
        And this agent's heartbeat shows that 1 groups are matching the agent


Scenario: Agent subscription to multiple group with policies after editing agent's tags (editing tags after agent provision)
    Given the Orb user has a registered account
        And the Orb user logs in
        And that an agent with test:true orb tag(s) already exists and is online
        And referred agent is subscribed to a group
        And that a sink already exists
        And 2 policies are applied to the group
        And this agent's heartbeat shows that 2 policies are successfully applied
        And an Agent Group is created with region:br orb tag(s)
        And 1 policies are applied to the group
    When edit the agent tags and use region:br, test:true orb tag(s)
    Then the container logs contain the message "completed RPC subscription to group" referred to each matching group within 10 seconds
        And this agent's heartbeat shows that 3 policies are successfully applied
        And this agent's heartbeat shows that 2 groups are matching the agent
        And the container logs contain the message "policy applied successfully" referred to each policy within 10 seconds
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds


Scenario: Agent subscription to group with policies after editing agent's tags (editing tags after agent provision)
    Given the Orb user has a registered account
        And the Orb user logs in
        And that an agent with 1 orb tag(s) already exists and is online
        And referred agent is subscribed to a group
        And that a sink already exists
        And 2 policies are applied to the group
        And this agent's heartbeat shows that 2 policies are successfully applied
        And an Agent Group is created with region:br orb tag(s)
        And 1 policies are applied to the group
    When edit the agent tags and use region:br orb tag(s)
    Then the container logs contain the message "completed RPC subscription to group" referred to each matching group within 10 seconds
        And this agent's heartbeat shows that 1 policies are successfully applied
        And this agent's heartbeat shows that 1 groups are matching the agent
        And the container logs contain the message "policy applied successfully" referred to each policy within 10 seconds
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds


Scenario: Insert tags in agents created without tags and apply policies to group matching new tags
    Given the Orb user has a registered account
        And the Orb user logs in
        And a new agent is created with 0 orb tag(s)
        And the agent container is started on port default
        And that a sink already exists
    When edit the agent tags and use 2 orb tag(s)
        And an Agent Group is created with same tag as the agent
        And 1 policies are applied to the group
    Then this agent's heartbeat shows that 1 policies are successfully applied
        And the container logs contain the message "completed RPC subscription to group" referred to each matching group within 10 seconds
        And this agent's heartbeat shows that 1 groups are matching the agent
        And the container logs contain the message "policy applied successfully" referred to each policy within 10 seconds
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds


Scenario: Edit agent name and apply policies to then
    Given the Orb user has a registered account
        And the Orb user logs in
        And that an agent with 5 orb tag(s) already exists and is online
        And an Agent Group is created with same tag as the agent
        And 1 agent must be matching on response field matching_agents
        And that a sink already exists
        And 1 policies are applied to the group
    When edit the agent name and edit agent tags using 3 orb tag(s)
    Then this agent's heartbeat shows that 1 policies are successfully applied
        And the container logs contain the message "policy applied successfully" referred to each policy within 10 seconds


Scenario: Editing tags of an Agent Group with policies (unsubscription - provision agent before editing)
    Given the Orb user has a registered account
        And the Orb user logs in
        And that an agent with region:br orb tag(s) already exists and is online
        And an Agent Group is created with same tag as the agent and without description
        And that a sink already exists
        And 2 policies are applied to the group
    When the name, tags, description of Agent Group is edited using: name=new_name/ tags=another:tag, ns1:true/ description=None
    Then 0 agent must be matching on response field matching_agents
        And the container logs should contain the message "completed RPC unsubscription to group" within 10 seconds
        And the agent status in Orb should be online


Scenario: Editing tags of an Agent Group with policies (subscription - provision agent before editing)
    Given the Orb user has a registered account
        And the Orb user logs in
        And that an agent with region:br, another:tag orb tag(s) already exists and is online
        And an Agent Group is created with ns1:true orb tag(s) and without description
        And that a sink already exists
        And 2 policies are applied to the group
    When the name, tags, description of Agent Group is edited using: name=new_name/ tags=region:br/ description=None
    Then 1 agent must be matching on response field matching_agents
        And the container logs should contain the message "completed RPC subscription to group" within 10 seconds
        And the agent status in Orb should be online
        And this agent's heartbeat shows that 1 groups are matching the agent
        And this agent's heartbeat shows that 2 policies are successfully applied


Scenario: Editing tags of an Agent Group with policies (provision agent after editing)
    Given the Orb user has a registered account
        And the Orb user logs in
        And an Agent Group is created with ns1:true orb tag(s) and without description
        And that a sink already exists
        And 2 policies are applied to the group
    When the name, tags, description of Agent Group is edited using: name=new_name/ tags=another:tag, ns1:true/ description=None
        And a new agent is created with region:us orb tag(s)
        And the agent container is started on port default
    Then 0 agent must be matching on response field matching_agents
        And the agent status in Orb should be online


Scenario: Editing tags of an Agent Group with policies (subscription - provision agent after editing)
    Given the Orb user has a registered account
        And the Orb user logs in
        And an Agent Group is created with ns1:true orb tag(s) and without description
        And that a sink already exists
        And 2 policies are applied to the group
    When the name, tags, description of Agent Group is edited using: name=new_name/ tags=region:br/ description=None
        And a new agent is created with region:br orb tag(s)
        And the agent container is started on port default
    Then 1 agent must be matching on response field matching_agents
        And the container logs should contain the message "completed RPC subscription to group" within 10 seconds
        And the agent status in Orb should be online
        And this agent's heartbeat shows that 1 groups are matching the agent
        And this agent's heartbeat shows that 2 policies are successfully applied


    Scenario: Editing tags of an Agent and Agent Group with policies (unsubscription - provision agent before editing)
    Given the Orb user has a registered account
        And the Orb user logs in
        And that an agent with region:br orb tag(s) already exists and is online
        And an Agent Group is created with same tag as the agent and without description
        And that a sink already exists
        And 2 policies are applied to the group
    When the name, tags, description of Agent Group is edited using: name=new_name/ tags=another:tag, ns1:true/ description=None
        And edit the agent tags and use region:us orb tag(s)
    Then 0 agent must be matching on response field matching_agents
        And the container logs should contain the message "completed RPC unsubscription to group" within 10 seconds
        And the agent status in Orb should be online


Scenario: Editing tags of an Agent and Agent Group with policies (subscription - provision agent before editing)
    Given the Orb user has a registered account
        And the Orb user logs in
        And that an agent with region:br, another:tag, test:true orb tag(s) already exists and is online
        And an Agent Group is created with ns1:true orb tag(s) and without description
        And that a sink already exists
        And 2 policies are applied to the group
    When edit the agent tags and use region:br, another:tag orb tag(s)
        And the name, tags, description of Agent Group is edited using: name=new_name/ tags=region:br/ description=None
    Then 1 agent must be matching on response field matching_agents
        And the container logs should contain the message "completed RPC subscription to group" within 10 seconds
        And the agent status in Orb should be online
        And this agent's heartbeat shows that 1 groups are matching the agent
        And this agent's heartbeat shows that 2 policies are successfully applied


Scenario: Editing tags of an Agent and Agent Group with policies (provision agent after editing)
    Given the Orb user has a registered account
        And the Orb user logs in
        And an Agent Group is created with ns1:true orb tag(s) and without description
        And that a sink already exists
        And 2 policies are applied to the group
    When the name, tags, description of Agent Group is edited using: name=new_name/ tags=another:tag, ns1:true/ description=None
        And a new agent is created with test:true orb tag(s)
        And edit the agent tags and use region:us orb tag(s)
        And the agent container is started on port default
    Then 0 agent must be matching on response field matching_agents
        And the agent status in Orb should be online


Scenario: Editing tags of an Agent and Agent Group with policies (subscription - provision agent after editing)
    Given the Orb user has a registered account
        And the Orb user logs in
        And an Agent Group is created with ns1:true orb tag(s) and without description
        And that a sink already exists
        And 2 policies are applied to the group
    When the name, tags, description of Agent Group is edited using: name=new_name/ tags=region:br/ description=None
        And a new agent is created with test:true orb tag(s)
        And edit the agent tags and use region:br orb tag(s)
        And the agent container is started on port default
    Then 1 agent must be matching on response field matching_agents
        And the container logs should contain the message "completed RPC subscription to group" within 10 seconds
        And the agent status in Orb should be online
        And this agent's heartbeat shows that 1 groups are matching the agent
        And this agent's heartbeat shows that 2 policies are successfully applied