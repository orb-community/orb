## Scenario: Create agent group without description 
## Steps:

1 - Create an agent groups with no description

- REST API Method: POST
- endpoint: /agent_groups
- header: {authorization:token}


## Expected Result:
- Request must have status code 201 (created) and the agent group must be created
