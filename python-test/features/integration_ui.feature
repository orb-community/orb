@integration_ui @AUTORETRY
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


  @smoke_ui
  Scenario: Apply one policies created via wizard to an agent (all steps UI)
    Given that the Orb user logs in Orb UI
      And that the user is on the orb Agents page
      And a new agent is created through the UI with 1 orb tag(s)
      And the agent container is started using the command provided by the UI on an available port
      And the agents list and the agents view should display agent's status as Online within 180 seconds
      And the user clicks on Agent Groups on left menu
      And a new agent group is created through the UI with same tags as the agent
      And this agent's heartbeat shows that 1 groups are matching the agent
      And the user clicks on Policy Management on left menu
      And a new policy is created through the UI with: handler=dns, description=policy_dns
      And the user clicks on Sink Management on left menu
      And a sink is created through the UI with 1 orb tag
    When a dataset is created through the UI
    Then the policy must have status running on agent view page (Active Policies/Datasets)
      And policy and dataset are clickable and redirect user to referred pages
      And group must be listed on agent view page (Active Groups)




  @smoke_ui
  Scenario: Remove dataset
    Given that the Orb user logs in Orb UI
      And that an agent with 1 orb tag(s) already exists and is online
      And referred agent is subscribed to 1 group
      And this agent's heartbeat shows that 1 groups are matching the agent
      And a new policy is created using: handler=dns, description='policy_dns'
      And that a sink already exists
      And a dataset is created through the UI
      And the policy must have status running on agent view page (Active Policies/Datasets)
      And policy and dataset are clickable and redirect user to referred pages
      And group must be listed on agent view page (Active Groups)
    When the dataset is removed
    Then policy must be removed from the agent


  @smoke_ui
  Scenario: Edit dataset
    Given that the Orb user logs in Orb UI
      And that an agent with 1 orb tag(s) already exists and is online
      And referred agent is subscribed to 1 group
      And this agent's heartbeat shows that 1 groups are matching the agent
      And a new policy is created using: handler=dns, description='policy_dns'
      And that a sink already exists
      And a dataset is created through the UI
      And a new sink is created
    When the dataset is edited and one more sink is inserted and name is changed
    Then 2 sinks are linked to the dataset and the new name is displayed
