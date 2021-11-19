@agents
Feature: agent provider
   
    Scenario: Provision agent
        Given A valid authentication
        When Create an agent
            And Build agent container 
        Then Agente should be online
            And Container logs should be sending capabilities
