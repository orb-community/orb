## Scenario: Create dataset 
## Steps:

1 - Create a dataset with no description

- REST API Method: POST
- endpoint: /policies/dataset
- header: {authorization:token}


## Expected Result:
- Request must have status code 201 (created) and the agent group must be created