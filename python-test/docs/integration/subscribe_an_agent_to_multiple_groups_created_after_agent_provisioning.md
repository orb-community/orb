## Scenario: Subscribe an agent to multiple groups created after agent provisioning 

## Steps: 
1. Provision an agent with tags
2. Create a group with at least one tag equal to agent
3. Create another group with at least one tag equal to agent
4. Check agent's logs and agent's heartbeat
 

## Expected Result: 
1 - Logs must display the message "completed RPC subscription to group" referred to both groups
2 - Agent's heartbeat must have 2 groups linked
