@agents
Feature: agent provider
   
    Scenario: Provision agent
        Given the Orb user logs in
        When a new agent is created
            And the agent container is started
        Then the agent status in Orb should be online
            And the container logs should contain the message "sending capabilities" within 10 seconds


    Scenario: Apply two policies to an agent
        Given the Orb user logs in
            And that an agent already exists and is online
            And referred agent is subscribed to a group
            And that a sink already exists
            And that a policy already exists
            And that a dataset using referred group, sink and policy ID already exists
        When a new policy is created
            And a new dataset is created using referred group, sink and policy ID
        Then this agent's heartbeat shows that all 2 policies have been successfully applied
            And the container logs contains the message "policy applied successfully" referred to each policy within 10 seconds
            And the container logs that were output after all policies have been applied contains the message "scraped metrics for policy" referred to each policy within 180 seconds
            And referred sink must have active state on response within 10 seconds
