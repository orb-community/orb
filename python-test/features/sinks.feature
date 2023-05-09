@sinks @AUTORETRY
Feature: sink creation

  @smoke
  Scenario: Create Sink using Prometheus
    Given that the user has the prometheus/grafana credentials
      And the Orb user has a registered account
      And the Orb user logs in
    When a new sink is created
    Then referred sink must have unknown state on response within 30 seconds


@smoke @sanity
Scenario: Create sink with name conflict
  Given that the user has the prometheus/grafana credentials
    And the Orb user has a registered account
    And the Orb user logs in
    And a new sink is created
  When a new sink is is requested to be created with the same name as an existent one
    Then the error message on response is failed to create Sink

@sanity
Scenario: Edit sink using an already existent name (conflict)
  Given that the user has the prometheus/grafana credentials
    And the Orb user has a registered account
    And the Orb user logs in
    And a new sink is created
    And a new sink is created
  When the name of last Sink is edited using an already existent one
  Then the error message on response is entity already exists


@sanity @sink_partial_update
Scenario: Partial Update: updating only sink name
  Given that the user has the prometheus/grafana credentials
    And the Orb user has a registered account
    And the Orb user logs in
    And a new sink is created
  When the name of this sink is updated
  Then the name updates to the new value and other fields remains the same


  @sanity @sink_partial_update
Scenario: Partial Update: updating only sink description
  Given that the user has the prometheus/grafana credentials
    And the Orb user has a registered account
    And the Orb user logs in
    And a new sink is created
  When the description of this sink is updated
  Then the description updates to the new value and other fields remains the same


@sanity @sink_partial_update
Scenario: Partial Update: updating only sink tags
  Given that the user has the prometheus/grafana credentials
    And the Orb user has a registered account
    And the Orb user logs in
    And a new sink is created
  When the tags of this sink is updated
  Then the tags updates to the new value and other fields remains the same


@sanity @sink_partial_update
Scenario: Partial Update: updating only sink configs
  Given that the user has the prometheus/grafana credentials
    And the Orb user has a registered account
    And the Orb user logs in
    And a new sink is created
  When the config of this sink is updated
  Then the config updates to the new value and other fields remains the same


@sanity @sink_partial_update
Scenario: Partial Update: updating only sink name and description
  Given that the user has the prometheus/grafana credentials
    And the Orb user has a registered account
    And the Orb user logs in
    And a new sink is created
  When the name and description of this sink is updated
  Then the name and description updates to the new value and other fields remains the same


@sanity @sink_partial_update
Scenario: Partial Update: updating only sink name and configs
  Given that the user has the prometheus/grafana credentials
    And the Orb user has a registered account
    And the Orb user logs in
    And a new sink is created
  When the name and config of this sink is updated
  Then the name and config updates to the new value and other fields remains the same


@sanity @sink_partial_update
Scenario: Partial Update: updating only sink name and tags
  Given that the user has the prometheus/grafana credentials
    And the Orb user has a registered account
    And the Orb user logs in
    And a new sink is created
  When the name and tags of this sink is updated
  Then the name and tags updates to the new value and other fields remains the same


@sanity @sink_partial_update
Scenario: Partial Update: updating only sink description and tags
  Given that the user has the prometheus/grafana credentials
    And the Orb user has a registered account
    And the Orb user logs in
    And a new sink is created
  When the description and tags of this sink is updated
  Then the description and tags updates to the new value and other fields remains the same


@sanity @sink_partial_update
Scenario: Partial Update: updating only sink description and configs
  Given that the user has the prometheus/grafana credentials
    And the Orb user has a registered account
    And the Orb user logs in
    And a new sink is created
  When the description and config of this sink is updated
  Then the description and config updates to the new value and other fields remains the same


@sanity @sink_partial_update
Scenario: Partial Update: updating only sink tags and configs
  Given that the user has the prometheus/grafana credentials
    And the Orb user has a registered account
    And the Orb user logs in
    And a new sink is created
  When the tags and config of this sink is updated
  Then the tags and config updates to the new value and other fields remains the same


@sanity @sink_partial_update
Scenario: Partial Update: updating only sink name, description and tags
  Given that the user has the prometheus/grafana credentials
    And the Orb user has a registered account
    And the Orb user logs in
    And a new sink is created
  When the name, description and tags of this sink is updated
  Then the name, description and tags updates to the new value and other fields remains the same


@sanity @sink_partial_update
Scenario: Partial Update: updating only sink name, description and configs
  Given that the user has the prometheus/grafana credentials
    And the Orb user has a registered account
    And the Orb user logs in
    And a new sink is created
  When the name, description and config of this sink is updated
  Then the name, description and config updates to the new value and other fields remains the same


@sanity @sink_partial_update
Scenario: Partial Update: updating only sink name, tags and configs
  Given that the user has the prometheus/grafana credentials
    And the Orb user has a registered account
    And the Orb user logs in
    And a new sink is created
  When the name, tags and config of this sink is updated
  Then the name, tags and config updates to the new value and other fields remains the same


@sanity @sink_partial_update
Scenario: Partial Update: updating only sink description, tags and configs
  Given that the user has the prometheus/grafana credentials
    And the Orb user has a registered account
    And the Orb user logs in
    And a new sink is created
  When the description, tags and config of this sink is updated
  Then the description, tags and config updates to the new value and other fields remains the same
