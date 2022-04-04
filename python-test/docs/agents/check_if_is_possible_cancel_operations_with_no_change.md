## Scenario: Check if is possible cancel operations with no change 
## Steps: 
1 - Create an agent

- REST API Method: POST
- endpoint: /agents
- header: {authorization:token}

2 - On agents' page (`orb.live/pages/fleet/agents`) click on edit button
3 - Change agents' name and click "next"
4 - Change agent's tag and click "next"
5 - Click "back" until return to agents' page

## Expected Result: 
- No changes must have been applied to the agent

