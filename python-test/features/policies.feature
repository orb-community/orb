@policies
Feature: policy creation

  Scenario: Create a policy
    Given the Orb user has a registered account
      And the Orb user logs in
      And that an agent already exists and is online
    When a new policy is created
    Then referred policy must be listed on the orb policies list