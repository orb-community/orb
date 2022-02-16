## Scenario: Create sink with invalid name (regex) 
## Steps:
1 - Create a sink using an invalid regex to sink name

- REST API Method: POST
- endpoint: /sinks
- header: {authorization:token}
- example of invalid regex:
    * name starting with non-alphabetic characters
    * name with just 1 letter
    * space-separated composite name

## Expected Result:
- Request must fail with status code 400 (bad request) and no sink must be created 
