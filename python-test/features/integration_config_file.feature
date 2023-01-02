@integration_config_files @AUTORETRY
Feature: Integration tests using agent provided via config file

########### provisioning agents without specify pktvisor configs on backend

@smoke @config_file @pktvisor_configs
Scenario: provisioning agent without specify pktvisor binary path and path to config file (config file - auto_provision=true)
    Given the Orb user has a registered account
        And the Orb user logs in
        And that a sink already exists
        And 1 Agent Group(s) is created with 1 orb tag(s) (lower case)
        And 3 simple policies flow are applied to the group
        And a new agent is created with 0 orb tag(s)
    When an agent(input_type:flow, settings: {"bind":"0.0.0.0", "port":"available_port"}) is self-provisioned via a configuration file on port available with matching 1 group agent tags and has status online. [Overwrite default: True. Paste only file: True. Use specif backend for pktvisor {"binary":"None", "config_file":"None"}]
        And pktvisor state is running
    Then 3 dataset(s) have validity valid and 0 have validity invalid in 30 seconds
        And this agent's heartbeat shows that 1 groups are matching the agent
        And the container logs should contain the message "completed RPC subscription to group" within 30 seconds
        And this agent's heartbeat shows that 3 policies are applied and all has status running
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
        And referred sink must have active state on response within 30 seconds
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And remove the agent .yaml generated on each scenario


@smoke @config_file @pktvisor_configs
Scenario: provisioning agent without specify pktvisor binary path (config file - auto_provision=true)
    Given the Orb user has a registered account
        And the Orb user logs in
        And that a sink already exists
        And 1 Agent Group(s) is created with 1 orb tag(s) (lower case)
        And 3 simple policies flow are applied to the group
        And a new agent is created with 0 orb tag(s)
    When an agent(input_type:flow, settings: {"bind":"0.0.0.0", "port":"available_port"}) is self-provisioned via a configuration file on port available with matching 1 group agent tags and has status online. [Overwrite default: True. Paste only file: True. Use specif backend for pktvisor {"config_file":"None"}]
        And pktvisor state is running
    Then 3 dataset(s) have validity valid and 0 have validity invalid in 30 seconds
        And this agent's heartbeat shows that 1 groups are matching the agent
        And the container logs should contain the message "completed RPC subscription to group" within 30 seconds
        And this agent's heartbeat shows that 3 policies are applied and all has status running
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
        And referred sink must have active state on response within 30 seconds
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And remove the agent .yaml generated on each scenario


@smoke @config_file @pktvisor_configs
Scenario: provisioning agent without specify pktvisor path to config file (config file - auto_provision=true)
    Given the Orb user has a registered account
        And the Orb user logs in
        And that a sink already exists
        And 1 Agent Group(s) is created with 1 orb tag(s) (lower case)
        And 3 simple policies flow are applied to the group
        And a new agent is created with 0 orb tag(s)
    When an agent(input_type:flow, settings: {"bind":"0.0.0.0", "port":"available_port"}) is self-provisioned via a configuration file on port available with matching 1 group agent tags and has status online. [Overwrite default: True. Paste only file: True. Use specif backend for pktvisor {"binary":"None"}]
        And pktvisor state is running
    Then 3 dataset(s) have validity valid and 0 have validity invalid in 30 seconds
        And this agent's heartbeat shows that 1 groups are matching the agent
        And the container logs should contain the message "completed RPC subscription to group" within 30 seconds
        And this agent's heartbeat shows that 3 policies are applied and all has status running
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
        And referred sink must have active state on response within 30 seconds
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And remove the agent .yaml generated on each scenario

@smoke @config_file @pktvisor_configs
Scenario: provisioning agent without specify pktvisor binary path and path to config file (config file - auto_provision=false)
    Given the Orb user has a registered account
        And the Orb user logs in
        And that a sink already exists
        And a new agent is created with 2 orb tag(s)
    When an agent(input_type:flow, settings: {"bind":"0.0.0.0", "port":"available_port"}) is provisioned via a configuration file on port available with 3 agent tags and has status online. [Overwrite default: True. Paste only file: True. Use specif backend for pktvisor {"binary":"None", "config_file":"None"}]
        And pktvisor state is running
        And edit the orb tags on agent and use 2 orb tag(s)
        And 1 Agent Group(s) is created with all tags contained in the agent
        And 3 simple policies same input_type as created via config file are applied to the group
    Then 3 dataset(s) have validity valid and 0 have validity invalid in 30 seconds
        And this agent's heartbeat shows that 1 groups are matching the agent
        And the container logs should contain the message "completed RPC subscription to group" within 30 seconds
        And this agent's heartbeat shows that 3 policies are applied and all has status running
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
        And referred sink must have active state on response within 30 seconds
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And remove the agent .yaml generated on each scenario



@smoke @config_file @pktvisor_configs
Scenario: provisioning agent without specify pktvisor binary path (config file - auto_provision=false)
    Given the Orb user has a registered account
        And the Orb user logs in
        And that a sink already exists
        And a new agent is created with 2 orb tag(s)
    When an agent(input_type:flow, settings: {"bind":"0.0.0.0", "port":"available_port"}) is provisioned via a configuration file on port available with 3 agent tags and has status online. [Overwrite default: True. Paste only file: True. Use specif backend for pktvisor {"config_file":"None"}]
        And pktvisor state is running
        And edit the orb tags on agent and use 2 orb tag(s)
        And 1 Agent Group(s) is created with all tags contained in the agent
        And 3 simple policies same input_type as created via config file are applied to the group
    Then 3 dataset(s) have validity valid and 0 have validity invalid in 30 seconds
        And this agent's heartbeat shows that 1 groups are matching the agent
        And the container logs should contain the message "completed RPC subscription to group" within 30 seconds
        And this agent's heartbeat shows that 3 policies are applied and all has status running
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
        And referred sink must have active state on response within 30 seconds
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And remove the agent .yaml generated on each scenario



