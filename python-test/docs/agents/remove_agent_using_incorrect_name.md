## Scenario: Remove agent using incorrect name 
## Steps: 
1 - Create an agent

- REST API Method: POST
- endpoint: /agents
- header: {authorization:token}

2 - On agents' page (`orb.live/pages/fleet/agents`) click on remove button
3 - Insert the name of the agent incorrectly on delete modal

## Expected Result:
- "I UNDERSTAND, DELETE THIS AGENT" button must not be enabled
- After user close the deletion modal, agent must not be deleted 

