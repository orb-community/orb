@metrics @AUTORETRY
Feature: Integration tests validating metric groups

#### netprobe

@smoke @metric_groups @metrics_netprobe
Scenario: netprobe handler with default metric groups configuration
    Given the Orb user has a registered account
        And the Orb user logs in
        And that a sink already exists
    When an agent(input_type:netprobe, settings: {"test_type":"ping", "targets": {"www.google.com": {"target": "www.google.com"}, "orb_community": {"target": "orb.community"}}}) is self-provisioned via a configuration file on port available with 1 agent tags and has status online. [Overwrite default: False. Paste only file: True]
        And pktvisor state is running
        And 1 Agent Group(s) is created with all tags contained in the agent
        And a netprobe policy netprobe with tap_selector matching all tag(s) of the tap from an agent, default metric_groups enabled, default metric_groups disabled and settings: default is applied to the group
    Then 1 dataset(s) have validity valid and 0 have validity invalid in 30 seconds
        And the container logs should contain the message "completed RPC subscription to group" within 30 seconds
        And this agent's heartbeat shows that 1 groups are matching the agent
        And this agent's heartbeat shows that 1 policies are applied and all has status running
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
        And referred sink must have active state on response within 30 seconds
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And remove the agent .yaml generated on each scenario


@smoke @metric_groups @metrics_netprobe
Scenario: netprobe handler with all metric groups enabled
    Given the Orb user has a registered account
        And the Orb user logs in
        And that a sink already exists
    When an agent(input_type:netprobe, settings: {"test_type":"ping", "targets": {"www.google.com": {"target": "www.google.com"}, "orb_community": {"target": "orb.community"}}}) is self-provisioned via a configuration file on port available with 1 agent tags and has status online. [Overwrite default: False. Paste only file: True]
        And pktvisor state is running
        And 1 Agent Group(s) is created with all tags contained in the agent
        And a netprobe policy netprobe with tap_selector matching all tag(s) of the tap from an agent, all metric_groups enabled, none metric_groups disabled and settings: default is applied to the group
    Then 1 dataset(s) have validity valid and 0 have validity invalid in 30 seconds
        And the container logs should contain the message "completed RPC subscription to group" within 30 seconds
        And this agent's heartbeat shows that 1 groups are matching the agent
        And this agent's heartbeat shows that 1 policies are applied and all has status running
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
        And referred sink must have active state on response within 30 seconds
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And remove the agent .yaml generated on each scenario


@smoke @metric_groups @metrics_netprobe
Scenario: netprobe handler with all metric groups disabled
    Given the Orb user has a registered account
        And the Orb user logs in
        And that a sink already exists
    When an agent(input_type:netprobe, settings: {"test_type":"ping", "targets": {"www.google.com": {"target": "www.google.com"}, "orb_community": {"target": "orb.community"}}}) is self-provisioned via a configuration file on port available with 1 agent tags and has status online. [Overwrite default: False. Paste only file: True]
        And pktvisor state is running
        And 1 Agent Group(s) is created with all tags contained in the agent
        And a netprobe policy netprobe with tap_selector matching all tag(s) of the tap from an agent, none metric_groups enabled, all metric_groups disabled and settings: default is applied to the group
    Then 1 dataset(s) have validity valid and 0 have validity invalid in 30 seconds
        And the container logs should contain the message "completed RPC subscription to group" within 30 seconds
        And this agent's heartbeat shows that 1 groups are matching the agent
        And this agent's heartbeat shows that 1 policies are applied and all has status running
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
        And referred sink must have active state on response within 30 seconds
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And remove the agent .yaml generated on each scenario


@smoke @metric_groups @metrics_netprobe
Scenario: netprobe handler with only counters metric groups enabled
    Given the Orb user has a registered account
        And the Orb user logs in
        And that a sink already exists
    When an agent(input_type:netprobe, settings: {"test_type":"ping", "targets": {"www.google.com": {"target": "www.google.com"}, "orb_community": {"target": "orb.community"}}}) is self-provisioned via a configuration file on port available with 1 agent tags and has status online. [Overwrite default: False. Paste only file: True]
        And pktvisor state is running
        And 1 Agent Group(s) is created with all tags contained in the agent
        And a netprobe policy netprobe with tap_selector matching all tag(s) of the tap from an agent, counters metric_groups enabled, quantiles, histograms metric_groups disabled and settings: default is applied to the group
    Then 1 dataset(s) have validity valid and 0 have validity invalid in 30 seconds
        And the container logs should contain the message "completed RPC subscription to group" within 30 seconds
        And this agent's heartbeat shows that 1 groups are matching the agent
        And this agent's heartbeat shows that 1 policies are applied and all has status running
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
        And referred sink must have active state on response within 30 seconds
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And remove the agent .yaml generated on each scenario


@smoke @metric_groups @metrics_netprobe
Scenario: netprobe handler with only quantiles metric groups enabled
    Given the Orb user has a registered account
        And the Orb user logs in
        And that a sink already exists
    When an agent(input_type:netprobe, settings: {"test_type":"ping", "targets": {"www.google.com": {"target": "www.google.com"}, "orb_community": {"target": "orb.community"}}}) is self-provisioned via a configuration file on port available with 1 agent tags and has status online. [Overwrite default: False. Paste only file: True]
        And pktvisor state is running
        And 1 Agent Group(s) is created with all tags contained in the agent
        And a netprobe policy netprobe with tap_selector matching all tag(s) of the tap from an agent, quantiles metric_groups enabled, counters, histograms metric_groups disabled and settings: default is applied to the group
    Then 1 dataset(s) have validity valid and 0 have validity invalid in 30 seconds
        And the container logs should contain the message "completed RPC subscription to group" within 30 seconds
        And this agent's heartbeat shows that 1 groups are matching the agent
        And this agent's heartbeat shows that 1 policies are applied and all has status running
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
        And referred sink must have active state on response within 30 seconds
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And remove the agent .yaml generated on each scenario


@smoke @metric_groups @metrics_netprobe
Scenario: netprobe handler with only histograms metric groups enabled
    Given the Orb user has a registered account
        And the Orb user logs in
        And that a sink already exists
    When an agent(input_type:netprobe, settings: {"test_type":"ping", "targets": {"www.google.com": {"target": "www.google.com"}, "orb_community": {"target": "orb.community"}}}) is self-provisioned via a configuration file on port available with 1 agent tags and has status online. [Overwrite default: False. Paste only file: True]
        And pktvisor state is running
        And 1 Agent Group(s) is created with all tags contained in the agent
        And a netprobe policy netprobe with tap_selector matching all tag(s) of the tap from an agent, histograms metric_groups enabled, quantiles, counters metric_groups disabled and settings: default is applied to the group
    Then 1 dataset(s) have validity valid and 0 have validity invalid in 30 seconds
        And the container logs should contain the message "completed RPC subscription to group" within 30 seconds
        And this agent's heartbeat shows that 1 groups are matching the agent
        And this agent's heartbeat shows that 1 policies are applied and all has status running
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
        And referred sink must have active state on response within 30 seconds
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And remove the agent .yaml generated on each scenario


