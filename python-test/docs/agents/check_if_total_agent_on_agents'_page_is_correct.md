## Scenario: Check if total agent on agents' page is correct 
## Steps: 
1 - Create multiple agents

- REST API Method: POST
- endpoint: /agents
- header: {authorization:token}

2 - Get all existing agents

- REST API Method: GET
- endpoint: /agents
 
3 - On agents' page (`orb.live/pages/fleet/agents`) check the total number of agents at the end of the agents table

4 - Count the number of existing agents

## Expected Result: 
- Total agents on API response, agents page and the real number must be the same 
