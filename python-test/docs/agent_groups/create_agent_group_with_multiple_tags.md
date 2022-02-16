## Scenario: Create agent group with multiple tags 
## Steps:
1 - Create an agent group with more than one pair (key:value) of tags

- REST API Method: POST
- endpoint: /agent_groups
- header: {authorization:token}


## Expected Result:
- Request must have status code 201 (created) and the agent group must be created
- Groups with multiple tags will only match with agents with the same multiple tags