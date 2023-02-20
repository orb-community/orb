## Scenario: Create agent with duplicate name 
## Steps: 
1 - Create an agent

- REST API Method: POST
- endpoint: /agents
- header: {authorization:token}

2 - Create another agent using the same agent name
 
## Expected Result: 
- First request must have status code 201 (created) and one agent must be created on orb
- Second request must fail with status code 409 (conflict) and no other agent must be created (make sure that first agent has not been modified)
