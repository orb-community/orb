## Scenario: Remove dataset using correct name 
## Steps:
1 - Create a dataset

- REST API Method: POST
- endpoint: /policies/dataset
- header: {authorization:token}

2 - On datasets' page (`orb.live/pages/datasets/list`) click on remove button
3 - Insert the name of the dataset correctly on delete modal
4 - Confirm the operation by clicking on "I UNDERSTAND, DELETE THIS DATASET" button

## Expected Result:
- Dataset must be deleted 
