@integration
Feature: Integration tests


@smoke
Scenario: Apply multiple advanced policies to an agent
    Given the Orb user has a registered account
        And the Orb user logs in
        And that an agent with 1 orb tag(s) already exists and is online
        And referred agent is subscribed to a group
        And that a sink already exists
    When 14 mixed policies are applied to the group
    Then this agent's heartbeat shows that 14 policies are successfully applied and has status running
        And the container logs contain the message "policy applied successfully" referred to each policy within 10 seconds
        And datasets related to all existing policies have validity valid


@smoke
Scenario: Apply two simple policies to an agent
    Given the Orb user has a registered account
        And the Orb user logs in
        And that an agent with 1 orb tag(s) already exists and is online
        And referred agent is subscribed to a group
        And that a sink already exists
    When 2 simple policies are applied to the group
    Then this agent's heartbeat shows that 2 policies are successfully applied and has status running
        And the container logs contain the message "policy applied successfully" referred to each policy within 10 seconds
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
        And referred sink must have active state on response within 10 seconds
        And datasets related to all existing policies have validity valid


@smoke
Scenario: apply one policy using multiple datasets to the same group
    Given the Orb user has a registered account
        And the Orb user logs in
        And that an agent with 2 orb tag(s) already exists and is online
        And referred agent is subscribed to a group
        And that a sink already exists
    When 2 simple policies are applied to the group by 3 datasets each
    Then this agent's heartbeat shows that 2 policies are successfully applied and has status running
        And 3 datasets are linked with each policy on agent's heartbeat
        And the container logs contain the message "policy applied successfully" referred to each policy within 180 seconds
        And referred sink must have active state on response within 180 seconds
        And datasets related to all existing policies have validity valid


@smoke
Scenario: Remove group to which agent is linked
    Given the Orb user has a registered account
        And the Orb user logs in
        And that an agent with 1 orb tag(s) already exists and is online
        And referred agent is subscribed to a group
        And that a sink already exists
        And 1 simple policies are applied to the group
        And this agent's heartbeat shows that 1 policies are successfully applied and has status running
    When the group to which the agent is linked is removed
    Then the container logs should contain the message "completed RPC unsubscription to group" within 10 seconds
        And this agent's heartbeat shows that 0 policies are successfully applied to the agent
        And this agent's heartbeat shows that 0 groups are matching the agent
        And dataset related have validity invalid


@smoke
Scenario: Remove policy from agent
    Given the Orb user has a registered account
        And the Orb user logs in
        And that an agent with 3 orb tag(s) already exists and is online
        And referred agent is subscribed to a group
        And that a sink already exists
        And 2 simple policies are applied to the group
        And this agent's heartbeat shows that 2 policies are successfully applied and has status running
    When one of applied policies is removed
    Then referred policy must not be listed on the orb policies list
        And datasets related to removed policy has validity invalid
        And datasets related to all existing policies have validity valid
        And this agent's heartbeat shows that 1 policies are successfully applied and has status running
        And container logs should inform that removed policy was stopped and removed within 10 seconds
        And the container logs that were output after the policy have been removed contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
        And the container logs that were output after the policy have been removed does not contain the message "scraped metrics for policy" referred to deleted policy anymore


@smoke
Scenario: Remove dataset from agent with just one dataset linked
    Given the Orb user has a registered account
        And the Orb user logs in
        And that an agent with 3 orb tag(s) already exists and is online
        And referred agent is subscribed to a group
        And that a sink already exists
        And 1 simple policies are applied to the group
        And this agent's heartbeat shows that 1 policies are successfully applied and has status running
    When a dataset linked to this agent is removed
    Then referred dataset must not be listed on the orb datasets list
        And this agent's heartbeat shows that 0 policies are successfully applied and has status running
        And container logs should inform that removed policy was stopped and removed within 10 seconds
        And the container logs that were output after removing dataset contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
        And the container logs that were output after removing dataset does not contain the message "scraped metrics for policy" referred to deleted policy anymore


