## Scenario: Check agent groups details 
## Steps:
1 - Create an agent group

- REST API Method: POST
- endpoint: /agent_groups
- header: {authorization:token}

2 - Get an agent group

- REST API Method: GET
- endpoint: /agent_groups/agent_group_id

## Expected Result:
- Status code must be 200 and the group name, description, matches against and tags must be returned on response

 
