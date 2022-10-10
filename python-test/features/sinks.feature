@sinks @AUTORETRY
Feature: sink creation

#  @smoke
  @MUTE
  Scenario: Create Sink using Prometheus
    Given that the user has the prometheus/grafana credentials
      And the Orb user has a registered account
      And the Orb user logs in
    When a new sink is created
    Then referred sink must have new state on response within 30 seconds


@smoke @sanity
Scenario: Create sink with name conflict
  Given that the user has the prometheus/grafana credentials
    And the Orb user has a registered account
    And the Orb user logs in
    And a new sink is created
  When a new sink is is requested to be created with the same name as an existent one
    Then the error message on response is failed to create Sink

@MUTE
Scenario: Edit sink using an already existent name (conflict)
  Given that the user has the prometheus/grafana credentials
    And the Orb user has a registered account
    And the Orb user logs in
    And a new sink is created
    And a new sink is created
  When the name of last Sink is edited using an already existent one
  Then the error message on response is entity already exists
