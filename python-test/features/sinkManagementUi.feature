@sink @AUTORETRY
Feature: Sink Management tests using orb UI

@smoke_ui
Scenario: Create a new Sink Management 
    Given that the Orb user logs in Orb UI
        And the user clicks on Sink Management on left menu
    When a sink is created through the UI with 1 orb tag
    Then the new sink is shown on the datatable
        And total number was increased in one unit

@smoke_ui
Scenario: Delete an Sink Management by name
    Given that the Orb user logs in Orb UI
        And the user clicks on Sink Management on left menu
    When delete a sink using filter by name with 1 orb tag
    Then the new sink is not shown on the datatable
        And total number was decreased in one unit

@smoke_ui
Scenario: Edit an Sink Management by name
    Given that the Orb user logs in Orb UI
        And the user clicks on Sink Management on left menu
    When update a sink using filter by name with 1 orb tag
    Then the new sink is shown on the datatable
        And total number was the same