@smoke @config_file @pktvisor_configs
Scenario: provisioning agent without specify pktvisor path to config file (config file - auto_provision=false)
    Given the Orb user has a registered account
        And the Orb user logs in
        And that a sink already exists
        And a new agent is created with 2 orb tag(s)
    When an agent(input_type:flow, settings: {"bind":"0.0.0.0", "port":"available_port"}) is provisioned via a configuration file on port available with 3 agent tags and has status online. [Overwrite default: True. Paste only file: True. Use specif backend for pktvisor {"binary":"None"}]
        And pktvisor state is running
        And edit the orb tags on agent and use 2 orb tag(s)
        And 1 Agent Group(s) is created with all tags contained in the agent
        And 3 simple policies same input_type as created via config file are applied to the group
    Then 3 dataset(s) have validity valid and 0 have validity invalid in 30 seconds
        And this agent's heartbeat shows that 1 groups are matching the agent
        And the container logs should contain the message "completed RPC subscription to group" within 30 seconds
        And this agent's heartbeat shows that 3 policies are applied and all has status running
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
        And referred sink must have active state on response within 30 seconds
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And remove the agent .yaml generated on each scenario


########### tap_selector

@smoke @config_file
Scenario: tap_selector - any - matching 0 of all tags from an agent
    Given the Orb user has a registered account
        And the Orb user logs in
        And that a sink already exists
    When an agent(input_type:pcap, settings: {"iface":"default"}) is self-provisioned via a configuration file on port available with 3 agent tags and has status online. [Overwrite default: False. Paste only file: True]
        And pktvisor state is running
        And 1 Agent Group(s) is created with all tags contained in the agent
        And a net policy pcap with tap_selector matching any of 0 agent tap tags ands settings: geoloc_notfound=False is applied to the group
    Then 1 dataset(s) have validity valid and 0 have validity invalid in 30 seconds
        And the container logs should contain the message "completed RPC subscription to group" within 30 seconds
        And this agent's heartbeat shows that 1 groups are matching the agent
        And this agent's heartbeat shows that 1 policies are applied and all has status failed_to_apply
        And the policy application error details must show that 422 no tap match found for specified 'input.tap_selector' tags
        And remove the agent .yaml generated on each scenario

@smoke @config_file
Scenario: tap_selector - any - matching 1 of all tags from an agent
    Given the Orb user has a registered account
        And the Orb user logs in
        And that a sink already exists
    When an agent(input_type:pcap, settings: {"iface":"default"}) is self-provisioned via a configuration file on port available with 3 agent tags and has status online. [Overwrite default: False. Paste only file: False]
        And pktvisor state is running
        And 1 Agent Group(s) is created with all tags contained in the agent
        And a net policy pcap with tap_selector matching any of 1 agent (1 tag matching) tap tags ands settings: geoloc_notfound=False is applied to the group
    Then 1 dataset(s) have validity valid and 0 have validity invalid in 30 seconds
        And the container logs should contain the message "completed RPC subscription to group" within 30 seconds
        And this agent's heartbeat shows that 1 groups are matching the agent
        And this agent's heartbeat shows that 1 policies are applied and all has status running
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
        And referred sink must have active state on response within 30 seconds
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And remove the agent .yaml generated on each scenario

@smoke @config_file
Scenario: tap_selector - any - matching 1 of all tags (plus 1 random tag) from an agent
    Given the Orb user has a registered account
        And the Orb user logs in
        And that a sink already exists
    When an agent(input_type:pcap, settings: {"iface":"default"}) is self-provisioned via a configuration file on port available with 3 agent tags and has status online. [Overwrite default: False. Paste only file: True]
        And pktvisor state is running
        And 1 Agent Group(s) is created with all tags contained in the agent
        And a net policy pcap with tap_selector matching any of 1 agent (1 tag matching + 1 random tag) tap tags ands settings: geoloc_notfound=False is applied to the group
    Then 1 dataset(s) have validity valid and 0 have validity invalid in 30 seconds
        And the container logs should contain the message "completed RPC subscription to group" within 30 seconds
        And this agent's heartbeat shows that 1 groups are matching the agent
        And this agent's heartbeat shows that 1 policies are applied and all has status running
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
        And referred sink must have active state on response within 30 seconds
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And remove the agent .yaml generated on each scenario

@smoke @config_file
Scenario: tap_selector - all - matching 0 of all tags from an agent
    Given the Orb user has a registered account
        And the Orb user logs in
        And that a sink already exists
    When an agent(input_type:pcap, settings: {"iface":"default"}) is self-provisioned via a configuration file on port available with 3 agent tags and has status online. [Overwrite default: False. Paste only file: False]
        And pktvisor state is running
        And 1 Agent Group(s) is created with all tags contained in the agent
        And a net policy pcap with tap_selector matching all of 0 agent tap tags ands settings: geoloc_notfound=False is applied to the group
    Then 1 dataset(s) have validity valid and 0 have validity invalid in 30 seconds
        And the container logs should contain the message "completed RPC subscription to group" within 30 seconds
        And this agent's heartbeat shows that 1 groups are matching the agent
        And this agent's heartbeat shows that 1 policies are applied and all has status failed_to_apply
        And the policy application error details must show that 422 no tap match found for specified 'input.tap_selector' tags
        And remove the agent .yaml generated on each scenario

@smoke @config_file
Scenario: tap_selector - all - matching 1 of all tags from an agent
    Given the Orb user has a registered account
        And the Orb user logs in
        And that a sink already exists
    When an agent(input_type:pcap, settings: {"iface":"default"}) is self-provisioned via a configuration file on port available with 3 agent tags and has status online. [Overwrite default: False. Paste only file: True]
        And pktvisor state is running
        And 1 Agent Group(s) is created with all tags contained in the agent
        And a net policy pcap with tap_selector matching all of 1 agent (1 tag matching) tap tags ands settings: geoloc_notfound=False is applied to the group
    Then 1 dataset(s) have validity valid and 0 have validity invalid in 30 seconds
        And the container logs should contain the message "completed RPC subscription to group" within 30 seconds
        And this agent's heartbeat shows that 1 groups are matching the agent
        And this agent's heartbeat shows that 1 policies are applied and all has status running
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
        And referred sink must have active state on response within 30 seconds
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And remove the agent .yaml generated on each scenario


