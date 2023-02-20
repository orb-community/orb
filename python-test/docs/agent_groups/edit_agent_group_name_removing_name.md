## Scenario: Edit Agent Group name removing name 
## Steps:
1 - Create an agent group

- REST API Method: POST
- endpoint: /agent_groups
- header: {authorization:token}

2- Edit this group name using None

- REST API Method: PUT
- endpoint: /agent_groups/group_id
- header: {authorization:token}


## Expected Result:
- Request must have status code 400 (error) and changes must not be applied
