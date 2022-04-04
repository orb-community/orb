## Scenario: Edit agent group tag 
## Steps:
1 - Create an agent group

- REST API Method: POST
- endpoint: /agent_groups
- header: {authorization:token}

2- Edit this group tag

- REST API Method: PUT
- endpoint: /agent_groups/group_id
- header: {authorization:token}


## Expected Result:
- Request must have status code 200 (ok) and changes must be applied
 