@smoke
Scenario: Remove dataset from agent with more than one dataset linked
    Given the Orb user has a registered account
        And the Orb user logs in
        And that an agent with 4 orb tag(s) already exists and is online
        And referred agent is subscribed to a group
        And that a sink already exists
        And 3 simple policies are applied to the group
        And this agent's heartbeat shows that 3 policies are successfully applied and has status running
    When a dataset linked to this agent is removed
    Then referred dataset must not be listed on the orb datasets list
        And this agent's heartbeat shows that 2 policies are successfully applied and has status running
        And container logs should inform that removed policy was stopped and removed within 10 seconds
        And the container logs that were output after removing dataset contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
        And the container logs that were output after removing dataset does not contain the message "scraped metrics for policy" referred to deleted policy anymore


@smoke
Scenario: Provision agent with tags matching an existent group
    Given the Orb user has a registered account
        And the Orb user logs in
        And an Agent Group is created with 2 orb tag(s)
    When a new agent is created with orb tags matching 1 existing group
        And the agent container is started on an available port
    Then the agent status in Orb should be online
        And the container logs should contain the message "completed RPC subscription to group" within 10 seconds
        And this agent's heartbeat shows that 1 groups are matching the agent


@smoke
Scenario: Provision agent with tag matching existing group linked to a valid dataset
    Given the Orb user has a registered account
        And the Orb user logs in
        And an Agent Group is created with 3 orb tag(s)
        And that a sink already exists
        And 2 simple policies are applied to the group
    When a new agent is created with orb tags matching 1 existing group
        And the agent container is started on an available port
    Then this agent's heartbeat shows that 2 policies are successfully applied and has status running
        And the container logs contain the message "policy applied successfully" referred to each policy within 10 seconds
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
        And referred sink must have active state on response within 10 seconds
        And datasets related to all existing policies have validity valid


@smoke
Scenario: Sink with invalid endpoint
    Given the Orb user has a registered account
        And the Orb user logs in
        And that an agent with 1 orb tag(s) already exists and is online
        And referred agent is subscribed to a group
        And that a sink with invalid endpoint already exists
        And 3 simple policies are applied to the group
        And that a policy using: handler=dns, description='policy_dns', host_specification=10.0.1.0/24,10.0.2.1/32,2001:db8::/64, bpf_filter_expression=udp port 53, pcap_source=libpcap, only_qname_suffix=[.foo.com/ .example.com], only_rcode=2 already exists
    When a new dataset is created using referred group, policy and 1 sink
    Then this agent's heartbeat shows that 4 policies are successfully applied and has status running
        And the container logs should contain the message "managing agent policy from core" within 10 seconds
        And the container logs contain the message "policy applied successfully" referred to each policy within 10 seconds
        And the container logs should contain the message "scraped metrics for policy" within 180 seconds
        And referred sink must have error state on response within 10 seconds
        And dataset related have validity valid


@smoke
Scenario: Sink with invalid username
    Given the Orb user has a registered account
        And the Orb user logs in
        And that an agent with 1 orb tag(s) already exists and is online
        And referred agent is subscribed to a group
        And that a sink with invalid username already exists
        And 3 simple policies are applied to the group
        And that a policy using: handler=dns, description='policy_dns', host_specification=10.0.1.0/24,10.0.2.1/32,2001:db8::/64, bpf_filter_expression=udp port 53, pcap_source=libpcap, only_qname_suffix=[.foo.com/ .example.com], only_rcode=3 already exists
    When a new dataset is created using referred group, policy and 1 sink
    Then the container logs should contain the message "managing agent policy from core" within 10 seconds
        And this agent's heartbeat shows that 4 policies are successfully applied and has status running
        And the container logs contain the message "policy applied successfully" referred to each policy within 10 seconds
        And the container logs should contain the message "scraped metrics for policy" within 180 seconds
        And referred sink must have error state on response within 10 seconds
        And dataset related have validity valid


