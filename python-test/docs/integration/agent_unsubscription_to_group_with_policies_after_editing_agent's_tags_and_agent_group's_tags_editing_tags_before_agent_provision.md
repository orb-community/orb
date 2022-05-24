## Scenario: Agent unsubscription to group with policies after editing agent's tags and agent group's tags (editing tags before agent provision) 
Steps:
-  
1. Create an agent with tags
2. Create a group with another tag
3. Create a sink
4. Create 1 policy
5. Create a dataset linking the group, the sink and the policy
6. Edit agent tags using the same tag as the group
7. Edit groups' tags using a different one
8. Provision the agent

Expected result:
-
- Agent heartbeat must show just one group matching
- Agent logs must show that agent is unsubscribed to the group
