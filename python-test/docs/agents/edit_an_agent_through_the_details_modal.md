## Scenario: Edit an agent through the details modal 
## Steps: 
1 - Create an agent

- REST API Method: POST
- endpoint: /agents
- header: {authorization:token}

2 - On agents' page (`orb.live/pages/fleet/agents`) click on details button
3 - Click on "edit" button

## Expected Result: 
- User should be redirected to this agent's edit page and should be able to make changes

