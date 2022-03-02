@sinks
Feature: sink creation

  @smoke
  Scenario: Create Sink using Prometheus
    Given that the user has the prometheus/grafana credentials
      And the Orb user has a registered account
      And the Orb user logs in
    When a new sink is created
    Then referred sink must have unknown state on response within 10 seconds

