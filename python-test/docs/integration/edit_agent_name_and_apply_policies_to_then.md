## Scenario: Edit agent name and apply policies to then 
Steps:
-  
1. Provision an agent with tags
2. Create a group with same tags as agent
3. Create a sink
4. Create 1 policy
5. Create a dataset linking the group, the sink and the policy
8. Edit agent name

Expected result:
-
- Agent heartbeat must show just one group matching
- Agent logs must show that agent is subscribed to the group
- The container logs contain the message "policy applied successfully" referred to the policy applied to the second group
- The container logs that were output after all policies have been applied contains the message "scraped metrics for policy" referred to each applied policy
