## Scenario: Check policies details 
## Steps:
1 - Create a policy

- REST API Method: POST
- endpoint:  /policies/agent/
- header: {authorization:token}

2 - Get a policy

- REST API Method: GET
- endpoint:  /policies/agent/

## Expected Result:
- Status code must be 200 and policy name, description, backend, input details and handler must be returned on response