@smoke @metric_groups @metrics_netprobe
Scenario: netprobe handler with counters and histograms metric groups enabled and quantiles disabled
    Given the Orb user has a registered account
        And the Orb user logs in
        And that a sink already exists
    When an agent(input_type:netprobe, settings: {"test_type":"ping", "targets": {"www.google.com": {"target": "www.google.com"}, "orb_community": {"target": "orb.community"}}}) is self-provisioned via a configuration file on port available with 1 agent tags and has status online. [Overwrite default: False. Paste only file: True]
        And pktvisor state is running
        And 1 Agent Group(s) is created with all tags contained in the agent
        And a netprobe policy netprobe with tap_selector matching all tag(s) of the tap from an agent, counters, histograms metric_groups enabled, quantiles metric_groups disabled and settings: default is applied to the group
    Then 1 dataset(s) have validity valid and 0 have validity invalid in 30 seconds
        And the container logs should contain the message "completed RPC subscription to group" within 30 seconds
        And this agent's heartbeat shows that 1 groups are matching the agent
        And this agent's heartbeat shows that 1 policies are applied and all has status running
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
        And referred sink must have active state on response within 30 seconds
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And remove the agent .yaml generated on each scenario


@smoke @metric_groups @metrics_netprobe
Scenario: netprobe handler with counters and quantiles metric groups enabled and histograms disabled
    Given the Orb user has a registered account
        And the Orb user logs in
        And that a sink already exists
    When an agent(input_type:netprobe, settings: {"test_type":"ping", "targets": {"www.google.com": {"target": "www.google.com"}, "orb_community": {"target": "orb.community"}}}) is self-provisioned via a configuration file on port available with 1 agent tags and has status online. [Overwrite default: False. Paste only file: True]
        And pktvisor state is running
        And 1 Agent Group(s) is created with all tags contained in the agent
        And a netprobe policy netprobe with tap_selector matching all tag(s) of the tap from an agent, counters, quantiles metric_groups enabled, histograms metric_groups disabled and settings: default is applied to the group
    Then 1 dataset(s) have validity valid and 0 have validity invalid in 30 seconds
        And the container logs should contain the message "completed RPC subscription to group" within 30 seconds
        And this agent's heartbeat shows that 1 groups are matching the agent
        And this agent's heartbeat shows that 1 policies are applied and all has status running
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
        And referred sink must have active state on response within 30 seconds
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And remove the agent .yaml generated on each scenario


@smoke @metric_groups @metrics_netprobe
Scenario: netprobe handler with histograms and quantiles metric groups enabled and counters disabled
    Given the Orb user has a registered account
        And the Orb user logs in
        And that a sink already exists
    When an agent(input_type:netprobe, settings: {"test_type":"ping", "targets": {"www.google.com": {"target": "www.google.com"}, "orb_community": {"target": "orb.community"}}}) is self-provisioned via a configuration file on port available with 1 agent tags and has status online. [Overwrite default: False. Paste only file: True]
        And pktvisor state is running
        And 1 Agent Group(s) is created with all tags contained in the agent
        And a netprobe policy netprobe with tap_selector matching all tag(s) of the tap from an agent, histograms, quantiles metric_groups enabled, counters metric_groups disabled and settings: default is applied to the group
    Then 1 dataset(s) have validity valid and 0 have validity invalid in 30 seconds
        And the container logs should contain the message "completed RPC subscription to group" within 30 seconds
        And this agent's heartbeat shows that 1 groups are matching the agent
        And this agent's heartbeat shows that 1 policies are applied and all has status running
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
        And referred sink must have active state on response within 30 seconds
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And remove the agent .yaml generated on each scenario

#### flow

@smoke @metric_groups @metrics_flow
Scenario: flow handler with default metric groups configuration
    Given the Orb user has a registered account
        And the Orb user logs in
        And that a sink already exists
    When an agent(input_type:flow, settings: {"bind":"0.0.0.0", "port":"available_port"}) is self-provisioned via a configuration file on port available with 1 agent tags and has status online. [Overwrite default: False. Paste only file: True]
        And pktvisor state is running
        And 1 Agent Group(s) is created with all tags contained in the agent
        And a flow policy flow with tap_selector matching all tag(s) of the tap from an agent, default metric_groups enabled, default metric_groups disabled and settings: default is applied to the group
    Then 1 dataset(s) have validity valid and 0 have validity invalid in 30 seconds
        And the container logs should contain the message "completed RPC subscription to group" within 30 seconds
        And this agent's heartbeat shows that 1 groups are matching the agent
        And this agent's heartbeat shows that 1 policies are applied and all has status running
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
        And referred sink must have active state on response within 30 seconds
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And remove the agent .yaml generated on each scenario


@smoke @metric_groups @metrics_flow
Scenario: flow handler with all metric groups enabled
    Given the Orb user has a registered account
        And the Orb user logs in
        And that a sink already exists
    When an agent(input_type:flow, settings: {"bind":"0.0.0.0", "port":"available_port"}) is self-provisioned via a configuration file on port available with 1 agent tags and has status online. [Overwrite default: False. Paste only file: True]
        And pktvisor state is running
        And 1 Agent Group(s) is created with all tags contained in the agent
        And a flow policy flow with tap_selector matching all tag(s) of the tap from an agent, all metric_groups enabled, none metric_groups disabled and settings: default is applied to the group
    Then 1 dataset(s) have validity valid and 0 have validity invalid in 30 seconds
        And the container logs should contain the message "completed RPC subscription to group" within 30 seconds
        And this agent's heartbeat shows that 1 groups are matching the agent
        And this agent's heartbeat shows that 1 policies are applied and all has status running
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
        And referred sink must have active state on response within 30 seconds
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And remove the agent .yaml generated on each scenario


