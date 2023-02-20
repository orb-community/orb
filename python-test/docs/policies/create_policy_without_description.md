## Scenario: Create policy without description 
## Steps:

1 - Create a policy without description

- REST API Method: POST
- endpoint: /policies/agent/
- header: {authorization:token}


## Expected Result:
- Request must have status code 201 (created) and the policy must be created