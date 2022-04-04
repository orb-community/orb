## Scenario: Remove agent using correct name 
## Steps: 
1 - Create an agent

- REST API Method: POST
- endpoint: /agents
- header: {authorization:token}

2 - On agents' page (`orb.live/pages/fleet/agents`) click on remove button
3 - Insert the name of the agent correctly on delete modal
4 - Confirm the operation by clicking on "I UNDERSTAND, DELETE THIS AGENT" button
 
## Expected Result: 
- Agent must be deleted 
