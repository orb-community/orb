## Scenario: Create policy with invalid name (regex) 
## Steps:
1 - Create an policy using an invalid regex to policy name

- REST API Method: POST
- endpoint: /policies/agent/
- header: {authorization:token}
- example of invalid regex:

* name starting with non-alphabetic characters
* name with just 1 letter
* space-separated composite name

## Expected Result:
- Request must fail with status code 400 (bad request) and no policy must be created 