@smoke @metric_groups @metrics_flow
Scenario: flow handler with all metric groups disabled
    Given the Orb user has a registered account
        And the Orb user logs in
        And that a sink already exists
    When an agent(input_type:flow, settings: {"bind":"0.0.0.0", "port":"available_port"}) is self-provisioned via a configuration file on port available with 1 agent tags and has status online. [Overwrite default: False. Paste only file: True]
        And pktvisor state is running
        And 1 Agent Group(s) is created with all tags contained in the agent
        And a flow policy flow with tap_selector matching all tag(s) of the tap from an agent, none metric_groups enabled, all metric_groups disabled and settings: default is applied to the group
    Then 1 dataset(s) have validity valid and 0 have validity invalid in 30 seconds
        And the container logs should contain the message "completed RPC subscription to group" within 30 seconds
        And this agent's heartbeat shows that 1 groups are matching the agent
        And this agent's heartbeat shows that 1 policies are applied and all has status running
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
        And referred sink must have active state on response within 30 seconds
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And remove the agent .yaml generated on each scenario


@smoke @metric_groups @metrics_flow
Scenario: flow handler with only cardinality metric groups enabled
    Given the Orb user has a registered account
        And the Orb user logs in
        And that a sink already exists
    When an agent(input_type:flow, settings: {"bind":"0.0.0.0", "port":"available_port"}) is self-provisioned via a configuration file on port available with 1 agent tags and has status online. [Overwrite default: False. Paste only file: True]
        And pktvisor state is running
        And 1 Agent Group(s) is created with all tags contained in the agent
        And a flow policy flow with tap_selector matching all tag(s) of the tap from an agent, cardinality metric_groups enabled, counters, by_packets, by_bytes, top_geo, conversations, top_ports, top_ips_ports, top_interfaces, top_ips metric_groups disabled and settings: default is applied to the group
    Then 1 dataset(s) have validity valid and 0 have validity invalid in 30 seconds
        And the container logs should contain the message "completed RPC subscription to group" within 30 seconds
        And this agent's heartbeat shows that 1 groups are matching the agent
        And this agent's heartbeat shows that 1 policies are applied and all has status running
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
        And referred sink must have active state on response within 30 seconds
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And remove the agent .yaml generated on each scenario


@smoke @metric_groups @metrics_flow
Scenario: flow handler with only counters metric groups enabled
    Given the Orb user has a registered account
        And the Orb user logs in
        And that a sink already exists
    When an agent(input_type:flow, settings: {"bind":"0.0.0.0", "port":"available_port"}) is self-provisioned via a configuration file on port available with 1 agent tags and has status online. [Overwrite default: False. Paste only file: True]
        And pktvisor state is running
        And 1 Agent Group(s) is created with all tags contained in the agent
        And a flow policy flow with tap_selector matching all tag(s) of the tap from an agent, counters metric_groups enabled, cardinality, by_packets, by_bytes, top_geo, conversations, top_ports, top_ips_ports, top_interfaces, top_ips metric_groups disabled and settings: default is applied to the group
    Then 1 dataset(s) have validity valid and 0 have validity invalid in 30 seconds
        And the container logs should contain the message "completed RPC subscription to group" within 30 seconds
        And this agent's heartbeat shows that 1 groups are matching the agent
        And this agent's heartbeat shows that 1 policies are applied and all has status running
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
        And referred sink must have active state on response within 30 seconds
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And remove the agent .yaml generated on each scenario


@smoke @metric_groups @metrics_flow
Scenario: flow handler with only by_packets metric groups enabled
    Given the Orb user has a registered account
        And the Orb user logs in
        And that a sink already exists
    When an agent(input_type:flow, settings: {"bind":"0.0.0.0", "port":"available_port"}) is self-provisioned via a configuration file on port available with 1 agent tags and has status online. [Overwrite default: False. Paste only file: True]
        And pktvisor state is running
        And 1 Agent Group(s) is created with all tags contained in the agent
        And a flow policy flow with tap_selector matching all tag(s) of the tap from an agent, by_packets metric_groups enabled, counters, cardinality, by_bytes, top_geo, conversations, top_ports, top_ips_ports, top_interfaces, top_ips metric_groups disabled and settings: default is applied to the group
    Then 1 dataset(s) have validity valid and 0 have validity invalid in 30 seconds
        And the container logs should contain the message "completed RPC subscription to group" within 30 seconds
        And this agent's heartbeat shows that 1 groups are matching the agent
        And this agent's heartbeat shows that 1 policies are applied and all has status running
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
        And referred sink must have active state on response within 30 seconds
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And remove the agent .yaml generated on each scenario


@smoke @metric_groups @metrics_flow
Scenario: flow handler with only by_bytes metric groups enabled
    Given the Orb user has a registered account
        And the Orb user logs in
        And that a sink already exists
    When an agent(input_type:flow, settings: {"bind":"0.0.0.0", "port":"available_port"}) is self-provisioned via a configuration file on port available with 1 agent tags and has status online. [Overwrite default: False. Paste only file: True]
        And pktvisor state is running
        And 1 Agent Group(s) is created with all tags contained in the agent
        And a flow policy flow with tap_selector matching all tag(s) of the tap from an agent, by_bytes metric_groups enabled, counters, by_packets, cardinality, top_geo, conversations, top_ports, top_ips_ports, top_interfaces, top_ips metric_groups disabled and settings: default is applied to the group
    Then 1 dataset(s) have validity valid and 0 have validity invalid in 30 seconds
        And the container logs should contain the message "completed RPC subscription to group" within 30 seconds
        And this agent's heartbeat shows that 1 groups are matching the agent
        And this agent's heartbeat shows that 1 policies are applied and all has status running
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
        And referred sink must have active state on response within 30 seconds
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And remove the agent .yaml generated on each scenario


@smoke @metric_groups @metrics_flow
Scenario: flow handler with only top_geo metric groups enabled
    Given the Orb user has a registered account
        And the Orb user logs in
        And that a sink already exists
    When an agent(input_type:flow, settings: {"bind":"0.0.0.0", "port":"available_port"}) is self-provisioned via a configuration file on port available with 1 agent tags and has status online. [Overwrite default: False. Paste only file: True]
        And pktvisor state is running
        And 1 Agent Group(s) is created with all tags contained in the agent
        And a flow policy flow with tap_selector matching all tag(s) of the tap from an agent, top_geo metric_groups enabled, counters, by_packets, by_bytes, cardinality, conversations, top_ports, top_ips_ports, top_interfaces, top_ips metric_groups disabled and settings: default is applied to the group
    Then 1 dataset(s) have validity valid and 0 have validity invalid in 30 seconds
        And the container logs should contain the message "completed RPC subscription to group" within 30 seconds
        And this agent's heartbeat shows that 1 groups are matching the agent
        And this agent's heartbeat shows that 1 policies are applied and all has status running
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
        And referred sink must have active state on response within 30 seconds
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And remove the agent .yaml generated on each scenario


