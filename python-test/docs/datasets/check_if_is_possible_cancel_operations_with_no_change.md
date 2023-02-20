## Scenario: Check if is possible cancel operations with no change 
## Steps:
1 - Create a dataset

- REST API Method: POST
- endpoint: /policies/dataset
- header: {authorization:token}

2 - On datasets' page (`orb.live/pages/datasets/list`) click on edit button
3 - Change groups' name and click "next"
4 - Change sink linked and click "next"

## Expected Result:
- No changes must have been applied to the dataset