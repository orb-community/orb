@group
Feature: Agent Group tests


@smoke_ui
Scenario: Create a new Agent Group
    Given that the Orb user logs in Orb UI
        And the user clicks on Agent Groups on left menu
    When a new agent group is created through the UI with 1 orb tag
    Then the new agent group is shown on the datatable

@smoke_ui
Scenario: Create a new Agent Group with decription
    Given that the Orb user logs in Orb UI
        And the user clicks on Agent Groups on left menu
    When a new agent group with description is created through the UI with 1 orb tag
    Then the new agent group is shown on the datatable
