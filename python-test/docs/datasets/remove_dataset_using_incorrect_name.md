## Scenario: Remove dataset using incorrect name 
## Steps:
1 - Create a dataset

- REST API Method: POST
- endpoint: /policies/dataset
- header: {authorization:token}

2 - On datasets' page (`orb.live/pages/datasets/list`) click on remove button
3 - Insert the name of the dataset correctly on delete modal

## Expected Result:
- "I UNDERSTAND, DELETE THIS DATASET" button must not be enabled
- After user close the deletion modal, dataset must not be deleted  
