## Scenario: Subscribe an agent to multiple groups created before agent provisioning 
## Steps:
1. Create a group with one tag
2. Create another group with 2 tags
3. Provision an agent with the same tags as the two groups
4. Check agent's logs and agent's heartbeat


## Expected Result:
1 - Logs must display the message "completed RPC subscription to group" referred to both groups
2 - Agent's heartbeat must have 2 groups linked