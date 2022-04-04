## Scenario: Create agent group without tag 

## Steps:
1 - Create an agent with no pair (key:value) of tags

- REST API Method: POST
- endpoint: /agent_groups
- header: {authorization:token}


## Expected Result:
- Request must fail with status code 400 (bad request) and the agent group must not be created
 
