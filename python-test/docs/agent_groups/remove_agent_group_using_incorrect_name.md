## Scenario: Remove agent group using incorrect name 
## Steps: 
1 - Create an agent group

- REST API Method: POST
- endpoint: /agent_groups
- header: {authorization:token}

2 - On agent groups' page (`orb.live/pages/fleet/groups`) click on remove button
3 - Insert the name of the group incorrectly on delete modal
 
## Expected Result: 
- "I UNDERSTAND, DELETE THIS AGENT GROUP" button must not be enabled
- After user close the deletion modal, agent group must not be deleted  
