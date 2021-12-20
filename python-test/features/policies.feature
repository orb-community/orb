@policies
Feature: policy creation

  Scenario: Create a policy
    Given the Orb user logs in
      And that an agent already exists and is online
    When a new policy is created
    Then referred policy must be listened on the orb policies list