@smoke @config_file
Scenario: tap_selector - all - matching all tags from an agent
    Given the Orb user has a registered account
        And the Orb user logs in
        And that a sink already exists
    When an agent(input_type:pcap, settings: {"iface":"default"}) is self-provisioned via a configuration file on port available with 3 agent tags and has status online. [Overwrite default: False. Paste only file: False]
        And pktvisor state is running
        And 1 Agent Group(s) is created with all tags contained in the agent
        And a net policy pcap with tap_selector matching all of an agent tap tags ands settings: geoloc_notfound=False is applied to the group
    Then 1 dataset(s) have validity valid and 0 have validity invalid in 30 seconds
        And the container logs should contain the message "completed RPC subscription to group" within 30 seconds
        And this agent's heartbeat shows that 1 groups are matching the agent
        And this agent's heartbeat shows that 1 policies are applied and all has status running
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
        And referred sink must have active state on response within 30 seconds
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And remove the agent .yaml generated on each scenario

########### pcap

@smoke @config_file
Scenario: agent pcap with only agent tags subscription to a group with policies created after provision the agent (config file - auto_provision=true)
    Given the Orb user has a registered account
        And the Orb user logs in
        And that a sink already exists
    When an agent(input_type:pcap, settings: {"iface":"default"}) is self-provisioned via a configuration file on port available with 3 agent tags and has status online. [Overwrite default: False. Paste only file: True]
        And pktvisor state is running
        And 1 Agent Group(s) is created with all tags contained in the agent
        And 3 simple policies same input_type as created via config file are applied to the group
    Then 3 dataset(s) have validity valid and 0 have validity invalid in 30 seconds
        And the container logs should contain the message "completed RPC subscription to group" within 30 seconds
        And this agent's heartbeat shows that 1 groups are matching the agent
        And this agent's heartbeat shows that 3 policies are applied and all has status running
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
        And referred sink must have active state on response within 30 seconds
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And remove the agent .yaml generated on each scenario


@smoke @config_file
Scenario: agent pcap with only agent tags subscription to a group with policies created before provision the agent (config file - auto_provision=true)
    Given the Orb user has a registered account
        And the Orb user logs in
        And that a sink already exists
        And 1 Agent Group(s) is created with 1 orb tag(s) (lower case)
        And 3 simple policies pcap are applied to the group
        And a new agent is created with 0 orb tag(s)
    When an agent(input_type:pcap, settings: {"iface":"default"}) is self-provisioned via a configuration file on port available with matching 1 group agent tags and has status online. [Overwrite default: False. Paste only file: False]
        And pktvisor state is running
    Then 3 dataset(s) have validity valid and 0 have validity invalid in 30 seconds
        And this agent's heartbeat shows that 1 groups are matching the agent
        And the container logs should contain the message "completed RPC subscription to group" within 30 seconds
        And this agent's heartbeat shows that 3 policies are applied and all has status running
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
        And referred sink must have active state on response within 30 seconds
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And remove the agent .yaml generated on each scenario


@smoke @config_file
Scenario: agent pcap with mixed tags subscription to a group with policies created after provision the agent (config file - auto_provision=true)
    Given the Orb user has a registered account
        And the Orb user logs in
        And that a sink already exists
    When an agent(input_type:pcap, settings: {"iface":"default"}) is self-provisioned via a configuration file on port available with 3 agent tags and has status online. [Overwrite default: False. Paste only file: True]
        And pktvisor state is running
        And edit the orb tags on agent and use 2 orb tag(s)
        And 1 Agent Group(s) is created with all tags contained in the agent
        And 3 simple policies same input_type as created via config file are applied to the group
    Then 3 dataset(s) have validity valid and 0 have validity invalid in 30 seconds
        And this agent's heartbeat shows that 1 groups are matching the agent
        And the container logs should contain the message "completed RPC subscription to group" within 30 seconds
        And this agent's heartbeat shows that 3 policies are applied and all has status running
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
        And referred sink must have active state on response within 30 seconds
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And remove the agent .yaml generated on each scenario


@smoke @config_file
Scenario: agent pcap with mixed tags subscription to a group with policies created before provision the agent (config file - auto_provision=true)
    Given the Orb user has a registered account
        And the Orb user logs in
        And that a sink already exists
        And 1 Agent Group(s) is created with 2 orb tag(s) (lower case)
        And 3 simple policies pcap are applied to the group
        And a new agent is created with 2 orb tag(s)
    When an agent(input_type:pcap, settings: {"iface":"default"}) is self-provisioned via a configuration file on port available with matching 1 group agent tags and has status online. [Overwrite default: False. Paste only file: False]
        And pktvisor state is running
    Then 3 dataset(s) have validity valid and 0 have validity invalid in 30 seconds
        And this agent's heartbeat shows that 1 groups are matching the agent
        And the container logs should contain the message "completed RPC subscription to group" within 30 seconds
        And this agent's heartbeat shows that 3 policies are applied and all has status running
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
        And referred sink must have active state on response within 30 seconds
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And remove the agent .yaml generated on each scenario

@smoke @config_file
Scenario: agent pcap with only agent tags subscription to a group with policies created after provision the agent (config file - auto_provision=false)
    Given the Orb user has a registered account
        And the Orb user logs in
        And that a sink already exists
        And a new agent is created with 0 orb tag(s)
    When an agent(input_type:pcap, settings: {"iface":"default"}) is provisioned via a configuration file on port available with 3 agent tags and has status online. [Overwrite default: False. Paste only file: True]
        And pktvisor state is running
        And 1 Agent Group(s) is created with all tags contained in the agent
        And 3 simple policies same input_type as created via config file are applied to the group
    Then 3 dataset(s) have validity valid and 0 have validity invalid in 30 seconds
        And the container logs should contain the message "completed RPC subscription to group" within 30 seconds
        And this agent's heartbeat shows that 1 groups are matching the agent
        And this agent's heartbeat shows that 3 policies are applied and all has status running
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
        And referred sink must have active state on response within 30 seconds
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And remove the agent .yaml generated on each scenario