@smoke @metric_groups @metrics_flow
Scenario: flow handler with only conversations metric groups enabled
    Given the Orb user has a registered account
        And the Orb user logs in
        And that a sink already exists
    When an agent(input_type:flow, settings: {"bind":"0.0.0.0", "port":"available_port"}) is self-provisioned via a configuration file on port available with 1 agent tags and has status online. [Overwrite default: False. Paste only file: True]
        And pktvisor state is running
        And 1 Agent Group(s) is created with all tags contained in the agent
        And a flow policy flow with tap_selector matching all tag(s) of the tap from an agent, conversations metric_groups enabled, counters, by_packets, by_bytes, top_geo, cardinality, top_ports, top_ips_ports, top_interfaces, top_ips metric_groups disabled and settings: default is applied to the group
    Then 1 dataset(s) have validity valid and 0 have validity invalid in 30 seconds
        And the container logs should contain the message "completed RPC subscription to group" within 30 seconds
        And this agent's heartbeat shows that 1 groups are matching the agent
        And this agent's heartbeat shows that 1 policies are applied and all has status running
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
        And referred sink must have active state on response within 30 seconds
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And remove the agent .yaml generated on each scenario


@smoke @metric_groups @metrics_flow
Scenario: flow handler with only top_ports metric groups enabled
    Given the Orb user has a registered account
        And the Orb user logs in
        And that a sink already exists
    When an agent(input_type:flow, settings: {"bind":"0.0.0.0", "port":"available_port"}) is self-provisioned via a configuration file on port available with 1 agent tags and has status online. [Overwrite default: False. Paste only file: True]
        And pktvisor state is running
        And 1 Agent Group(s) is created with all tags contained in the agent
        And a flow policy flow with tap_selector matching all tag(s) of the tap from an agent, top_ports metric_groups enabled, conversations, counters, by_packets, by_bytes, top_geo, cardinality, top_ips_ports, top_interfaces, top_ips metric_groups disabled and settings: default is applied to the group
    Then 1 dataset(s) have validity valid and 0 have validity invalid in 30 seconds
        And the container logs should contain the message "completed RPC subscription to group" within 30 seconds
        And this agent's heartbeat shows that 1 groups are matching the agent
        And this agent's heartbeat shows that 1 policies are applied and all has status running
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
        And referred sink must have active state on response within 30 seconds
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And remove the agent .yaml generated on each scenario


@smoke @metric_groups @metrics_flow
Scenario: flow handler with only top_ips_ports metric groups enabled
    Given the Orb user has a registered account
        And the Orb user logs in
        And that a sink already exists
    When an agent(input_type:flow, settings: {"bind":"0.0.0.0", "port":"available_port"}) is self-provisioned via a configuration file on port available with 1 agent tags and has status online. [Overwrite default: False. Paste only file: True]
        And pktvisor state is running
        And 1 Agent Group(s) is created with all tags contained in the agent
        And a flow policy flow with tap_selector matching all tag(s) of the tap from an agent, top_ips_ports metric_groups enabled, conversations, counters, by_packets, by_bytes, top_geo, cardinality, top_ports, top_interfaces, top_ips metric_groups disabled and settings: default is applied to the group
    Then 1 dataset(s) have validity valid and 0 have validity invalid in 30 seconds
        And the container logs should contain the message "completed RPC subscription to group" within 30 seconds
        And this agent's heartbeat shows that 1 groups are matching the agent
        And this agent's heartbeat shows that 1 policies are applied and all has status running
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
        And referred sink must have active state on response within 30 seconds
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And remove the agent .yaml generated on each scenario


@smoke @metric_groups @metrics_flow
Scenario: flow handler with only top_interfaces metric groups enabled
    Given the Orb user has a registered account
        And the Orb user logs in
        And that a sink already exists
    When an agent(input_type:flow, settings: {"bind":"0.0.0.0", "port":"available_port"}) is self-provisioned via a configuration file on port available with 1 agent tags and has status online. [Overwrite default: False. Paste only file: True]
        And pktvisor state is running
        And 1 Agent Group(s) is created with all tags contained in the agent
        And a flow policy flow with tap_selector matching all tag(s) of the tap from an agent, top_interfaces metric_groups enabled, conversations, counters, by_packets, by_bytes, top_geo, cardinality, top_ports, top_ips_ports, top_ips metric_groups disabled and settings: default is applied to the group
    Then 1 dataset(s) have validity valid and 0 have validity invalid in 30 seconds
        And the container logs should contain the message "completed RPC subscription to group" within 30 seconds
        And this agent's heartbeat shows that 1 groups are matching the agent
        And this agent's heartbeat shows that 1 policies are applied and all has status running
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
        And referred sink must have active state on response within 30 seconds
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And remove the agent .yaml generated on each scenario


@smoke @metric_groups @metrics_flow
Scenario: flow handler with only top_ips metric groups enabled
    Given the Orb user has a registered account
        And the Orb user logs in
        And that a sink already exists
    When an agent(input_type:flow, settings: {"bind":"0.0.0.0", "port":"available_port"}) is self-provisioned via a configuration file on port available with 1 agent tags and has status online. [Overwrite default: False. Paste only file: True]
        And pktvisor state is running
        And 1 Agent Group(s) is created with all tags contained in the agent
        And a flow policy flow with tap_selector matching all tag(s) of the tap from an agent, top_ips metric_groups enabled, conversations, counters, by_packets, by_bytes, top_geo, cardinality, top_ports, top_ips_ports, top_interfaces metric_groups disabled and settings: default is applied to the group
    Then 1 dataset(s) have validity valid and 0 have validity invalid in 30 seconds
        And the container logs should contain the message "completed RPC subscription to group" within 30 seconds
        And this agent's heartbeat shows that 1 groups are matching the agent
        And this agent's heartbeat shows that 1 policies are applied and all has status running
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
        And referred sink must have active state on response within 30 seconds
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And remove the agent .yaml generated on each scenario


#### pcap

@smoke @metric_groups @metrics_pcap
Scenario: pcap handler with default metric groups configuration
    Given the Orb user has a registered account
        And the Orb user logs in
        And that a sink already exists
    When an agent(input_type:pcap, settings: {"iface":"default"}) is self-provisioned via a configuration file on port available with 1 agent tags and has status online. [Overwrite default: False. Paste only file: True]
        And pktvisor state is running
        And 1 Agent Group(s) is created with all tags contained in the agent
        And a pcap policy pcap with tap_selector matching all tag(s) of the tap from an agent, default metric_groups enabled, default metric_groups disabled and settings: default is applied to the group
    Then 1 dataset(s) have validity valid and 0 have validity invalid in 30 seconds
        And the container logs should contain the message "completed RPC subscription to group" within 30 seconds
        And this agent's heartbeat shows that 1 groups are matching the agent
        And this agent's heartbeat shows that 1 policies are applied and all has status running
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
        And referred sink must have active state on response within 30 seconds
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And remove the agent .yaml generated on each scenario


