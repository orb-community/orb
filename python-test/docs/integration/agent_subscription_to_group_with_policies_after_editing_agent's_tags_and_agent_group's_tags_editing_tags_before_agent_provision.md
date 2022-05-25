## Scenario: Agent subscription to group with policies after editing orb agent's tags and agent group's tags (editing tags before agent provision) 
Steps:
-  
1. Create an agent with tags
2. Create a group with same tags as agent
3. Create a sink
4. Create 1 policy
5. Create a dataset linking the group, the sink and the policy
6. Edit groups' tags changing the value
7. Edit agent orb tags to match with new groups tags
8. Provision the agent

Expected result:
-
- Agent heartbeat must show just one group matching
- Agent logs must show that agent is subscribed to the group
