## Scenario: Create agent with invalid name (regex) 
## Steps: 
1 - Create an agent using an invalid regex to agent name

- REST API Method: POST
- endpoint: /agents
- header: {authorization:token}
- example of invalid regex:

    * name starting with non-alphabetic characters
    * name with just 1 letter
    * space-separated composite name
 
## Expected Result: 
- Request must fail with status code 400 (bad request) and no agent must be created 
