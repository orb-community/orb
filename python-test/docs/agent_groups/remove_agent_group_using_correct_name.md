## Scenario: Remove agent group using correct name 
## Steps:
1 - Create an agent group

- REST API Method: POST
- endpoint: /agent_groups
- header: {authorization:token}

2 - On agent groups' page (`orb.live/pages/fleet/groups`) click on remove button
3 - Insert the name of the group correctly on delete modal
4 - Confirm the operation by clicking on "I UNDERSTAND, DELETE THIS AGENT GROUP" button

## Expected Result:
- Agent group must be deleted 

 
