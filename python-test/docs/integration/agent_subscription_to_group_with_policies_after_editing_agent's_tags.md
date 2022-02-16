## Scenario: Agent subscription to group with policies after editing agent's tags 
Steps:
-  
1. Provision an agent with tags
2. Create a group with same tags as agent
3. Create a sink
4. Create 1 policy
5. Create a dataset linking the group, the sink and the policy
6. Create another group with different tags
7. Create another policy and apply to the group
8. Edit agent orb tags to match with second group

Expected result:
-
- Agent heartbeat must show just one group matching
- Agent logs must show that agent is unsubscribed from the first group
- Agent logs must show that agent is subscribed to the second group
- The container logs contain the message "policy applied successfully" referred to the policy applied to the second group
- The container logs that were output after all policies have been applied contains the message "scraped metrics for policy" referred to each applied policy


## Scenario: Agent subscription to multiple groups with policies after editing agent's tags
Steps:
-  
1. Provision an agent with tags
2. Create a group with same tags as agent
3. Create a sink
4. Create 1 policy
5. Create a dataset linking the group, the sink and the policy
6. Create another group with different tags
7. Create another policy and apply to the group
8. Edit agent orb tags to match with both groups

Expected result:
-
- Agent heartbeat must show 2 group matching
- Agent logs must show that agent is unsubscribed from the first group
- Agent logs must show that agent is subscribed to the second group
- The container logs contain the message "policy applied successfully" referred to the policy applied to both groups
- The container logs that were output after all policies have been applied contains the message "scraped metrics for policy" referred to each applied policy
