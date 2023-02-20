## Scenario: provision an agent before create an agent group

Steps:
-  
1. Provision an agent with tags
2. Create a group with same tags as the agent

Expected result:
-
- The orb-agent container logs contain the message "completed RPC subscription to group"
- Group has one agent matching
- Agent status is online