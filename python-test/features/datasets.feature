@datasets
Feature: datasets creation

  @smoke
  Scenario: Create Dataset
    Given the Orb user has a registered account
      And the Orb user logs in
      And that an agent with 1 orb tag(s) already exists and is online
      And referred agent is subscribed to a group
      And that a sink already exists
      And that a policy using: handler=dns, description='policy_dns', host_specification=10.0.1.0/24,10.0.2.1/32,2001:db8::/64, bpf_filter_expression=udp port 53, pcap_source=libpcap, only_qname_suffix=[.foo.com/ .example.com], only_rcode=0 already exists
    When a new dataset is created using referred group, policy and 1 sink
    Then the container logs should contain the message "managing agent policy from core" within 10 seconds
      And the container logs should contain the message "policy applied successfully" within 10 seconds
      And the container logs should contain the message "scraped metrics for policy" within 180 seconds
      And referred sink must have active state on response within 10 seconds
      And datasets related to all existing policies have validity valid