@smoke @metric_groups @metrics_pcap
Scenario: pcap handler with all metric groups enabled
    Given the Orb user has a registered account
        And the Orb user logs in
        And that a sink already exists
    When an agent(input_type:pcap, settings: {"iface":"default"}) is self-provisioned via a configuration file on port available with 1 agent tags and has status online. [Overwrite default: False. Paste only file: True]
        And pktvisor state is running
        And 1 Agent Group(s) is created with all tags contained in the agent
        And a pcap policy pcap with tap_selector matching all tag(s) of the tap from an agent, all metric_groups enabled, none metric_groups disabled and settings: default is applied to the group
    Then 1 dataset(s) have validity valid and 0 have validity invalid in 30 seconds
        And the container logs should contain the message "completed RPC subscription to group" within 30 seconds
        And this agent's heartbeat shows that 1 groups are matching the agent
        And this agent's heartbeat shows that 1 policies are applied and all has status running
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
        And referred sink must have active state on response within 30 seconds
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And remove the agent .yaml generated on each scenario


@smoke @metric_groups @metrics_pcap
Scenario: pcap handler with all metric groups disabled
    Given the Orb user has a registered account
        And the Orb user logs in
        And that a sink already exists
    When an agent(input_type:pcap, settings: {"iface":"default"}) is self-provisioned via a configuration file on port available with 1 agent tags and has status online. [Overwrite default: False. Paste only file: True]
        And pktvisor state is running
        And 1 Agent Group(s) is created with all tags contained in the agent
        And a pcap policy pcap with tap_selector matching all tag(s) of the tap from an agent, none metric_groups enabled, all metric_groups disabled and settings: default is applied to the group
    Then 1 dataset(s) have validity valid and 0 have validity invalid in 30 seconds
        And the container logs should contain the message "completed RPC subscription to group" within 30 seconds
        And this agent's heartbeat shows that 1 groups are matching the agent
        And this agent's heartbeat shows that 1 policies are applied and all has status running
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
        And referred sink must have active state on response within 30 seconds
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And remove the agent .yaml generated on each scenario


#### bgp

@smoke @metric_groups @metrics_bgp
Scenario: bgp handler with default metric groups configuration
    Given the Orb user has a registered account
        And the Orb user logs in
        And that a sink already exists
    When an agent(input_type:pcap, settings: {"iface":"default"}) is self-provisioned via a configuration file on port available with 1 agent tags and has status online. [Overwrite default: False. Paste only file: True]
        And pktvisor state is running
        And 1 Agent Group(s) is created with all tags contained in the agent
        And a bgp policy pcap with tap_selector matching all tag(s) of the tap from an agent, default metric_groups enabled, default metric_groups disabled and settings: default is applied to the group
    Then 1 dataset(s) have validity valid and 0 have validity invalid in 30 seconds
        And the container logs should contain the message "completed RPC subscription to group" within 30 seconds
        And this agent's heartbeat shows that 1 groups are matching the agent
        And this agent's heartbeat shows that 1 policies are applied and all has status running
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
        And referred sink must have active state on response within 30 seconds
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And remove the agent .yaml generated on each scenario


@smoke @metric_groups @metrics_bgp
Scenario: bgp handler with all metric groups enabled
    Given the Orb user has a registered account
        And the Orb user logs in
        And that a sink already exists
    When an agent(input_type:pcap, settings: {"iface":"default"}) is self-provisioned via a configuration file on port available with 1 agent tags and has status online. [Overwrite default: False. Paste only file: True]
        And pktvisor state is running
        And 1 Agent Group(s) is created with all tags contained in the agent
        And a bgp policy pcap with tap_selector matching all tag(s) of the tap from an agent, all metric_groups enabled, none metric_groups disabled and settings: default is applied to the group
    Then 1 dataset(s) have validity valid and 0 have validity invalid in 30 seconds
        And the container logs should contain the message "completed RPC subscription to group" within 30 seconds
        And this agent's heartbeat shows that 1 groups are matching the agent
        And this agent's heartbeat shows that 1 policies are applied and all has status running
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
        And referred sink must have active state on response within 30 seconds
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And remove the agent .yaml generated on each scenario


@smoke @metric_groups @metrics_bgp
Scenario: bgp handler with all metric groups disabled
    Given the Orb user has a registered account
        And the Orb user logs in
        And that a sink already exists
    When an agent(input_type:pcap, settings: {"iface":"default"}) is self-provisioned via a configuration file on port available with 1 agent tags and has status online. [Overwrite default: False. Paste only file: True]
        And pktvisor state is running
        And 1 Agent Group(s) is created with all tags contained in the agent
        And a bgp policy pcap with tap_selector matching all tag(s) of the tap from an agent, none metric_groups enabled, all metric_groups disabled and settings: default is applied to the group
    Then 1 dataset(s) have validity valid and 0 have validity invalid in 30 seconds
        And the container logs should contain the message "completed RPC subscription to group" within 30 seconds
        And this agent's heartbeat shows that 1 groups are matching the agent
        And this agent's heartbeat shows that 1 policies are applied and all has status running
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
        And referred sink must have active state on response within 30 seconds
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And remove the agent .yaml generated on each scenario


#### dhcp

@smoke @metric_groups @metrics_dhcp
Scenario: dhcp handler with default metric groups configuration
    Given the Orb user has a registered account
        And the Orb user logs in
        And that a sink already exists
    When an agent(input_type:pcap, settings: {"iface":"default"}) is self-provisioned via a configuration file on port available with 1 agent tags and has status online. [Overwrite default: False. Paste only file: True]
        And pktvisor state is running
        And 1 Agent Group(s) is created with all tags contained in the agent
        And a dhcp policy pcap with tap_selector matching all tag(s) of the tap from an agent, default metric_groups enabled, default metric_groups disabled and settings: default is applied to the group
    Then 1 dataset(s) have validity valid and 0 have validity invalid in 30 seconds
        And the container logs should contain the message "completed RPC subscription to group" within 30 seconds
        And this agent's heartbeat shows that 1 groups are matching the agent
        And this agent's heartbeat shows that 1 policies are applied and all has status running
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
        And referred sink must have active state on response within 30 seconds
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And remove the agent .yaml generated on each scenario


