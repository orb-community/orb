## Scenario: Edit Agent Group tags to unsubscribe agent 
Steps:
-  
1. Provision an agent with tags
2. Create a group with same tags
3. Edit groups' tags changing the value

Expected result:
-
- Agent heartbeat must show 0 group matching
- Agent logs must show that agent is unsubscribed to the group