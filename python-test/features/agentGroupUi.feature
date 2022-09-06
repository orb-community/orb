@group
Feature: Agent Group tests using Orb UI


#@smoke_ui
#Scenario: Create a new Agent Group
#    Given that the Orb user logs in Orb UI
#        And the user clicks on Agent Groups on left menu
#   When a new agent group is created through the UI with 1 orb tag
#    Then the new agent group is shown on the datatable

#@smoke_ui
#Scenario: Create a new Agent Group with description
#    Given that the Orb user logs in Orb UI
#        And the user clicks on Agent Groups on left menu
#    When a new agent group with description is created through the UI with 1 orb tag
#    Then the new agent group is shown on the datatable

#@smoke_ui
#Scenario: Delete an Agent Group
#    Given that the Orb user logs in Orb UI
#        And the user clicks on Agent Groups on left menu
#    When delete the agent group using filter by name with 1 orb tag
#    Then the new agent group is not shown on the datatable
#    And total number was decreased in one unit

@smoke_ui
Scenario: Update an Agent Group by Name
    Given that the Orb user logs in Orb UI
        And the user clicks on Agent Groups on left menu
    When update the agent group using filter by name with 1 orb tag
    Then the new agent group is shown on the datatable
    And total number was the same