#@smoke @config_file
@MUTE
Scenario: agent pcap with only agent tags subscription to a group with policies created before provision the agent (config file - auto_provision=false)
    Given the Orb user has a registered account
        And the Orb user logs in
        And that a sink already exists
        And 1 Agent Group(s) is created with 1 orb tag(s) (lower case)
        And 3 simple policies pcap are applied to the group
        And a new agent is created with 0 orb tag(s)
    When an agent(input_type:pcap, settings: {"iface":"default"}) is provisioned via a configuration file on port available with matching 1 group agent tags and has status online. [Overwrite default: False. Paste only file: False]
        And pktvisor state is running
    Then 3 dataset(s) have validity valid and 0 have validity invalid in 30 seconds
        And this agent's heartbeat shows that 1 groups are matching the agent
        And the container logs should contain the message "completed RPC subscription to group" within 30 seconds
        And this agent's heartbeat shows that 3 policies are applied and all has status running
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
        And referred sink must have active state on response within 30 seconds
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And remove the agent .yaml generated on each scenario

@smoke @config_file
Scenario: agent pcap with mixed tags subscription to a group with policies created after provision the agent (config file - auto_provision=false)
    Given the Orb user has a registered account
        And the Orb user logs in
        And that a sink already exists
        And a new agent is created with 2 orb tag(s)
    When an agent(input_type:pcap, settings: {"iface":"default"}) is provisioned via a configuration file on port available with 3 agent tags and has status online. [Overwrite default: False. Paste only file: False]
        And pktvisor state is running
        And edit the orb tags on agent and use 2 orb tag(s)
        And 1 Agent Group(s) is created with all tags contained in the agent
        And 3 simple policies same input_type as created via config file are applied to the group
    Then 3 dataset(s) have validity valid and 0 have validity invalid in 30 seconds
        And this agent's heartbeat shows that 1 groups are matching the agent
        And the container logs should contain the message "completed RPC subscription to group" within 30 seconds
        And this agent's heartbeat shows that 3 policies are applied and all has status running
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
        And referred sink must have active state on response within 30 seconds
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And remove the agent .yaml generated on each scenario

#@smoke @config_file
@MUTE
Scenario: agent pcap with mixed tags subscription to a group with policies created before provision the agent (config file - auto_provision=false)
    Given the Orb user has a registered account
        And the Orb user logs in
        And that a sink already exists
        And 1 Agent Group(s) is created with 2 orb tag(s) (lower case)
        And 3 simple policies pcap are applied to the group
        And a new agent is created with 2 orb tag(s)
    When an agent(input_type:pcap, settings: {"iface":"default"}) is provisioned via a configuration file on port available with matching 1 group agent tags and has status online. [Overwrite default: False. Paste only file: False]
        And pktvisor state is running
    Then 3 dataset(s) have validity valid and 0 have validity invalid in 30 seconds
        And this agent's heartbeat shows that 1 groups are matching the agent
        And the container logs should contain the message "completed RPC subscription to group" within 30 seconds
        And this agent's heartbeat shows that 3 policies are applied and all has status running
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
        And referred sink must have active state on response within 30 seconds
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And remove the agent .yaml generated on each scenario


########### flow

@smoke @config_file
Scenario: agent flow with only agent tags subscription to a group with policies created after provision the agent (config file - auto_provision=true)
    Given the Orb user has a registered account
        And the Orb user logs in
        And that a sink already exists
    When an agent(input_type:flow, settings: {"bind":"0.0.0.0", "port":"available_port"}) is self-provisioned via a configuration file on port available with 3 agent tags and has status online. [Overwrite default: True. Paste only file: False]
        And pktvisor state is running
        And 1 Agent Group(s) is created with all tags contained in the agent
        And 3 simple policies same input_type as created via config file are applied to the group
    Then 3 dataset(s) have validity valid and 0 have validity invalid in 30 seconds
        And the container logs should contain the message "completed RPC subscription to group" within 30 seconds
        And this agent's heartbeat shows that 1 groups are matching the agent
        And this agent's heartbeat shows that 3 policies are applied and all has status running
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
        And referred sink must have active state on response within 30 seconds
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And remove the agent .yaml generated on each scenario


@smoke @config_file
Scenario: agent flow with only agent tags subscription to a group with policies created before provision the agent (config file - auto_provision=true)
    Given the Orb user has a registered account
        And the Orb user logs in
        And that a sink already exists
        And 1 Agent Group(s) is created with 1 orb tag(s) (lower case)
        And 3 simple policies flow are applied to the group
        And a new agent is created with 0 orb tag(s)
    When an agent(input_type:flow, settings: {"bind":"0.0.0.0", "port":"available_port"}) is self-provisioned via a configuration file on port available with matching 1 group agent tags and has status online. [Overwrite default: False. Paste only file: False]
        And pktvisor state is running
    Then 3 dataset(s) have validity valid and 0 have validity invalid in 30 seconds
        And this agent's heartbeat shows that 1 groups are matching the agent
        And the container logs should contain the message "completed RPC subscription to group" within 30 seconds
        And this agent's heartbeat shows that 3 policies are applied and all has status running
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
        And referred sink must have active state on response within 30 seconds
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And remove the agent .yaml generated on each scenario


@smoke @config_file
Scenario: agent flow with mixed tags subscription to a group with policies created after provision the agent (config file - auto_provision=true)
    Given the Orb user has a registered account
        And the Orb user logs in
        And that a sink already exists
    When an agent(input_type:flow, settings: {"bind":"0.0.0.0", "port":"available_port"}) is self-provisioned via a configuration file on port available with 3 agent tags and has status online. [Overwrite default: False. Paste only file: False]
        And pktvisor state is running
        And edit the orb tags on agent and use 2 orb tag(s)
        And 1 Agent Group(s) is created with all tags contained in the agent
        And 3 simple policies same input_type as created via config file are applied to the group
    Then 3 dataset(s) have validity valid and 0 have validity invalid in 30 seconds
        And this agent's heartbeat shows that 1 groups are matching the agent
        And the container logs should contain the message "completed RPC subscription to group" within 30 seconds
        And this agent's heartbeat shows that 3 policies are applied and all has status running
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
        And referred sink must have active state on response within 30 seconds
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And remove the agent .yaml generated on each scenario


