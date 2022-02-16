## Scenario: Test datasets filter 
## Steps:
1 - Create multiple datasets

- REST API Method: POST
- endpoint: /policies/dataset
- header: {authorization:token}

2 - On datasets' page (`orb.live/pages/datasets/list`) use the filter:

* Name
* Search by


## Expected Result:

- All filters must be working properly
