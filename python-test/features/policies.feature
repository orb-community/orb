@policies
Feature: policy creation

  @smoke
  Scenario: Create a policy with dns handler, description, host specification, bpf filter, pcap source, only qname suffix and only rcode
    Given the Orb user has a registered account
      And the Orb user logs in
      And that an agent with 1 orb tag(s) already exists and is online
    When a new policy is created using: handler=dns, description='policy_dns', host_specification=10.0.1.0/24,10.0.2.1/32,2001:db8::/64, bpf_filter_expression=udp port 53, pcap_source=libpcap, only_qname_suffix=[.foo.com/ .example.com], only_rcode=0
    Then referred policy must be listed on the orb policies list


  @smoke
  Scenario: Create a policy with dns handler, host specification, bpf filter, pcap source, only qname suffix and only rcode
    Given the Orb user has a registered account
      And the Orb user logs in
      And that an agent with 1 orb tag(s) already exists and is online
    When a new policy is created using: handler=dns, host_specification=10.0.1.0/24,10.0.2.1/32,2001:db8::/64, bpf_filter_expression=udp port 53, pcap_source=libpcap, only_qname_suffix=[.foo.com/ .example.com], only_rcode=2
    Then referred policy must be listed on the orb policies list


  @smoke
  Scenario: Create a policy with dns handler, bpf filter, pcap source, only qname suffix and only rcode
    Given the Orb user has a registered account
      And the Orb user logs in
      And that an agent with 1 orb tag(s) already exists and is online
    When a new policy is created using: handler=dns, bpf_filter_expression=udp port 53, pcap_source=af_packet, only_qname_suffix=[.foo.com/ .example.com], only_rcode=3
    Then referred policy must be listed on the orb policies list


  @smoke
  Scenario: Create a policy with dns handler, pcap source, only qname suffix and only rcode
    Given the Orb user has a registered account
      And the Orb user logs in
      And that an agent with 1 orb tag(s) already exists and is online
    When a new policy is created using: handler=dns, pcap_source=af_packet, only_qname_suffix=[.foo.com/ .example.com], only_rcode=5
    Then referred policy must be listed on the orb policies list


  @smoke
  Scenario: Create a policy with dns handler, only qname suffix
    Given the Orb user has a registered account
      And the Orb user logs in
      And that an agent with 1 orb tag(s) already exists and is online
    When a new policy is created using: handler=dns, only_qname_suffix=[.foo.com/ .example.com]
    Then referred policy must be listed on the orb policies list


  @smoke
  Scenario: Create a policy with dhcp handler, description, host specification, bpf filter and pcap source
    Given the Orb user has a registered account
      And the Orb user logs in
      And that an agent with 1 orb tag(s) already exists and is online
    When a new policy is created using: handler=dhcp, description='policy_dhcp', host_specification=10.0.1.0/24,10.0.2.1/32,2001:db8::/64, bpf_filter_expression=udp port 53, pcap_source=libpcap
    Then referred policy must be listed on the orb policies list


  @smoke
  Scenario: Create a policy with net handler, description, host specification, bpf filter and pcap source
    Given the Orb user has a registered account
      And the Orb user logs in
      And that an agent with 1 orb tag(s) already exists and is online
    When a new policy is created using: handler=net, description='policy_net', host_specification=10.0.1.0/24,10.0.2.1/32,2001:db8::/64, bpf_filter_expression=udp port 53, pcap_source=libpcap
    Then referred policy must be listed on the orb policies list