@smoke @config_file
Scenario: agent flow with mixed tags subscription to a group with policies created before provision the agent (config file - auto_provision=true)
    Given the Orb user has a registered account
        And the Orb user logs in
        And that a sink already exists
        And 1 Agent Group(s) is created with 2 orb tag(s) (lower case)
        And 3 simple policies flow are applied to the group
        And a new agent is created with 2 orb tag(s)
    When an agent(input_type:flow, settings: {"bind":"0.0.0.0", "port":"available_port"}) is self-provisioned via a configuration file on port available with matching 1 group agent tags and has status online. [Overwrite default: False. Paste only file: False]
        And pktvisor state is running
    Then 3 dataset(s) have validity valid and 0 have validity invalid in 30 seconds
        And this agent's heartbeat shows that 1 groups are matching the agent
        And the container logs should contain the message "completed RPC subscription to group" within 30 seconds
        And this agent's heartbeat shows that 3 policies are applied and all has status running
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
        And referred sink must have active state on response within 30 seconds
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And remove the agent .yaml generated on each scenario

@smoke @config_file
Scenario: agent flow with only agent tags subscription to a group with policies created after provision the agent (config file - auto_provision=false)
    Given the Orb user has a registered account
        And the Orb user logs in
        And that a sink already exists
        And a new agent is created with 0 orb tag(s)
    When an agent(input_type:flow, settings: {"bind":"0.0.0.0", "port":"available_port"}) is provisioned via a configuration file on port available with 3 agent tags and has status online. [Overwrite default: False. Paste only file: False]
        And pktvisor state is running
        And 1 Agent Group(s) is created with all tags contained in the agent
        And 3 simple policies same input_type as created via config file are applied to the group
    Then 3 dataset(s) have validity valid and 0 have validity invalid in 30 seconds
        And the container logs should contain the message "completed RPC subscription to group" within 30 seconds
        And this agent's heartbeat shows that 1 groups are matching the agent
        And this agent's heartbeat shows that 3 policies are applied and all has status running
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
        And referred sink must have active state on response within 30 seconds
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And remove the agent .yaml generated on each scenario

#@smoke @config_file
@MUTE
Scenario: agent flow with only agent tags subscription to a group with policies created before provision the agent (config file - auto_provision=false)
    Given the Orb user has a registered account
        And the Orb user logs in
        And that a sink already exists
        And 1 Agent Group(s) is created with 1 orb tag(s) (lower case)
        And 3 simple policies flow are applied to the group
        And a new agent is created with 0 orb tag(s)
    When an agent(input_type:flow, settings: {"bind":"0.0.0.0", "port":"available_port"}) is provisioned via a configuration file on port available with matching 1 group agent tags and has status online. [Overwrite default: False. Paste only file: False]
        And pktvisor state is running
    Then 3 dataset(s) have validity valid and 0 have validity invalid in 30 seconds
        And this agent's heartbeat shows that 1 groups are matching the agent
        And the container logs should contain the message "completed RPC subscription to group" within 30 seconds
        And this agent's heartbeat shows that 3 policies are applied and all has status running
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
        And referred sink must have active state on response within 30 seconds
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And remove the agent .yaml generated on each scenario

@smoke @config_file
Scenario: agent flow with mixed tags subscription to a group with policies created after provision the agent (config file - auto_provision=false)
    Given the Orb user has a registered account
        And the Orb user logs in
        And that a sink already exists
        And a new agent is created with 2 orb tag(s)
    When an agent(input_type:flow, settings: {"bind":"0.0.0.0", "port":"available_port"}) is provisioned via a configuration file on port available with 3 agent tags and has status online. [Overwrite default: False. Paste only file: False]
        And pktvisor state is running
        And edit the orb tags on agent and use 2 orb tag(s)
        And 1 Agent Group(s) is created with all tags contained in the agent
        And 3 simple policies same input_type as created via config file are applied to the group
    Then 3 dataset(s) have validity valid and 0 have validity invalid in 30 seconds
        And this agent's heartbeat shows that 1 groups are matching the agent
        And the container logs should contain the message "completed RPC subscription to group" within 30 seconds
        And this agent's heartbeat shows that 3 policies are applied and all has status running
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
        And referred sink must have active state on response within 30 seconds
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And remove the agent .yaml generated on each scenario

#@smoke @config_file
@MUTE
Scenario: agent flow with mixed tags subscription to a group with policies created before provision the agent (config file - auto_provision=false)
    Given the Orb user has a registered account
        And the Orb user logs in
        And that a sink already exists
        And 1 Agent Group(s) is created with 2 orb tag(s) (lower case)
        And 3 simple policies flow are applied to the group
        And a new agent is created with 2 orb tag(s)
    When an agent(input_type:flow, settings: {"bind":"0.0.0.0", "port":"available_port"}) is provisioned via a configuration file on port available with matching 1 group agent tags and has status online. [Overwrite default: False. Paste only file: False]
        And pktvisor state is running
    Then 3 dataset(s) have validity valid and 0 have validity invalid in 30 seconds
        And this agent's heartbeat shows that 1 groups are matching the agent
        And the container logs should contain the message "completed RPC subscription to group" within 30 seconds
        And this agent's heartbeat shows that 3 policies are applied and all has status running
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
        And referred sink must have active state on response within 30 seconds
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And remove the agent .yaml generated on each scenario



########### dnstap

@smoke @config_file
Scenario: agent dnstap with only agent tags subscription to a group with policies created after provision the agent (config file - auto_provision=true)
    Given the Orb user has a registered account
        And the Orb user logs in
        And that a sink already exists
    When an agent(input_type:dnstap, settings: {"tcp":"0.0.0.0:available_port", "only_hosts":"0.0.0.0/32"}) is self-provisioned via a configuration file on port available with 3 agent tags and has status online. [Overwrite default: False. Paste only file: False]
        And pktvisor state is running
        And 1 Agent Group(s) is created with all tags contained in the agent
        And 3 simple policies same input_type as created via config file are applied to the group
    Then 3 dataset(s) have validity valid and 0 have validity invalid in 30 seconds
        And the container logs should contain the message "completed RPC subscription to group" within 30 seconds
        And this agent's heartbeat shows that 1 groups are matching the agent
        And this agent's heartbeat shows that 3 policies are applied and all has status running
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
        And referred sink must have active state on response within 30 seconds
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And remove the agent .yaml generated on each scenario