@smoke
Scenario: Sink with invalid password
    Given the Orb user has a registered account
        And the Orb user logs in
        And that an agent with 1 orb tag(s) already exists and is online
        And referred agent is subscribed to a group
        And that a sink with invalid password already exists
        And 3 simple policies are applied to the group
        And that a policy using: handler=dns, description='policy_dns', host_specification=10.0.1.0/24,10.0.2.1/32,2001:db8::/64, bpf_filter_expression=udp port 53, pcap_source=libpcap, only_qname_suffix=[.foo.com/ .example.com], only_rcode=5 already exists
    When a new dataset is created using referred group, policy and 1 sink
    Then this agent's heartbeat shows that 4 policies are successfully applied and has status running
        And the container logs should contain the message "managing agent policy from core" within 10 seconds
        And the container logs contain the message "policy applied successfully" referred to each policy within 10 seconds
        And the container logs should contain the message "scraped metrics for policy" within 180 seconds
        And referred sink must have error state on response within 10 seconds
        And dataset related have validity valid


@smoke
Scenario: Agent subscription to multiple groups created after provisioning agent
    Given the Orb user has a registered account
        And the Orb user logs in
        And a new agent is created with 3 orb tag(s)
        And the agent container is started on an available port
    When an Agent Group is created with 2 tags contained in the agent
        And an Agent Group is created with 1 tags contained in the agent
    Then the container logs contain the message "completed RPC subscription to group" referred to each matching group within 10 seconds
        And this agent's heartbeat shows that 2 groups are matching the agent


@smoke
Scenario: Agent subscription to multiple groups created before provisioning agent
    Given the Orb user has a registered account
        And the Orb user logs in
        And an Agent Group is created with 2 orb tag(s)
        And an Agent Group is created with 1 orb tag(s)
        And an Agent Group is created with 1 orb tag(s)
    When a new agent is created with orb tags matching all existing group
        And the agent container is started on an available port
    Then the container logs contain the message "completed RPC subscription to group" referred to each matching group within 10 seconds
        And this agent's heartbeat shows that 3 groups are matching the agent


@smoke
Scenario: Agent subscription to group after editing agent's tags (agent provisioned before editing and group created after)
    Given the Orb user has a registered account
        And the Orb user logs in
        And a new agent is created with 1 orb tag(s)
        And the agent container is started on an available port
    When edit the agent tags and use 2 orb tag(s)
        And an Agent Group is created with all tags contained in the agent
    Then this agent's heartbeat shows that 1 groups are matching the agent
        And the container logs contain the message "completed RPC subscription to group" referred to each matching group within 30 seconds


@smoke
Scenario: Agent subscription to group after editing agent's tags (editing tags after agent provision)
    Given the Orb user has a registered account
        And the Orb user logs in
        And an Agent Group is created with 2 orb tag(s)
        And an Agent Group is created with 1 orb tag(s)
        And a new agent is created with orb tags matching 1 existing group
        And the agent container is started on an available port
    When edit the agent tags and use orb tags matching 2 existing group
    Then the container logs contain the message "completed RPC subscription to group" referred to each matching group within 30 seconds
        And this agent's heartbeat shows that 2 groups are matching the agent


@smoke
Scenario: Agent subscription to group after editing agent's tags (editing tags before agent provision)
    Given the Orb user has a registered account
        And the Orb user logs in
        And a new agent is created with 2 orb tag(s)
        And edit the agent tags and use 1 orb tag(s)
        And the agent container is started on an available port
    When an Agent Group is created with 1 tags contained in the agent
    Then the container logs contain the message "completed RPC subscription to group" referred to each matching group within 10 seconds
        And this agent's heartbeat shows that 1 groups are matching the agent


