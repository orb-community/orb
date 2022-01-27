## Scenario: Create dataset with invalid name (regex) 
## Steps:
1 - Create a dataset using an invalid regex to dataset name

- REST API Method: POST
- endpoint: /policies/dataset
- header: {authorization:token}
- example of invalid regex:

* name starting with non-alphabetic characters
* name with just 1 letter
* space-separated composite name

## Expected Result:
- Request must fail with status code 400 (bad request) and no dataset must be created 