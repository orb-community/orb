@datasets
Feature: datasets creation

  Scenario: Create Dataset
    Given the Orb user has a registered account
    And the Orb user logs in
    And that an agent already exists and is online
    And referred agent is subscribed to a group
    And that a sink already exists
    And that a policy already exists
    When a new dataset is created using referred group, sink and policy ID
    Then the container logs should contain the message "managing agent policy from core" within 10 seconds
    And the container logs should contain the message "policy applied successfully" within 10 seconds
    And the container logs should contain the message "scraped metrics for policy" within 180 seconds
    And referred sink must have active state on response within 10 seconds