@smoke
Scenario: Agent subscription to multiple group with policies after editing agent's tags (editing tags after agent provision)
    Given the Orb user has a registered account
        And the Orb user logs in
        And that an agent with 1 orb tag(s) already exists and is online
        And referred agent is subscribed to a group
        And that a sink already exists
        And 2 simple policies are applied to the group
        And this agent's heartbeat shows that 2 policies are successfully applied and has status running
        And an Agent Group is created with 1 orb tag(s)
        And 1 simple policies are applied to the group
    When edit the agent tags and use orb tags matching all existing group
    Then the container logs contain the message "completed RPC subscription to group" referred to each matching group within 10 seconds
        And this agent's heartbeat shows that 3 policies are successfully applied and has status running
        And this agent's heartbeat shows that 2 groups are matching the agent
        And the container logs contain the message "policy applied successfully" referred to each policy within 10 seconds
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds


@smoke
Scenario: Agent subscription to group with policies after editing agent's tags (editing tags after agent provision)
    Given the Orb user has a registered account
        And the Orb user logs in
        And that an agent with 1 orb tag(s) already exists and is online
        And referred agent is subscribed to a group
        And that a sink already exists
        And 2 simple policies are applied to the group
        And this agent's heartbeat shows that 2 policies are successfully applied and has status running
        And an Agent Group is created with 1 orb tag(s)
        And 1 simple policies are applied to the group
    When edit the agent tags and use orb tags matching last existing group
    Then the container logs contain the message "completed RPC subscription to group" referred to each matching group within 10 seconds
        And this agent's heartbeat shows that 1 groups are matching the agent
        And this agent's heartbeat shows that 1 policies are successfully applied and has status running
        And the container logs contain the message "policy applied successfully" referred to each policy within 10 seconds
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds


@smoke
Scenario: Insert tags in agents created without tags and apply policies to group matching new tags
    Given the Orb user has a registered account
        And the Orb user logs in
        And a new agent is created with 0 orb tag(s)
        And the agent container is started on an available port
        And that a sink already exists
    When edit the agent tags and use 2 orb tag(s)
        And an Agent Group is created with same tag as the agent and without description
        And 1 simple policies are applied to the group
    Then this agent's heartbeat shows that 1 policies are successfully applied and has status running
        And the container logs contain the message "completed RPC subscription to group" referred to each matching group within 10 seconds
        And this agent's heartbeat shows that 1 groups are matching the agent
        And the container logs contain the message "policy applied successfully" referred to each policy within 10 seconds
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds


@smoke
Scenario: Edit agent name and apply policies to then
    Given the Orb user has a registered account
        And the Orb user logs in
        And that an agent with 5 orb tag(s) already exists and is online
        And an Agent Group is created with all tags contained in the agent
        And 1 agent must be matching on response field matching_agents
        And that a sink already exists
        And 1 simple policies are applied to the group
    When edit the agent name and edit agent tags using 3 orb tag(s)
    Then this agent's heartbeat shows that 1 policies are successfully applied and has status running
        And the container logs contain the message "policy applied successfully" referred to each policy within 10 seconds


@smoke
Scenario: Editing tags of an Agent Group with policies (unsubscription - provision agent before editing)
    Given the Orb user has a registered account
        And the Orb user logs in
        And that an agent with 1 orb tag(s) already exists and is online
        And an Agent Group is created with same tag as the agent and without description
        And that a sink already exists
        And 2 simple policies are applied to the group
    When the name, tags, description of Agent Group is edited using: name=new_name/ tags=2 orb tag(s)/ description=None
    Then 0 agent must be matching on response field matching_agents
        And the container logs should contain the message "completed RPC unsubscription to group" within 10 seconds
        And the agent status in Orb should be online


