## Scenario: Edit dataset sink 
## Steps: 
1 - Create a dataset

- REST API Method: POST
- endpoint: /policies/dataset
- header: {authorization:token}

2- Edit this dataset sink

- REST API Method: PUT
- endpoint: /policies/dataset/dataset_id
- header: {authorization:token}


## Expected Result:
- Request must have status code 200 (ok) and changes must be applied