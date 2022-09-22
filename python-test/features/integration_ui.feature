Feature: Integration UI tests

  @smoke_ui
  Scenario: Apply one policies to an agent (only dataset UI)
    Given that the Orb user logs in Orb UI
      And that an agent with 1 orb tag(s) already exists and is online
      And referred agent is subscribed to 1 group
      And this agent's heartbeat shows that 1 groups are matching the agent
      And a new policy is created using: handler=dns, description='policy_dns'
      And that a sink already exists
    When a dataset is created through the UI
    Then the policy must have status running on agent view page (Active Policies/Datasets)
      And policy and dataset are clickable and redirect user to referred pages
      And group must be listed on agent view page (Active Groups)