@smoke
Scenario: Editing tags of an Agent Group with policies (subscription - provision agent before editing)
    Given the Orb user has a registered account
        And the Orb user logs in
        And that an agent with 2 orb tag(s) already exists and is online
        And an Agent Group is created with 1 orb tag(s) and without description
        And that a sink already exists
        And 2 simple policies are applied to the group
    When the name, tags, description of Agent Group is edited using: name=new_name/ tags=matching the agent/ description=None
    Then 1 agent must be matching on response field matching_agents
        And the container logs should contain the message "completed RPC subscription to group" within 10 seconds
        And the agent status in Orb should be online
        And this agent's heartbeat shows that 1 groups are matching the agent
        And this agent's heartbeat shows that 2 policies are successfully applied and has status running


@smoke
Scenario: Editing tags of an Agent Group with policies (provision agent after editing)
    Given the Orb user has a registered account
        And the Orb user logs in
        And an Agent Group is created with 1 orb tag(s) and without description
        And that a sink already exists
        And 2 simple policies are applied to the group
        And a new agent is created with orb tags matching 1 existing group
        And 1 agent must be matching on response field matching_agents
    When the name, tags, description of Agent Group is edited using: name=new_name/ tags=2 orb tag/ description=None
        And the agent container is started on an available port
    Then 0 agent must be matching on response field matching_agents
        And the agent status in Orb should be online


@smoke
Scenario: Editing tags of an Agent Group with policies (subscription - provision agent after editing)
    Given the Orb user has a registered account
        And the Orb user logs in
        And an Agent Group is created with 1 orb tag(s) and without description
        And that a sink already exists
        And 2 simple policies are applied to the group
    When the name, tags, description of Agent Group is edited using: name=new_name/ tags=2 orb tag(s)/ description=None
        And a new agent is created with orb tags matching 1 existing group
        And the agent container is started on an available port
    Then 1 agent must be matching on response field matching_agents
        And the container logs should contain the message "completed RPC subscription to group" within 10 seconds
        And the agent status in Orb should be online
        And this agent's heartbeat shows that 1 groups are matching the agent
        And this agent's heartbeat shows that 2 policies are successfully applied and has status running


@smoke
Scenario: Editing tags of an Agent and Agent Group with policies (unsubscription - provision agent before editing)
    Given the Orb user has a registered account
        And the Orb user logs in
        And that an agent with 1 orb tag(s) already exists and is online
        And an Agent Group is created with same tag as the agent and without description
        And that a sink already exists
        And 2 simple policies are applied to the group
    When the name, tags, description of Agent Group is edited using: name=new_name/ tags=2 orb tag(s)/ description=None
        And edit the agent tags and use 1 orb tag(s)
    Then 0 agent must be matching on response field matching_agents
        And the container logs should contain the message "completed RPC unsubscription to group" within 10 seconds
        And this agent's heartbeat shows that 0 groups are matching the agent
        And this agent's heartbeat shows that 0 policies are successfully applied to the agent
        And the agent status in Orb should be online


@smoke
Scenario: Editing tags of an Agent and Agent Group with policies (subscription - provision agent before editing)
    Given the Orb user has a registered account
        And the Orb user logs in
        And that an agent with 3 orb tag(s) already exists and is online
        And an Agent Group is created with 1 orb tag(s) and without description
        And that a sink already exists
        And 2 simple policies are applied to the group
    When edit the agent tags and use 2 orb tag(s)
        And the name, tags, description of Agent Group is edited using: name=new_name/ tags=matching the agent/ description=None
    Then 1 agent must be matching on response field matching_agents
        And the container logs should contain the message "completed RPC subscription to group" within 10 seconds
        And the agent status in Orb should be online
        And this agent's heartbeat shows that 1 groups are matching the agent
        And this agent's heartbeat shows that 2 policies are successfully applied and has status running


@smoke
Scenario: Editing tags of an Agent and Agent Group with policies (provision agent after editing)
    Given the Orb user has a registered account
        And the Orb user logs in
        And an Agent Group is created with 1 orb tag(s) and without description
        And that a sink already exists
        And 2 simple policies are applied to the group
    When the name, tags, description of Agent Group is edited using: name=new_name/ tags=2 orb tag(s)/ description=None
        And a new agent is created with 1 orb tag(s)
        And edit the agent tags and use 1 orb tag(s)
        And the agent container is started on an available port
    Then 0 agent must be matching on response field matching_agents
        And the agent status in Orb should be online


