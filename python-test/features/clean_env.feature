Feature: clean env

@cleanup
Scenario: cleanup orb
  Given the Orb user logs in
  Then cleanup agents
  Then cleanup agent group
  Then cleanup sinks
  Then cleanup policies
  Then cleanup datasets
