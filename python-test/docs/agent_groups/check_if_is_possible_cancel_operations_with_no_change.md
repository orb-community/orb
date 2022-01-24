## Scenario: Check if is possible cancel operations with no change 
## Steps:
1 - Create an agent group

- REST API Method: POST
- endpoint: /agent_groups
- header: {authorization:token}

2 - On agent groups' page (`orb.live/pages/fleet/groups`) click on edit button
3 - Change groups' name and click "next"
4 - Change groups' description and click "next"
4 - Change groups' tag and click "next"
5 - Click "back" until return to agent groups' page

## Expected Result:
- No changes must have been applied to the agent group
 
