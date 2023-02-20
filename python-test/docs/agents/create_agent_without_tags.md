## Scenario: Create agent without tags 

## Steps:
1 - Create an agent with no pair (key:value) of tags

- REST API Method: POST
- endpoint: /agents
- header: {authorization:token}


## Expected Result:
- Request must have status code 201 (created) and the agent must be created
 
