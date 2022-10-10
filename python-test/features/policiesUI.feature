@policies_ui
Feature: Create policies using orb ui

  @MUTE
  Scenario: Create policy with tap pcap and handler dns through wizard editor
    Given that the Orb user logs in Orb UI
      And that an agent with 1 orb tag(s) already exists and is online
      And the user clicks on Policy Management on left menu
    When a new policy is created through the UI with: handler=dns, description=policy_dns, host_specification=10.0.1.0/24,10.0.2.1/32,2001:db8::/64, bpf_filter_expression=udp port 53, pcap_source=libpcap, only_qname_suffix=[.foo.com/ .example.com], only_rcode=2, exclude_noerror=true
    Then user must be directed to the policy view page
      And created policy must have the chosen parameters
      And created policy must is displayed on policy pages

  @smoke_ui
  Scenario: Create policy with tap pcap and handler net through wizard editor
    Given that the Orb user logs in Orb UI
      And that an agent with 1 orb tag(s) already exists and is online
      And the user clicks on Policy Management on left menu
    When a new policy is created through the UI with: handler=net, description=policy_net
      Then user must be directed to the policy view page
      And created policy must have the chosen parameters
      And created policy must is displayed on policy pages

  @smoke_ui
  Scenario: Create policy with tap pcap and handler dhcp through wizard editor
    Given that the Orb user logs in Orb UI
      And that an agent with 1 orb tag(s) already exists and is online
      And the user clicks on Policy Management on left menu
    When a new policy is created through the UI with: handler=dhcp, description=policy_dhcp
    Then user must be directed to the policy view page
      And created policy must have the chosen parameters
      And created policy must is displayed on policy pages

  @smoke_ui
  Scenario: Remove policy from Orb UI
    Given that the Orb user logs in Orb UI
      And that an agent with 1 orb tag(s) already exists and is online
      And the user clicks on Policy Management on left menu
      And a new policy is created through the UI with: handler=dhcp, description=policy_dhcp
      And user must be directed to the policy view page
      And the user clicks on Policy Management on left menu
    When remove policy from Orb UI
    Then created policy must is not displayed on policy pages
