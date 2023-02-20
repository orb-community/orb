## Scenario: Create policy with duplicate name 
## Steps:
1 - Create a policy

- REST API Method: POST
- endpoint: /policies/agent/
- header: {authorization:token}

2 - Create another policy using the same policy name

## Expected Result:
- First request must have status code 201 (created) and one policy must be created on orb
- Second request must fail with status code 409 (conflict) and no other policy must be created (make sure that first policy has not been modified)
