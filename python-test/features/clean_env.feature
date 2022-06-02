Feature: clean env

@cleanup
Scenario: cleanup agents
  Given the Orb user logs in
  Then cleanup agents

@cleanup
Scenario: cleanup agent groups
  Given the Orb user logs in
  Then cleanup agent group

@cleanup
Scenario: cleanup sinks
  Given the Orb user logs in
  Then cleanup sinks

@cleanup
Scenario: cleanup policies
  Given the Orb user logs in
  Then cleanup policies

@cleanup
Scenario: cleanup datasets
  Given the Orb user logs in
  Then cleanup datasets

@cleanup
Scenario: cleanup yaml file
  Given the Orb user logs in
  Then remove all the agents .yaml generated on test process
