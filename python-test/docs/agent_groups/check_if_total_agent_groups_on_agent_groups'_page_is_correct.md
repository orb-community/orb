## Scenario: Check if total agent groups on agent groups' page is correct 
## Steps:
1 - Create multiple agent groups

- REST API Method: POST
- endpoint: /agent_groups
- header: {authorization:token}

2 - Get all existing agent groups

- REST API Method: GET
- endpoint: /agent_groups

3 - On agent groups' page (`orb.live/pages/fleet/groups`) check the total number of agent groups at the end of the agent groups table

4 - Count the number of existing agent groups

## Expected Result:
- Total agent groups on API response, agent groups page and the real number must be the same 