@smoke @config_file
Scenario: agent dnstap with only agent tags subscription to a group with policies created before provision the agent (config file - auto_provision=true)
    Given the Orb user has a registered account
        And the Orb user logs in
        And that a sink already exists
        And 1 Agent Group(s) is created with 1 orb tag(s) (lower case)
        And 3 simple policies dnstap are applied to the group
        And a new agent is created with 0 orb tag(s)
    When an agent(input_type:dnstap, settings: {"tcp":"0.0.0.0:available_port", "only_hosts":"0.0.0.0/32"}) is self-provisioned via a configuration file on port available with matching 1 group agent tags and has status online. [Overwrite default: False. Paste only file: False]
        And pktvisor state is running
    Then 3 dataset(s) have validity valid and 0 have validity invalid in 30 seconds
        And this agent's heartbeat shows that 1 groups are matching the agent
        And the container logs should contain the message "completed RPC subscription to group" within 30 seconds
        And this agent's heartbeat shows that 3 policies are applied and all has status running
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
        And referred sink must have active state on response within 30 seconds
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And remove the agent .yaml generated on each scenario


@smoke @config_file
Scenario: agent dnstap with mixed tags subscription to a group with policies created after provision the agent (config file - auto_provision=true)
    Given the Orb user has a registered account
        And the Orb user logs in
        And that a sink already exists
    When an agent(input_type:dnstap, settings: {"tcp":"0.0.0.0:available_port", "only_hosts":"0.0.0.0/32"}) is self-provisioned via a configuration file on port available with 3 agent tags and has status online. [Overwrite default: False. Paste only file: False]
        And pktvisor state is running
        And edit the orb tags on agent and use 2 orb tag(s)
        And 1 Agent Group(s) is created with all tags contained in the agent
        And 3 simple policies same input_type as created via config file are applied to the group
    Then 3 dataset(s) have validity valid and 0 have validity invalid in 30 seconds
        And this agent's heartbeat shows that 1 groups are matching the agent
        And the container logs should contain the message "completed RPC subscription to group" within 30 seconds
        And this agent's heartbeat shows that 3 policies are applied and all has status running
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
        And referred sink must have active state on response within 30 seconds
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And remove the agent .yaml generated on each scenario


@smoke @config_file
Scenario: agent dnstap with mixed tags subscription to a group with policies created before provision the agent (config file - auto_provision=true)
    Given the Orb user has a registered account
        And the Orb user logs in
        And that a sink already exists
        And 1 Agent Group(s) is created with 2 orb tag(s) (lower case)
        And 3 simple policies dnstap are applied to the group
        And a new agent is created with 2 orb tag(s)
    When an agent(input_type:dnstap, settings: {"tcp":"0.0.0.0:available_port", "only_hosts":"0.0.0.0/32"}) is self-provisioned via a configuration file on port available with matching 1 group agent tags and has status online. [Overwrite default: False. Paste only file: False]
        And pktvisor state is running
    Then 3 dataset(s) have validity valid and 0 have validity invalid in 30 seconds
        And this agent's heartbeat shows that 1 groups are matching the agent
        And the container logs should contain the message "completed RPC subscription to group" within 30 seconds
        And this agent's heartbeat shows that 3 policies are applied and all has status running
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
        And referred sink must have active state on response within 30 seconds
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And remove the agent .yaml generated on each scenario

@smoke @config_file
Scenario: agent dnstap with only agent tags subscription to a group with policies created after provision the agent (config file - auto_provision=false)
    Given the Orb user has a registered account
        And the Orb user logs in
        And that a sink already exists
        And a new agent is created with 0 orb tag(s)
    When an agent(input_type:dnstap, settings: {"tcp":"0.0.0.0:available_port", "only_hosts":"0.0.0.0/32"}) is provisioned via a configuration file on port available with 3 agent tags and has status online. [Overwrite default: False. Paste only file: False]
        And pktvisor state is running
        And 1 Agent Group(s) is created with all tags contained in the agent
        And 3 simple policies same input_type as created via config file are applied to the group
    Then 3 dataset(s) have validity valid and 0 have validity invalid in 30 seconds
        And the container logs should contain the message "completed RPC subscription to group" within 30 seconds
        And this agent's heartbeat shows that 1 groups are matching the agent
        And this agent's heartbeat shows that 3 policies are applied and all has status running
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
        And referred sink must have active state on response within 30 seconds
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And remove the agent .yaml generated on each scenario

