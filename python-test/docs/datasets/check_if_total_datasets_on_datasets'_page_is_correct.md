## Scenario: Check if total datasets on datasets' page is correct 
## Steps:
1 - Create multiple datasets

- REST API Method: POST
- endpoint: /policies/dataset
- header: {authorization:token}

2 - Get all existing datasets

- REST API Method: GET
- endpoint: /policies/dataset

3 - On datasets' page (`orb.live/pages/datasets/list`) check the total number of datasets at the end of the dataset table

4 - Count the number of existing datasets

## Expected Result:
- Total datasets on API response, datasets page and the real number must be the same 

