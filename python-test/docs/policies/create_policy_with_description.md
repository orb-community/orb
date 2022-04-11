## Scenario: Create policy with description 
## Steps:

1 - Create a policy with description

- REST API Method: POST
- endpoint: /policies/agent/
- header: {authorization:token}


## Expected Result:
- Request must have status code 201 (created) and the policy must be created
