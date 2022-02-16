## Scenario: Edit an agent group through the details modal 
## Steps:
1 - Create an agent group

- REST API Method: POST
- endpoint: /agent_groups
- header: {authorization:token}

2 - On agent groups' page (`orb.live/pages/fleet/groups`) click on details button
3 - Click on "edit" button

## Expected Result:
- User should be redirected to this agent group's edit page and should be able to make changes