#@smoke @config_file
@MUTE
Scenario: agent dnstap with only agent tags subscription to a group with policies created before provision the agent (config file - auto_provision=false)
    Given the Orb user has a registered account
        And the Orb user logs in
        And that a sink already exists
        And 1 Agent Group(s) is created with 1 orb tag(s) (lower case)
        And 3 simple policies dnstap are applied to the group
        And a new agent is created with 0 orb tag(s)
    When an agent(input_type:dnstap, settings: {{"tcp":"0.0.0.0:available_port", "only_hosts":"0.0.0.0/32"}) is provisioned via a configuration file on port available with matching 1 group agent tags and has status online. [Overwrite default: False. Paste only file: False]
        And pktvisor state is running
    Then 3 dataset(s) have validity valid and 0 have validity invalid in 30 seconds
        And this agent's heartbeat shows that 1 groups are matching the agent
        And the container logs should contain the message "completed RPC subscription to group" within 30 seconds
        And this agent's heartbeat shows that 3 policies are applied and all has status running
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
        And referred sink must have active state on response within 30 seconds
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And remove the agent .yaml generated on each scenario

@smoke @config_file
Scenario: agent dnstap with mixed tags subscription to a group with policies created after provision the agent (config file - auto_provision=false)
    Given the Orb user has a registered account
        And the Orb user logs in
        And that a sink already exists
        And a new agent is created with 2 orb tag(s)
    When an agent(input_type:dnstap, settings: {"tcp":"0.0.0.0:available_port", "only_hosts":"0.0.0.0/32"}) is provisioned via a configuration file on port available with 3 agent tags and has status online. [Overwrite default: False. Paste only file: False]
        And pktvisor state is running
        And edit the orb tags on agent and use 2 orb tag(s)
        And 1 Agent Group(s) is created with all tags contained in the agent
        And 3 simple policies same input_type as created via config file are applied to the group
    Then 3 dataset(s) have validity valid and 0 have validity invalid in 30 seconds
        And this agent's heartbeat shows that 1 groups are matching the agent
        And the container logs should contain the message "completed RPC subscription to group" within 30 seconds
        And this agent's heartbeat shows that 3 policies are applied and all has status running
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
        And referred sink must have active state on response within 30 seconds
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And remove the agent .yaml generated on each scenario

#@smoke @config_file
@MUTE
Scenario: agent dnstap with mixed tags subscription to a group with policies created before provision the agent (config file - auto_provision=false)
    Given the Orb user has a registered account
        And the Orb user logs in
        And that a sink already exists
        And 1 Agent Group(s) is created with 2 orb tag(s) (lower case)
        And 3 simple policies dnstap are applied to the group
        And a new agent is created with 2 orb tag(s)
    When an agent(input_type:dnstap, settings: {"tcp":"0.0.0.0:available_port", "only_hosts":"0.0.0.0/32"}) is provisioned via a configuration file on port available with matching 1 group agent tags and has status online. [Overwrite default: False. Paste only file: False]
        And pktvisor state is running
    Then 3 dataset(s) have validity valid and 0 have validity invalid in 30 seconds
        And this agent's heartbeat shows that 1 groups are matching the agent
        And the container logs should contain the message "completed RPC subscription to group" within 30 seconds
        And this agent's heartbeat shows that 3 policies are applied and all has status running
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
        And referred sink must have active state on response within 30 seconds
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And remove the agent .yaml generated on each scenario


########### netprobe


@smoke @config_file @netprobe
Scenario: agent netprobe with only agent tags subscription to a group with policies created after provision the agent (config file - auto_provision=true)
    Given the Orb user has a registered account
        And the Orb user logs in
        And that a sink already exists
    When an agent(input_type:netprobe, settings: {"test_type":"ping", "packets_per_test":3, "interval_msec":3000, "timeout_msec":1500, "packets_interval_msec":50, "packet_payload_size":56, "targets": {"www.google.com": {"target": "www.google.com"}, "orb_community": {"target": "orb.community"}}}) is self-provisioned via a configuration file on port available with 3 agent tags and has status online. [Overwrite default: False. Paste only file: True]
        And pktvisor state is running
        And 1 Agent Group(s) is created with all tags contained in the agent
        And 3 simple policies same input_type as created via config file are applied to the group
    Then 3 dataset(s) have validity valid and 0 have validity invalid in 30 seconds
        And the container logs should contain the message "completed RPC subscription to group" within 30 seconds
        And this agent's heartbeat shows that 1 groups are matching the agent
        And this agent's heartbeat shows that 3 policies are applied and all has status running
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
        And referred sink must have active state on response within 30 seconds
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And remove the agent .yaml generated on each scenario


@smoke @config_file @netprobe
Scenario: agent netprobe with only agent tags subscription to a group with policies created before provision the agent (config file - auto_provision=true)
    Given the Orb user has a registered account
        And the Orb user logs in
        And that a sink already exists
        And 1 Agent Group(s) is created with 1 orb tag(s) (lower case)
        And 3 advanced policies netprobe are applied to the group
        And a new agent is created with 0 orb tag(s)
    When an agent(input_type:netprobe, settings: {"test_type":"ping", "packets_per_test":3, "interval_msec":3000, "timeout_msec":1500, "packets_interval_msec":50, "packet_payload_size":56, "targets": {"www.google.com": {"target": "www.google.com"}, "orb_community": {"target": "orb.community"}}}) is self-provisioned via a configuration file on port available with matching 1 group agent tags and has status online. [Overwrite default: False. Paste only file: False]
        And pktvisor state is running
    Then 3 dataset(s) have validity valid and 0 have validity invalid in 30 seconds
        And this agent's heartbeat shows that 1 groups are matching the agent
        And the container logs should contain the message "completed RPC subscription to group" within 30 seconds
        And this agent's heartbeat shows that 3 policies are applied and all has status running
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
        And referred sink must have active state on response within 30 seconds
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And remove the agent .yaml generated on each scenario


@smoke @config_file @netprobe
Scenario: agent netprobe with mixed tags subscription to a group with policies created after provision the agent (config file - auto_provision=true)
    Given the Orb user has a registered account
        And the Orb user logs in
        And that a sink already exists
    When an agent(input_type:netprobe, settings: {"test_type":"ping", "packets_per_test":3, "interval_msec":3000, "timeout_msec":1500, "packets_interval_msec":50, "packet_payload_size":56, "targets": {"www.google.com": {"target": "www.google.com"}, "orb_community": {"target": "orb.community"}}}) is self-provisioned via a configuration file on port available with 3 agent tags and has status online. [Overwrite default: False. Paste only file: False]
        And pktvisor state is running
        And edit the orb tags on agent and use 2 orb tag(s)
        And 1 Agent Group(s) is created with all tags contained in the agent
        And 3 simple policies same input_type as created via config file are applied to the group
    Then 3 dataset(s) have validity valid and 0 have validity invalid in 30 seconds
        And this agent's heartbeat shows that 1 groups are matching the agent
        And the container logs should contain the message "completed RPC subscription to group" within 30 seconds
        And this agent's heartbeat shows that 3 policies are applied and all has status running
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
        And referred sink must have active state on response within 30 seconds
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And remove the agent .yaml generated on each scenario


@smoke @config_file @netprobe
Scenario: agent netprobe with mixed tags subscription to a group with policies created before provision the agent (config file - auto_provision=true)
    Given the Orb user has a registered account
        And the Orb user logs in
        And that a sink already exists
        And 1 Agent Group(s) is created with 2 orb tag(s) (lower case)
        And 3 simple policies netprobe are applied to the group
        And a new agent is created with 2 orb tag(s)
    When an agent(input_type:netprobe, settings: {"test_type":"ping", "packets_per_test":3, "interval_msec":3000, "timeout_msec":1500, "packets_interval_msec":50, "packet_payload_size":56, "targets": {"www.google.com": {"target": "www.google.com"}, "orb_community": {"target": "orb.community"}}}) is self-provisioned via a configuration file on port available with matching 1 group agent tags and has status online. [Overwrite default: False. Paste only file: False]
        And pktvisor state is running
    Then 3 dataset(s) have validity valid and 0 have validity invalid in 30 seconds
        And this agent's heartbeat shows that 1 groups are matching the agent
        And the container logs should contain the message "completed RPC subscription to group" within 30 seconds
        And this agent's heartbeat shows that 3 policies are applied and all has status running
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
        And referred sink must have active state on response within 30 seconds
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And remove the agent .yaml generated on each scenario


@smoke @config_file @netprobe
Scenario: agent netprobe with only agent tags subscription to a group with policies created after provision the agent (config file - auto_provision=false)
    Given the Orb user has a registered account
        And the Orb user logs in
        And that a sink already exists
        And a new agent is created with 0 orb tag(s)
    When an agent(input_type:netprobe, settings: {"test_type":"ping", "packets_per_test":3, "interval_msec":3000, "timeout_msec":1500, "packets_interval_msec":50, "packet_payload_size":56, "targets": {"www.google.com": {"target": "www.google.com"}, "orb_community": {"target": "orb.community"}}}) is provisioned via a configuration file on port available with 3 agent tags and has status online. [Overwrite default: True. Paste only file: True]
        And pktvisor state is running
        And 1 Agent Group(s) is created with all tags contained in the agent
        And 3 simple policies same input_type as created via config file are applied to the group
    Then 3 dataset(s) have validity valid and 0 have validity invalid in 30 seconds
        And the container logs should contain the message "completed RPC subscription to group" within 30 seconds
        And this agent's heartbeat shows that 1 groups are matching the agent
        And this agent's heartbeat shows that 3 policies are applied and all has status running
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
        And referred sink must have active state on response within 30 seconds
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And remove the agent .yaml generated on each scenario


#@smoke @config_file @netprobe
@MUTE
Scenario: agent netprobe with only agent tags subscription to a group with policies created before provision the agent (config file - auto_provision=false)
    Given the Orb user has a registered account
        And the Orb user logs in
        And that a sink already exists
        And 1 Agent Group(s) is created with 1 orb tag(s) (lower case)
        And 3 simple policies netprobe are applied to the group
        And a new agent is created with 0 orb tag(s)
    When an agent(input_type:netprobe, settings: {"test_type":"ping", "packets_per_test":3, "interval_msec":3000, "timeout_msec":1500, "packets_interval_msec":50, "packet_payload_size":56, "targets": {"www.google.com": {"target": "www.google.com"}, "orb_community": {"target": "orb.community"}}}) is provisioned via a configuration file on port available with matching 1 group agent tags and has status online. [Overwrite default: False. Paste only file: False]
        And pktvisor state is running
    Then 3 dataset(s) have validity valid and 0 have validity invalid in 30 seconds
        And this agent's heartbeat shows that 1 groups are matching the agent
        And the container logs should contain the message "completed RPC subscription to group" within 30 seconds
        And this agent's heartbeat shows that 3 policies are applied and all has status running
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
        And referred sink must have active state on response within 30 seconds
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And remove the agent .yaml generated on each scenario


@smoke @config_file @netprobe
Scenario: agent netprobe with mixed tags subscription to a group with policies created after provision the agent (config file - auto_provision=false)
    Given the Orb user has a registered account
        And the Orb user logs in
        And that a sink already exists
        And a new agent is created with 2 orb tag(s)
    When an agent(input_type:netprobe, settings: {"test_type":"ping", "packets_per_test":3, "interval_msec":3000, "timeout_msec":1500, "packets_interval_msec":50, "packet_payload_size":56, "targets": {"www.google.com": {"target": "www.google.com"}, "orb_community": {"target": "orb.community"}}}) is provisioned via a configuration file on port available with 3 agent tags and has status online. [Overwrite default: False. Paste only file: False]
        And pktvisor state is running
        And edit the orb tags on agent and use 2 orb tag(s)
        And 1 Agent Group(s) is created with all tags contained in the agent
        And 3 simple policies same input_type as created via config file are applied to the group
    Then 3 dataset(s) have validity valid and 0 have validity invalid in 30 seconds
        And this agent's heartbeat shows that 1 groups are matching the agent
        And the container logs should contain the message "completed RPC subscription to group" within 30 seconds
        And this agent's heartbeat shows that 3 policies are applied and all has status running
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
        And referred sink must have active state on response within 30 seconds
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And remove the agent .yaml generated on each scenario


#@smoke @config_file @netprobe
@MUTE
Scenario: agent netprobe with mixed tags subscription to a group with policies created before provision the agent (config file - auto_provision=false)
    Given the Orb user has a registered account
        And the Orb user logs in
        And that a sink already exists
        And 1 Agent Group(s) is created with 2 orb tag(s) (lower case)
        And 3 simple policies netprobe are applied to the group
        And a new agent is created with 2 orb tag(s)
    When an agent(input_type:netprobe, settings: {"test_type":"ping", "packets_per_test":3, "interval_msec":3000, "timeout_msec":1500, "packets_interval_msec":50, "packet_payload_size":56, "targets": {"www.google.com": {"target": "www.google.com"}, "orb_community": {"target": "orb.community"}}}) is provisioned via a configuration file on port available with matching 1 group agent tags and has status online. [Overwrite default: False. Paste only file: False]
        And pktvisor state is running
    Then 3 dataset(s) have validity valid and 0 have validity invalid in 30 seconds
        And this agent's heartbeat shows that 1 groups are matching the agent
        And the container logs should contain the message "completed RPC subscription to group" within 30 seconds
        And this agent's heartbeat shows that 3 policies are applied and all has status running
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
        And referred sink must have active state on response within 30 seconds
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And remove the agent .yaml generated on each scenario
