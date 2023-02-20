## Scenario: Create agent with multiple tags 
## Steps: 
1 - Create an agent with more than one pair (key:value) of tags

- REST API Method: POST
- endpoint: /agents
- header: {authorization:token}


## Expected Result:
- Request must have status code 201 (created) and the agent must be created
- Agent with multiple tags will match each tag individually
