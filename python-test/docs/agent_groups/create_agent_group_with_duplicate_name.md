## Scenario: Create agent group with duplicate name 
## Steps:
1 - Create an agent group

- REST API Method: POST
- endpoint: /agent_groups
- header: {authorization:token}

2 - Create another agent group using the same agent group name

## Expected Result:
- First request must have status code 201 (created) and one group must be created on orb
- Second request must fail with status code 409 (conflict) and no other group must be created (make sure that first group has not been modified)
