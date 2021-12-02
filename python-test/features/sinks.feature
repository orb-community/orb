@sinks
Feature: sink creation

  Scenario: Create Sink using Prometheus
    Given that there is a remote write endpoint to send Prometheus metrics to Grafana Cloud
      And that the user have a Grafana Cloud Prometheus username
      And that the user have a Grafana Cloud API Key with a role with metrics push privileges
      And that the user is logged in on orb account
    When a new sink is created
    Then referred sink must have unknown state on response

