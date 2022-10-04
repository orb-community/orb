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


@smoke
Scenario: Create duplicated net policy without insert new name
    Given the Orb user has a registered account
        And the Orb user logs in
        And that an agent with 1 orb tag(s) already exists and is online
        And a new policy is created using: handler=net, description='policy_net'
    When try to duplicate this policy 4 times without set new name
    Then 3 policies must be successfully duplicated and 1 must return an error


@smoke
Scenario: Create duplicated dhcp policy without insert new name
    Given the Orb user has a registered account
        And the Orb user logs in
        And that an agent with 1 orb tag(s) already exists and is online
        And a new policy is created using: handler=dhcp, description='policy_dhcp'
    When try to duplicate this policy 4 times without set new name
    Then 3 policies must be successfully duplicated and 1 must return an error


@smoke
Scenario: Create duplicated dns policy without insert new name
    Given the Orb user has a registered account
        And the Orb user logs in
        And that an agent with 1 orb tag(s) already exists and is online
        And a new policy is created using: handler=dns, description='policy_dns', host_specification=10.0.1.0/24,10.0.2.1/32,2001:db8::/64, bpf_filter_expression=udp port 53, pcap_source=libpcap, only_qname_suffix=[.foo.com/ .example.com], only_rcode=0
    When try to duplicate this policy 4 times without set new name
    Then 3 policies must be successfully duplicated and 1 must return an error


@smoke
Scenario: Create 4 duplicated policy with new name
    Given the Orb user has a registered account
        And the Orb user logs in
        And that an agent with 1 orb tag(s) already exists and is online
        And a new policy is created using: handler=dns, description='policy_dns', host_specification=10.0.1.0/24,10.0.2.1/32,2001:db8::/64, bpf_filter_expression=udp port 53, pcap_source=libpcap, only_qname_suffix=[.foo.com/ .example.com], only_rcode=0
    When try to duplicate this policy 4 times with a random new name
    Then 4 policies must be successfully duplicated and 0 must return an error


@smoke
Scenario: Create 3 duplicated dns policy without insert new name and 1 with new name
    Given the Orb user has a registered account
        And the Orb user logs in
        And that an agent with 1 orb tag(s) already exists and is online
        And a new policy is created using: handler=dns, description='policy_dns', host_specification=10.0.1.0/24,10.0.2.1/32,2001:db8::/64, bpf_filter_expression=udp port 53, pcap_source=libpcap, only_qname_suffix=[.foo.com/ .example.com], only_rcode=0
    When try to duplicate this policy 3 times without set new name
        And 3 policies must be successfully duplicated and 0 must return an error
        And try to duplicate this policy 2 times with a random new name
    Then 2 policies must be successfully duplicated and 0 must return an error


@smoke @sanity
Scenario: Create policy with name conflict
    Given the Orb user has a registered account
        And the Orb user logs in
        And a new policy is created using: handler=dns, description='policy_dns'
    When a new policy is requested to be created with the same name as an existent one and: handler=dhcp, description='policy_dns'
    Then the error message on response is failed to create policy

@smoke @sanity
Scenario: Edit policy using an already existent name (conflict)
    Given the Orb user has a registered account
        And the Orb user logs in
        And a new policy is created using: handler=dns, description='policy_dns'
    When editing a policy using name=conflict
    Then the error message on response is entity already exists