@smoke
Scenario: Editing tags of an Agent and Agent Group with policies (subscription - provision agent after editing)
    Given the Orb user has a registered account
        And the Orb user logs in
        And an Agent Group is created with 1 orb tag(s) and without description
        And that a sink already exists
        And 2 simple policies are applied to the group
    When the name, tags, description of Agent Group is edited using: name=new_name/ tags=1 orb tag/ description=None
        And a new agent is created with 2 orb tag(s)
        And edit the agent tags and use orb tags matching 1 existing group
        And the agent container is started on an available port
    Then 1 agent must be matching on response field matching_agents
        And the container logs should contain the message "completed RPC subscription to group" within 10 seconds
        And the agent status in Orb should be online
        And this agent's heartbeat shows that 1 groups are matching the agent
        And this agent's heartbeat shows that 2 policies are successfully applied and has status running



@smoke
Scenario: Edit an advanced policy with handler dns changing the handler to net
    Given the Orb user has a registered account
        And the Orb user logs in
        And that a sink already exists
        And that an agent with 1 orb tag(s) already exists and is online
        And an Agent Group is created with all tags contained in the agent
        And a new policy is created using: handler=dns, description='policy_dns', bpf_filter_expression=udp port 53, pcap_source=libpcap, only_qname_suffix=[.foo.com/ .example.com], only_rcode=0
        And a new dataset is created using referred group, policy and 1 sink
    When editing a policy using name=my_policy, handler=net, only_qname_suffix=None, only_rcode=None
    Then policy version must be 1
        And policy name must be my_policy
        And policy handler must be net
        And policy only_qname_suffix must be None
        And policy only_rcode must be None
        And this agent's heartbeat shows that 1 policies are successfully applied and has status running



@smoke
Scenario: Edit an advanced policy with handler dns changing the handler to dhcp
    Given the Orb user has a registered account
        And the Orb user logs in
        And that a sink already exists
        And that an agent with 1 orb tag(s) already exists and is online
        And an Agent Group is created with all tags contained in the agent
        And a new policy is created using: handler=dns, host_specification=10.0.1.0/24,10.0.2.1/32,2001:db8::/64, bpf_filter_expression=udp port 53, pcap_source=libpcap, only_qname_suffix=[.foo.com/ .example.com], only_rcode=2
        And a new dataset is created using referred group, policy and 1 sink
    When editing a policy using name=second_policy, handler=dhcp, only_qname_suffix=None, only_rcode=None
    Then policy version must be 1
        And policy name must be second_policy
        And policy handler must be dhcp
        And this agent's heartbeat shows that 1 policies are successfully applied and has status running


@smoke
Scenario: Edit a simple policy with handler dhcp changing the handler to net
    Given the Orb user has a registered account
        And the Orb user logs in
        And that a sink already exists
        And that an agent with 1 orb tag(s) already exists and is online
        And an Agent Group is created with all tags contained in the agent
        And a new policy is created using: handler=dhcp
        And a new dataset is created using referred group, policy and 1 sink
    When editing a policy using handler=net, description="policy_net"
    Then policy version must be 1
        And policy handler must be net
        And this agent's heartbeat shows that 1 policies are successfully applied and has status running


