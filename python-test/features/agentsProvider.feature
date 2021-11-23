@agents
Feature: agent provider
   
    Scenario: Provision agent
        Given A valid authentication
        When Create an agent
            And Run agent container 
        Then Agent should be online
            And Container logs should be sending capabilities
