@login_ui
Feature: login on orb

  Scenario: Login on orb site
    Given that the Orb user logs in Orb UI
    Then the user should have access to orb home page
