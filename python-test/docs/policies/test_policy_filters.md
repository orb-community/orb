## Scenario: Test policy filters 
## Steps:
1 - Create multiple policies

- REST API Method: POST
- endpoint: /policies/agent/
- header: {authorization:token}

2 - On policies' page (`orb.live/pages/datasets/policies`) use the filter:

   * Name
   * Description
   * Version
   * Search by


## Expected Result:

- All filters must be working properly