@smoke @metric_groups @metrics_dhcp
Scenario: dhcp handler with all metric groups enabled
    Given the Orb user has a registered account
        And the Orb user logs in
        And that a sink already exists
    When an agent(input_type:pcap, settings: {"iface":"default"}) is self-provisioned via a configuration file on port available with 1 agent tags and has status online. [Overwrite default: False. Paste only file: True]
        And pktvisor state is running
        And 1 Agent Group(s) is created with all tags contained in the agent
        And a dhcp policy pcap with tap_selector matching all tag(s) of the tap from an agent, all metric_groups enabled, none metric_groups disabled and settings: default is applied to the group
    Then 1 dataset(s) have validity valid and 0 have validity invalid in 30 seconds
        And the container logs should contain the message "completed RPC subscription to group" within 30 seconds
        And this agent's heartbeat shows that 1 groups are matching the agent
        And this agent's heartbeat shows that 1 policies are applied and all has status running
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
        And referred sink must have active state on response within 30 seconds
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And remove the agent .yaml generated on each scenario


@smoke @metric_groups @metrics_dhcp
Scenario: dhcp handler with all metric groups disabled
    Given the Orb user has a registered account
        And the Orb user logs in
        And that a sink already exists
    When an agent(input_type:pcap, settings: {"iface":"default"}) is self-provisioned via a configuration file on port available with 1 agent tags and has status online. [Overwrite default: False. Paste only file: True]
        And pktvisor state is running
        And 1 Agent Group(s) is created with all tags contained in the agent
        And a dhcp policy pcap with tap_selector matching all tag(s) of the tap from an agent, none metric_groups enabled, all metric_groups disabled and settings: default is applied to the group
    Then 1 dataset(s) have validity valid and 0 have validity invalid in 30 seconds
        And the container logs should contain the message "completed RPC subscription to group" within 30 seconds
        And this agent's heartbeat shows that 1 groups are matching the agent
        And this agent's heartbeat shows that 1 policies are applied and all has status running
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
        And referred sink must have active state on response within 30 seconds
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And remove the agent .yaml generated on each scenario


#### net

@smoke @metric_groups @metrics_net
Scenario: net handler with default metric groups configuration
    Given the Orb user has a registered account
        And the Orb user logs in
        And that a sink already exists
    When an agent(input_type:pcap, settings: {"iface":"default"}) is self-provisioned via a configuration file on port available with 1 agent tags and has status online. [Overwrite default: False. Paste only file: True]
        And pktvisor state is running
        And 1 Agent Group(s) is created with all tags contained in the agent
        And a net policy pcap with tap_selector matching all tag(s) of the tap from an agent, default metric_groups enabled, default metric_groups disabled and settings: default is applied to the group
    Then 1 dataset(s) have validity valid and 0 have validity invalid in 30 seconds
        And the container logs should contain the message "completed RPC subscription to group" within 30 seconds
        And this agent's heartbeat shows that 1 groups are matching the agent
        And this agent's heartbeat shows that 1 policies are applied and all has status running
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
        And referred sink must have active state on response within 30 seconds
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And remove the agent .yaml generated on each scenario


@smoke @metric_groups @metrics_net
Scenario: net handler with all metric groups enabled
    Given the Orb user has a registered account
        And the Orb user logs in
        And that a sink already exists
    When an agent(input_type:pcap, settings: {"iface":"default"}) is self-provisioned via a configuration file on port available with 1 agent tags and has status online. [Overwrite default: False. Paste only file: True]
        And pktvisor state is running
        And 1 Agent Group(s) is created with all tags contained in the agent
        And a net policy pcap with tap_selector matching all tag(s) of the tap from an agent, all metric_groups enabled, none metric_groups disabled and settings: default is applied to the group
    Then 1 dataset(s) have validity valid and 0 have validity invalid in 30 seconds
        And the container logs should contain the message "completed RPC subscription to group" within 30 seconds
        And this agent's heartbeat shows that 1 groups are matching the agent
        And this agent's heartbeat shows that 1 policies are applied and all has status running
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
        And referred sink must have active state on response within 30 seconds
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And remove the agent .yaml generated on each scenario


@smoke @metric_groups @metrics_net
Scenario: net handler with all metric groups disabled
    Given the Orb user has a registered account
        And the Orb user logs in
        And that a sink already exists
    When an agent(input_type:pcap, settings: {"iface":"default"}) is self-provisioned via a configuration file on port available with 1 agent tags and has status online. [Overwrite default: False. Paste only file: True]
        And pktvisor state is running
        And 1 Agent Group(s) is created with all tags contained in the agent
        And a net policy pcap with tap_selector matching all tag(s) of the tap from an agent, none metric_groups enabled, all metric_groups disabled and settings: default is applied to the group
    Then 1 dataset(s) have validity valid and 0 have validity invalid in 30 seconds
        And the container logs should contain the message "completed RPC subscription to group" within 30 seconds
        And this agent's heartbeat shows that 1 groups are matching the agent
        And this agent's heartbeat shows that 1 policies are applied and all has status running
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
        And referred sink must have active state on response within 30 seconds
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And remove the agent .yaml generated on each scenario


@smoke @metric_groups @metrics_net
Scenario: net handler with only cardinality metric groups enabled
    Given the Orb user has a registered account
        And the Orb user logs in
        And that a sink already exists
    When an agent(input_type:pcap, settings: {"iface":"default"}) is self-provisioned via a configuration file on port available with 1 agent tags and has status online. [Overwrite default: False. Paste only file: True]
        And pktvisor state is running
        And 1 Agent Group(s) is created with all tags contained in the agent
        And a net policy pcap with tap_selector matching all tag(s) of the tap from an agent, cardinality metric_groups enabled, counters, top_geo, top_ips metric_groups disabled and settings: default is applied to the group
    Then 1 dataset(s) have validity valid and 0 have validity invalid in 30 seconds
        And the container logs should contain the message "completed RPC subscription to group" within 30 seconds
        And this agent's heartbeat shows that 1 groups are matching the agent
        And this agent's heartbeat shows that 1 policies are applied and all has status running
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
        And referred sink must have active state on response within 30 seconds
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And remove the agent .yaml generated on each scenario


@smoke @metric_groups @metrics_net
Scenario: net handler with only counters metric groups enabled
    Given the Orb user has a registered account
        And the Orb user logs in
        And that a sink already exists
    When an agent(input_type:pcap, settings: {"iface":"default"}) is self-provisioned via a configuration file on port available with 1 agent tags and has status online. [Overwrite default: False. Paste only file: True]
        And pktvisor state is running
        And 1 Agent Group(s) is created with all tags contained in the agent
        And a net policy pcap with tap_selector matching all tag(s) of the tap from an agent, counters metric_groups enabled, cardinality, top_geo, top_ips metric_groups disabled and settings: default is applied to the group
    Then 1 dataset(s) have validity valid and 0 have validity invalid in 30 seconds
        And the container logs should contain the message "completed RPC subscription to group" within 30 seconds
        And this agent's heartbeat shows that 1 groups are matching the agent
        And this agent's heartbeat shows that 1 policies are applied and all has status running
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
        And referred sink must have active state on response within 30 seconds
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And remove the agent .yaml generated on each scenario


