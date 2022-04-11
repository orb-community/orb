## Scenario: Edit Agent Group tags to subscribe agent

Steps:
-  
1. Provision an agent with tags
2. Create a group with different tags
3. Edit groups' tags changing the value to match with agent


Expected result:
-
- Agent heartbeat must show just one group matching
- Agent logs must show that agent is subscribed to the group