@smoke
Scenario: Edit a simple policy with handler net changing the handler to dns and inserting advanced parameters
    Given the Orb user has a registered account
        And the Orb user logs in
        And that a sink already exists
        And that an agent with 1 orb tag(s) already exists and is online
        And an Agent Group is created with all tags contained in the agent
        And a new policy is created using: handler=net
        And a new dataset is created using referred group, policy and 1 sink
    When editing a policy using handler=dns, host_specification=10.0.1.0/24,10.0.2.1/32,2001:db8::/64, bpf_filter_expression=udp port 53, pcap_source=libpcap, only_qname_suffix=[.foo.com/ .example.com], only_rcode=2
    Then policy version must be 1
        And policy handler must be dns
        And policy host_specification must be 10.0.1.0/24,10.0.2.1/32,2001:db8::/64
        And policy bpf_filter_expression must be udp port 53
        And policy only_qname_suffix must be ['.foo.com', '.example.com']
        And policy only_rcode must be 2
        And this agent's heartbeat shows that 1 policies are successfully applied and has status running


@smoke
Scenario: remove 1 sink from a dataset with 2 sinks
    Given the Orb user has a registered account
        And the Orb user logs in
        And that 2 sinks already exists
        And that an agent with 1 orb tag(s) already exists and is online
        And an Agent Group is created with all tags contained in the agent
        And a new policy is created using: handler=dhcp
        And a new dataset is created using referred group, policy and 2 sinks
    When remove 1 of the linked sinks from orb
    Then dataset related have validity valid
        And this agent's heartbeat shows that 1 policies are successfully applied and has status running


@smoke
Scenario: remove 1 sink from a dataset with 1 sinks
    Given the Orb user has a registered account
        And the Orb user logs in
        And that 2 sinks already exists
        And that an agent with 1 orb tag(s) already exists and is online
        And an Agent Group is created with all tags contained in the agent
        And a new policy is created using: handler=dhcp
        And a new dataset is created using referred group, policy and 1 sinks
    When remove 1 of the linked sinks from orb
    Then dataset related have validity invalid
        And this agent's heartbeat shows that 0 policies are successfully applied to the agent
        And the container logs should contain the message "completed RPC subscription to group" within 10 seconds


@smoke
Scenario: remove one sink from a dataset with 1 sinks, edit the dataset and insert another sink
    Given the Orb user has a registered account
        And the Orb user logs in
        And that 2 sinks already exists
        And that an agent with 1 orb tag(s) already exists and is online
        And an Agent Group is created with all tags contained in the agent
        And a new policy is created using: handler=dns
        And a new dataset is created using referred group, policy and 1 sinks
        And remove 1 of the linked sinks from orb
        And dataset related have validity invalid
        And this agent's heartbeat shows that 0 policies are successfully applied to the agent
    When the dataset is edited and 1 sinks are linked
    Then dataset related have validity valid
        And this agent's heartbeat shows that 1 policies are successfully applied and has status running


@smoke
Scenario: agent with only agent tags subscription to a group with policies created after provision the agent (config file)
    Given the Orb user has a registered account
        And the Orb user logs in
        And that a sink already exists
    When an agent is self-provisioned via a configuration file on port available with 3 agent tags and has status online
        And an Agent Group is created with all tags contained in the agent
        And 3 simple policies are applied to the group
    Then dataset related have validity valid
        And this agent's heartbeat shows that 1 groups are matching the agent
        And the container logs should contain the message "completed RPC subscription to group" within 10 seconds
        And this agent's heartbeat shows that 3 policies are successfully applied and has status running
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
        And referred sink must have active state on response within 10 seconds
        And the container logs contain the message "policy applied successfully" referred to each policy within 10 seconds
        And remove all the agents .yaml generated on test process


@smoke
Scenario: agent with only agent tags subscription to a group with policies created before provision the agent
    Given the Orb user has a registered account
        And the Orb user logs in
        And that a sink already exists
        And an Agent Group is created with 1 orb tag(s)
        And 3 simple policies are applied to the group
    When an agent is self-provisioned via a configuration file on port available with matching 1 group agent tags and has status online
    Then dataset related have validity valid
        And this agent's heartbeat shows that 1 groups are matching the agent
        And the container logs should contain the message "completed RPC subscription to group" within 10 seconds
        And this agent's heartbeat shows that 3 policies are successfully applied and has status running
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
        And referred sink must have active state on response within 10 seconds
        And the container logs contain the message "policy applied successfully" referred to each policy within 10 seconds
        And remove all the agents .yaml generated on test process
        And remove all the agents .yaml generated on test process