@smoke @metric_groups @metrics_net
Scenario: net handler with only top_geo metric groups enabled
    Given the Orb user has a registered account
        And the Orb user logs in
        And that a sink already exists
    When an agent(input_type:pcap, settings: {"iface":"default"}) is self-provisioned via a configuration file on port available with 1 agent tags and has status online. [Overwrite default: False. Paste only file: True]
        And pktvisor state is running
        And 1 Agent Group(s) is created with all tags contained in the agent
        And a net policy pcap with tap_selector matching all tag(s) of the tap from an agent, top_geo metric_groups enabled, counters, cardinality, top_ips metric_groups disabled and settings: default is applied to the group
    Then 1 dataset(s) have validity valid and 0 have validity invalid in 30 seconds
        And the container logs should contain the message "completed RPC subscription to group" within 30 seconds
        And this agent's heartbeat shows that 1 groups are matching the agent
        And this agent's heartbeat shows that 1 policies are applied and all has status running
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
        And referred sink must have active state on response within 30 seconds
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And remove the agent .yaml generated on each scenario


@smoke @metric_groups @metrics_net
Scenario: net handler with only top_ips metric groups enabled
    Given the Orb user has a registered account
        And the Orb user logs in
        And that a sink already exists
    When an agent(input_type:pcap, settings: {"iface":"default"}) is self-provisioned via a configuration file on port available with 1 agent tags and has status online. [Overwrite default: False. Paste only file: True]
        And pktvisor state is running
        And 1 Agent Group(s) is created with all tags contained in the agent
        And a net policy pcap with tap_selector matching all tag(s) of the tap from an agent, top_ips metric_groups enabled, counters, top_geo, cardinality metric_groups disabled and settings: default is applied to the group
    Then 1 dataset(s) have validity valid and 0 have validity invalid in 30 seconds
        And the container logs should contain the message "completed RPC subscription to group" within 30 seconds
        And this agent's heartbeat shows that 1 groups are matching the agent
        And this agent's heartbeat shows that 1 policies are applied and all has status running
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
        And referred sink must have active state on response within 30 seconds
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And remove the agent .yaml generated on each scenario


#### dns

@smoke @metric_groups @metrics_dns
Scenario: dns handler with default metric groups configuration
    Given the Orb user has a registered account
        And the Orb user logs in
        And that a sink already exists
    When an agent(input_type:pcap, settings: {"iface":"default"}) is self-provisioned via a configuration file on port available with 1 agent tags and has status online. [Overwrite default: False. Paste only file: True]
        And pktvisor state is running
        And 1 Agent Group(s) is created with all tags contained in the agent
        And a dns policy pcap with tap_selector matching all tag(s) of the tap from an agent, default metric_groups enabled, default metric_groups disabled and settings: default is applied to the group
    Then 1 dataset(s) have validity valid and 0 have validity invalid in 30 seconds
        And the container logs should contain the message "completed RPC subscription to group" within 30 seconds
        And this agent's heartbeat shows that 1 groups are matching the agent
        And this agent's heartbeat shows that 1 policies are applied and all has status running
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
        And referred sink must have active state on response within 30 seconds
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And remove the agent .yaml generated on each scenario


@smoke @metric_groups @metrics_dns
Scenario: dns handler with all metric groups enabled
    Given the Orb user has a registered account
        And the Orb user logs in
        And that a sink already exists
    When an agent(input_type:pcap, settings: {"iface":"default"}) is self-provisioned via a configuration file on port available with 1 agent tags and has status online. [Overwrite default: False. Paste only file: True]
        And pktvisor state is running
        And 1 Agent Group(s) is created with all tags contained in the agent
        And a dns policy pcap with tap_selector matching all tag(s) of the tap from an agent, all metric_groups enabled, none metric_groups disabled and settings: default is applied to the group
    Then 1 dataset(s) have validity valid and 0 have validity invalid in 30 seconds
        And the container logs should contain the message "completed RPC subscription to group" within 30 seconds
        And this agent's heartbeat shows that 1 groups are matching the agent
        And this agent's heartbeat shows that 1 policies are applied and all has status running
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
        And referred sink must have active state on response within 30 seconds
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And remove the agent .yaml generated on each scenario


@smoke @metric_groups @metrics_dns
Scenario: dns handler with all metric groups disabled
    Given the Orb user has a registered account
        And the Orb user logs in
        And that a sink already exists
    When an agent(input_type:pcap, settings: {"iface":"default"}) is self-provisioned via a configuration file on port available with 1 agent tags and has status online. [Overwrite default: False. Paste only file: True]
        And pktvisor state is running
        And 1 Agent Group(s) is created with all tags contained in the agent
        And a dns policy pcap with tap_selector matching all tag(s) of the tap from an agent, none metric_groups enabled, all metric_groups disabled and settings: default is applied to the group
    Then 1 dataset(s) have validity valid and 0 have validity invalid in 30 seconds
        And the container logs should contain the message "completed RPC subscription to group" within 30 seconds
        And this agent's heartbeat shows that 1 groups are matching the agent
        And this agent's heartbeat shows that 1 policies are applied and all has status running
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
        And referred sink must have active state on response within 30 seconds
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And remove the agent .yaml generated on each scenario


@smoke @metric_groups @metrics_dns
Scenario: dns handler with only top_ecs metric groups enabled
    Given the Orb user has a registered account
        And the Orb user logs in
        And that a sink already exists
    When an agent(input_type:pcap, settings: {"iface":"default"}) is self-provisioned via a configuration file on port available with 1 agent tags and has status online. [Overwrite default: False. Paste only file: True]
        And pktvisor state is running
        And 1 Agent Group(s) is created with all tags contained in the agent
        And a dns policy pcap with tap_selector matching all tag(s) of the tap from an agent, top_ecs metric_groups enabled, top_qnames_details, cardinality, counters, dns_transaction, top_qnames, top_ports metric_groups disabled and settings: default is applied to the group
    Then 1 dataset(s) have validity valid and 0 have validity invalid in 30 seconds
        And the container logs should contain the message "completed RPC subscription to group" within 30 seconds
        And this agent's heartbeat shows that 1 groups are matching the agent
        And this agent's heartbeat shows that 1 policies are applied and all has status running
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
        And referred sink must have active state on response within 30 seconds
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And remove the agent .yaml generated on each scenario


