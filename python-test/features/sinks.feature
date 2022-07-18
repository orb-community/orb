@sinks
Feature: sink creation

#  @smoke
  @MUTE
  Scenario: Create Sink using Prometheus
    Given that the user has the prometheus/grafana credentials
      And the Orb user has a registered account
      And the Orb user logs in
    When a new sink is created
    Then referred sink must have new state on response within 30 seconds
