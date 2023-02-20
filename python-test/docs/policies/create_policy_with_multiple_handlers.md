## Scenario: Create policy with multiple handlers 
## Steps:

1 - Create a policy with dns, net and dhcp handler

- REST API Method: POST
- endpoint: /policies/agent/
- header: {authorization:token}


## Expected Result:
- Request must have status code 201 (created) and the policy must be created