## Scenario: Agent unsubscription to group with policies after editing orb agent's tags and agent group's tags (editing tags after agent provision) 
Steps:
-  
1. Provision an agent with tags
2. Create a group with another tag
3. Create a sink
4. Create 1 policy
5. Create a dataset linking the group, the sink and the policy
6. Edit orb tags on agent using the same tag as the group
7. Edit groups' tags using a different one


Expected result:
-
- Agent heartbeat must show just one group matching
- Agent logs must show that agent is unsubscribed to the group
