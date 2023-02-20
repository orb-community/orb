## Scenario: Check datasets details 
## Steps:
1 - Create a dataset

- REST API Method: POST
- endpoint: /policies/dataset
- header: {authorization:token}

2 - Get a dataset

- REST API Method: GET
- endpoint: /policies/dataset/dataset_id

## Expected Result:
- Status code must be 200 and the dataset name, validity, agent group linked, agent policy linked and sink linked must be returned on response

 