@smoke
Scenario: agent with mixed tags subscription to a group with policies created after provision the agent (config file)
    Given the Orb user has a registered account
        And the Orb user logs in
        And that a sink already exists
    When an agent is self-provisioned via a configuration file on port available with 3 agent tags and has status online
        And edit the agent tags and use 2 orb tag(s)
        And an Agent Group is created with all tags contained in the agent
        And 3 simple policies are applied to the group
    Then dataset related have validity valid
        And this agent's heartbeat shows that 1 groups are matching the agent
        And the container logs should contain the message "completed RPC subscription to group" within 10 seconds
        And this agent's heartbeat shows that 3 policies are successfully applied and has status running
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
        And referred sink must have active state on response within 10 seconds
        And the container logs contain the message "policy applied successfully" referred to each policy within 10 seconds
        And remove all the agents .yaml generated on test process


@smoke
Scenario: agent with mixed tags subscription to a group with policies created before provision the agent (config file)
    Given the Orb user has a registered account
        And the Orb user logs in
        And that a sink already exists
        And an Agent Group is created with 2 orb tag(s)
        And 3 simple policies are applied to the group
    When an agent is self-provisioned via a configuration file on port available with 1 agent tags and has status online
        And edit the agent tags and use orb tags matching 1 existing group
    Then dataset related have validity valid
        And this agent's heartbeat shows that 1 groups are matching the agent
        And the container logs should contain the message "completed RPC subscription to group" within 10 seconds
        And this agent's heartbeat shows that 3 policies are successfully applied and has status running
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
        And referred sink must have active state on response within 10 seconds
        And the container logs contain the message "policy applied successfully" referred to each policy within 10 seconds
        And remove all the agents .yaml generated on test process

@smoke
Scenario: Remotely restart agents with policies applied
    Given the Orb user has a registered account
        And the Orb user logs in
        And that an agent with 1 orb tag(s) already exists and is online
        And referred agent is subscribed to a group
        And that a sink already exists
        And 2 simple policies are applied to the group
        And this agent's heartbeat shows that 2 policies are successfully applied and has status running
    When remotely restart the agent
    Then the container logs should contain the message "restarting all backends" within 5 seconds
        And the container logs that were output after reset the agent contain the message "removing policies" within 5 seconds
        And the container logs that were output after reset the agent contain the message "resetting backend" within 5 seconds
        And the container logs that were output after reset the agent contain the message "pktvisor process stopped" within 5 seconds
        And the container logs that were output after reset the agent contain the message "reapplying policies" within 5 seconds
        And the container logs that were output after reset the agent contain the message "all backends were restarted" within 5 seconds
        And the container logs that were output after reset the agent contain the message "policy applied successfully" referred to each applied policy within 10 seconds
        And the container logs that were output after reset the agent contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds

@smoke
Scenario: Remotely restart agents without policies applied
    Given the Orb user has a registered account
        And the Orb user logs in
        And that an agent with 1 orb tag(s) already exists and is online
        And referred agent is subscribed to a group
        And that a sink already exists
    When remotely restart the agent
        And the container logs that were output after reset the agent contain the message "resetting backend" within 5 seconds
        And the container logs that were output after reset the agent contain the message "pktvisor process stopped" within 5 seconds
        And the container logs that were output after reset the agent contain the message "all backends were restarted" within 5 seconds
        And 2 simple policies are applied to the group
    Then the container logs should contain the message "restarting all backends" within 5 seconds
        And this agent's heartbeat shows that 2 policies are successfully applied and has status running
        And the container logs that were output after reset the agent contain the message "policy applied successfully" referred to each applied policy within 20 seconds
        And the container logs that were output after reset the agent contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds