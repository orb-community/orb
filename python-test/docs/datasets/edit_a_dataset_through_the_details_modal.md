## Scenario: Edit a dataset through the details modal 
## Steps:
1 - Create a dataset

- REST API Method: POST
- endpoint: /policies/dataset
- header: {authorization:token}

2 - On datasets' page (`orb.live/pages/datasets/list`) click on details button
3 - Click on "edit" button

## Expected Result:
- User should be redirected to this dataset's edit page and should be able to make changes