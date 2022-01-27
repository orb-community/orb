## Scenario: provision an agent after create an agent group

Steps:
-  
1. Create a group with tags
2. Provision an agent with same tags as group

Expected result:
-
- The orb-agent container logs contain the message "completed RPC subscription to group"
- Group has one agent matching
- Agent status is online