@smoke @metric_groups @metrics_dns
Scenario: dns handler with only top_qnames_details metric groups enabled
    Given the Orb user has a registered account
        And the Orb user logs in
        And that a sink already exists
    When an agent(input_type:pcap, settings: {"iface":"default"}) is self-provisioned via a configuration file on port available with 1 agent tags and has status online. [Overwrite default: False. Paste only file: True]
        And pktvisor state is running
        And 1 Agent Group(s) is created with all tags contained in the agent
        And a dns policy pcap with tap_selector matching all tag(s) of the tap from an agent, top_qnames_details metric_groups enabled, top_ecs, cardinality, counters, dns_transaction, top_qnames, top_ports metric_groups disabled and settings: default is applied to the group
    Then 1 dataset(s) have validity valid and 0 have validity invalid in 30 seconds
        And the container logs should contain the message "completed RPC subscription to group" within 30 seconds
        And this agent's heartbeat shows that 1 groups are matching the agent
        And this agent's heartbeat shows that 1 policies are applied and all has status running
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
        And referred sink must have active state on response within 30 seconds
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And remove the agent .yaml generated on each scenario


@smoke @metric_groups @metrics_dns
Scenario: dns handler with only cardinality metric groups enabled
    Given the Orb user has a registered account
        And the Orb user logs in
        And that a sink already exists
    When an agent(input_type:pcap, settings: {"iface":"default"}) is self-provisioned via a configuration file on port available with 1 agent tags and has status online. [Overwrite default: False. Paste only file: True]
        And pktvisor state is running
        And 1 Agent Group(s) is created with all tags contained in the agent
        And a dns policy pcap with tap_selector matching all tag(s) of the tap from an agent, cardinality metric_groups enabled, top_qnames_details, top_ecs, counters, dns_transaction, top_qnames, top_ports metric_groups disabled and settings: default is applied to the group
    Then 1 dataset(s) have validity valid and 0 have validity invalid in 30 seconds
        And the container logs should contain the message "completed RPC subscription to group" within 30 seconds
        And this agent's heartbeat shows that 1 groups are matching the agent
        And this agent's heartbeat shows that 1 policies are applied and all has status running
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
        And referred sink must have active state on response within 30 seconds
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And remove the agent .yaml generated on each scenario


@smoke @metric_groups @metrics_dns
Scenario: dns handler with only counters metric groups enabled
    Given the Orb user has a registered account
        And the Orb user logs in
        And that a sink already exists
    When an agent(input_type:pcap, settings: {"iface":"default"}) is self-provisioned via a configuration file on port available with 1 agent tags and has status online. [Overwrite default: False. Paste only file: True]
        And pktvisor state is running
        And 1 Agent Group(s) is created with all tags contained in the agent
        And a dns policy pcap with tap_selector matching all tag(s) of the tap from an agent, counters metric_groups enabled, top_qnames_details, cardinality, top_ecs, dns_transaction, top_qnames, top_ports metric_groups disabled and settings: default is applied to the group
    Then 1 dataset(s) have validity valid and 0 have validity invalid in 30 seconds
        And the container logs should contain the message "completed RPC subscription to group" within 30 seconds
        And this agent's heartbeat shows that 1 groups are matching the agent
        And this agent's heartbeat shows that 1 policies are applied and all has status running
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
        And referred sink must have active state on response within 30 seconds
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And remove the agent .yaml generated on each scenario


@smoke @metric_groups @metrics_dns
Scenario: dns handler with only dns_transaction metric groups enabled
    Given the Orb user has a registered account
        And the Orb user logs in
        And that a sink already exists
    When an agent(input_type:pcap, settings: {"iface":"default"}) is self-provisioned via a configuration file on port available with 1 agent tags and has status online. [Overwrite default: False. Paste only file: True]
        And pktvisor state is running
        And 1 Agent Group(s) is created with all tags contained in the agent
        And a dns policy pcap with tap_selector matching all tag(s) of the tap from an agent, dns_transaction metric_groups enabled, top_qnames_details, cardinality, counters, top_ecs, top_qnames, top_ports metric_groups disabled and settings: default is applied to the group
    Then 1 dataset(s) have validity valid and 0 have validity invalid in 30 seconds
        And the container logs should contain the message "completed RPC subscription to group" within 30 seconds
        And this agent's heartbeat shows that 1 groups are matching the agent
        And this agent's heartbeat shows that 1 policies are applied and all has status running
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
        And referred sink must have active state on response within 30 seconds
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And remove the agent .yaml generated on each scenario


@smoke @metric_groups @metrics_dns
Scenario: dns handler with only top_qnames metric groups enabled
    Given the Orb user has a registered account
        And the Orb user logs in
        And that a sink already exists
    When an agent(input_type:pcap, settings: {"iface":"default"}) is self-provisioned via a configuration file on port available with 1 agent tags and has status online. [Overwrite default: False. Paste only file: True]
        And pktvisor state is running
        And 1 Agent Group(s) is created with all tags contained in the agent
        And a dns policy pcap with tap_selector matching all tag(s) of the tap from an agent, top_qnames metric_groups enabled, top_qnames_details, cardinality, counters, dns_transaction, top_ecs, top_ports metric_groups disabled and settings: default is applied to the group
    Then 1 dataset(s) have validity valid and 0 have validity invalid in 30 seconds
        And the container logs should contain the message "completed RPC subscription to group" within 30 seconds
        And this agent's heartbeat shows that 1 groups are matching the agent
        And this agent's heartbeat shows that 1 policies are applied and all has status running
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
        And referred sink must have active state on response within 30 seconds
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And remove the agent .yaml generated on each scenario


@smoke @metric_groups @metrics_dns
Scenario: dns handler with only top_ports metric groups enabled
    Given the Orb user has a registered account
        And the Orb user logs in
        And that a sink already exists
    When an agent(input_type:pcap, settings: {"iface":"default"}) is self-provisioned via a configuration file on port available with 1 agent tags and has status online. [Overwrite default: False. Paste only file: True]
        And pktvisor state is running
        And 1 Agent Group(s) is created with all tags contained in the agent
        And a dns policy pcap with tap_selector matching all tag(s) of the tap from an agent, top_ports metric_groups enabled, top_qnames_details, cardinality, counters, dns_transaction, top_qnames, top_ecs metric_groups disabled and settings: default is applied to the group
    Then 1 dataset(s) have validity valid and 0 have validity invalid in 30 seconds
        And the container logs should contain the message "completed RPC subscription to group" within 30 seconds
        And this agent's heartbeat shows that 1 groups are matching the agent
        And this agent's heartbeat shows that 1 policies are applied and all has status running
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
        And referred sink must have active state on response within 30 seconds
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And remove the agent .yaml generated on each scenario