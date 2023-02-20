## Scenario: Edit Agent Group name, description and tags 
## Steps:
1 - Create an agent group

- REST API Method: POST
- endpoint: /agent_groups
- header: {authorization:token}

2- Edit this group name, description and tags

- REST API Method: PUT
- endpoint: /agent_groups/group_id
- header: {authorization:token}


## Expected Result:
- Request must have status code 200 (ok) and changes must be applied
