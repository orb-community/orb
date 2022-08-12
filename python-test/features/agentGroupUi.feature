@group
Feature: Agent Group tests


@smoke_ui
Scenario: Create a new Agent Group
    Given the Orb user logs in through the UI
        And the user clicks on new agent group on left menu
    When a new agent group is created through the UI with 1 orb tag
    Then the new agent group is shown on the datatable

@smoke_ui
Scenario: Create a new Agent Group with decription
    Given the Orb user logs in through the UI
        And the user clicks on new agent group on left menu
    When a new agent group with description is created through the UI with 1 orb tag
    Then the new agent group is shown on the datatable
    
@smoke_ui
Scenario: Delete an Agent Group
    Given the Orb user logs in through the UI
        And the user clicks on new agent group on left menu
    When delete the agent group using filter by name with 1 orb tag
    Then the agent group is not shown on the datatable
    And total number was decreased in one unit
