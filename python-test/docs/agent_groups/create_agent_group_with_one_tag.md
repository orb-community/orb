## Scenario: Create agent group with one tag 
## Steps:

1 - Create an agent groups with one pair (key:value) of tags

- REST API Method: POST
- endpoint: /agent_groups
- header: {authorization:token}


## Expected Result:
- Request must have status code 201 (created) and the agent group must be created