## Scenario: Agent subscription to group with policies after editing agent group's tags (editing tags after agent provision) 
Steps:
-  
1. Provision an agent with tags
2. Create a group with different tags
3. Create a sink
4. Create 1 policy
5. Create a dataset linking the group, the sink and the policy
6. Edit groups' tags changing the value to match with agent


Expected result:
-
- Agent heartbeat must show just one group matching
- Agent logs must show that agent is subscribed to the group