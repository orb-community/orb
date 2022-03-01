@policies
Feature: policy creation

  Scenario: Create a policy dns with description, host specification, bpf filter, pcap source, only qname suffix and only rcode
    Given the Orb user has a registered account
      And the Orb user logs in
      And that an agent with 1 orb tag(s) already exists and is online
    When a new policy is created using: handler=dns, description='policy_dns', host_specification=10.0.1.0/24,10.0.2.1/32,2001:db8::/64, bpf_filter_expression=udp port 53, pcap_source=libpcap, only_qname_suffix=[.foo.com/ .example.com], only_rcode=0
    Then referred policy must be listed on the orb policies list


  Scenario: Create a policy dns with host specification, bpf filter, pcap source, only qname suffix and only rcode
    Given the Orb user has a registered account
      And the Orb user logs in
      And that an agent with 1 orb tag(s) already exists and is online
    When a new policy is created using: handler=dns, host_specification=10.0.1.0/24,10.0.2.1/32,2001:db8::/64, bpf_filter_expression=udp port 53, pcap_source=libpcap, only_qname_suffix=[.foo.com/ .example.com], only_rcode=2
    Then referred policy must be listed on the orb policies list


  Scenario: Create a policy dns with bpf filter, pcap source, only qname suffix and only rcode
    Given the Orb user has a registered account
      And the Orb user logs in
      And that an agent with 1 orb tag(s) already exists and is online
    When a new policy is created using: handler=dns, bpf_filter_expression=udp port 53, pcap_source=af_packet, only_qname_suffix=[.foo.com/ .example.com], only_rcode=3
    Then referred policy must be listed on the orb policies list


  Scenario: Create a policy dns with pcap source, only qname suffix and only rcode
    Given the Orb user has a registered account
      And the Orb user logs in
      And that an agent with 1 orb tag(s) already exists and is online
    When a new policy is created using: handler=dns, pcap_source=af_packet, only_qname_suffix=[.foo.com/ .example.com], only_rcode=5
    Then referred policy must be listed on the orb policies list


  Scenario: Create a policy dns with only qname suffix
    Given the Orb user has a registered account
      And the Orb user logs in
      And that an agent with 1 orb tag(s) already exists and is online
    When a new policy is created using: handler=dns, only_qname_suffix=[.foo.com/ .example.com]
    Then referred policy must be listed on the orb policies list





  Scenario: Create a policy dhcp with description, host specification, bpf filter, pcap source, only qname suffix and only rcode
    Given the Orb user has a registered account
      And the Orb user logs in
      And that an agent with 1 orb tag(s) already exists and is online
    When a new policy is created using: handler=dhcp, description='policy_dhcp', host_specification=10.0.1.0/24,10.0.2.1/32,2001:db8::/64, bpf_filter_expression=udp port 53, pcap_source=libpcap, only_qname_suffix=[.foo.com/ .example.com], only_rcode=0
    Then referred policy must be listed on the orb policies list


  Scenario: Create a policy dhcp with host specification, bpf filter, pcap source, only qname suffix and only rcode
    Given the Orb user has a registered account
      And the Orb user logs in
      And that an agent with 1 orb tag(s) already exists and is online
    When a new policy is created using: handler=dhcp, host_specification=10.0.1.0/24,10.0.2.1/32,2001:db8::/64, bpf_filter_expression=udp port 53, pcap_source=libpcap, only_qname_suffix=[.foo.com/ .example.com], only_rcode=2
    Then referred policy must be listed on the orb policies list


  Scenario: Create a policy dhcp with bpf filter, pcap source, only qname suffix and only rcode
    Given the Orb user has a registered account
      And the Orb user logs in
      And that an agent with 1 orb tag(s) already exists and is online
    When a new policy is created using: handler=dhcp, bpf_filter_expression=udp port 53, pcap_source=af_packet, only_qname_suffix=[.foo.com/ .example.com], only_rcode=3
    Then referred policy must be listed on the orb policies list


  Scenario: Create a policy dhcp with pcap source, only qname suffix and only rcode
    Given the Orb user has a registered account
      And the Orb user logs in
      And that an agent with 1 orb tag(s) already exists and is online
    When a new policy is created using: handler=dhcp, pcap_source=af_packet, only_qname_suffix=[.foo.com/ .example.com], only_rcode=5
    Then referred policy must be listed on the orb policies list


  Scenario: Create a policy dhcp with only qname suffix
    Given the Orb user has a registered account
      And the Orb user logs in
      And that an agent with 1 orb tag(s) already exists and is online
    When a new policy is created using: handler=dhcp, only_qname_suffix=[.foo.com/ .example.com]
    Then referred policy must be listed on the orb policies list





  Scenario: Create a policy net with description, host specification, bpf filter, pcap source, only qname suffix and only rcode
    Given the Orb user has a registered account
      And the Orb user logs in
      And that an agent with 1 orb tag(s) already exists and is online
    When a new policy is created using: handler=net, description='policy_net', host_specification=10.0.1.0/24,10.0.2.1/32,2001:db8::/64, bpf_filter_expression=udp port 53, pcap_source=libpcap, only_qname_suffix=[.foo.com/ .example.com], only_rcode=0
    Then referred policy must be listed on the orb policies list


  Scenario: Create a policy net with host specification, bpf filter, pcap source, only qname suffix and only rcode
    Given the Orb user has a registered account
      And the Orb user logs in
      And that an agent with 1 orb tag(s) already exists and is online
    When a new policy is created using: handler=net, host_specification=10.0.1.0/24,10.0.2.1/32,2001:db8::/64, bpf_filter_expression=udp port 53, pcap_source=libpcap, only_qname_suffix=[.foo.com/ .example.com], only_rcode=2
    Then referred policy must be listed on the orb policies list


  Scenario: Create a policy net with bpf filter, pcap source, only qname suffix and only rcode
    Given the Orb user has a registered account
      And the Orb user logs in
      And that an agent with 1 orb tag(s) already exists and is online
    When a new policy is created using: handler=net, bpf_filter_expression=udp port 53, pcap_source=af_packet, only_qname_suffix=[.foo.com/ .example.com], only_rcode=3
    Then referred policy must be listed on the orb policies list


  Scenario: Create a policy net with pcap source, only qname suffix and only rcode
    Given the Orb user has a registered account
      And the Orb user logs in
      And that an agent with 1 orb tag(s) already exists and is online
    When a new policy is created using: handler=net, pcap_source=af_packet, only_qname_suffix=[.foo.com/ .example.com], only_rcode=5
    Then referred policy must be listed on the orb policies list


  Scenario: Create a policy net with only qname suffix
    Given the Orb user has a registered account
      And the Orb user logs in
      And that an agent with 1 orb tag(s) already exists and is online
    When a new policy is created using: handler=net, only_qname_suffix=[.foo.com/ .example.com]
    Then referred policy must be listed on the orb policies list