@sink
Feature: Sink Management tests

@smoke_ui
Scenario: Create a new Sink Management
    Given that the Orb user logs in Orb UI
        And the user clicks on Sink Management on left menu
    When a sink is created through the UI with 1 orb tag
    Then the new sink is shown on the datatable
        And total number was increased in one unit
