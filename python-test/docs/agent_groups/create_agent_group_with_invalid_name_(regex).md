## Scenario: Create agent group with invalid name (regex) 
## Steps:
1 - Create an agent group using an invalid regex to agent group name

- REST API Method: POST
- endpoint: /agent_groups
- header: {authorization:token}
- example of invalid regex:

 * name starting with non-alphabetic characters
 * name with just 1 letter
 * space-separated composite name

## Expected Result:
- Request must fail with status code 400 (bad request) and no group must be created 
